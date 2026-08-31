package additional

import (
	"math"
	"sort"
	"time"

	"ev_routing/internal/dto"
)

// SpecificConsumptionKWhPer100Km is E(T).
func SpecificConsumptionKWhPer100Km(temperature, e0 float64) float64 {
	switch {
	case temperature < 25:
		return e0 / (1 - 5.1e-5*math.Pow(25-temperature, 2.45))
	case temperature > 25:
		return e0 / (1 - 1.2e-3*math.Pow(temperature-25, 2))
	default:
		return e0
	}
}

// energyForDistanceKWh is E(T)/100 * S_ij.
func energyForDistanceKWh(distanceKm, consumptionPer100Km float64) float64 {
	return consumptionPer100Km / 100 * distanceKm
}

// candidateCount is k = 30(1 - R_EV/(R_EV+D_remain)).
func candidateCount(routeRequest *dto.RouteRequestDTO) int {
	consumption := SpecificConsumptionKWhPer100Km(routeRequest.Temperature, routeRequest.SpendOpt)
	remainingRangeKm := routeRequest.AccLevel / consumption * 100
	remainingTripKm := haversineDistanceKm(routeRequest.StartCoords, routeRequest.FinishCoords)

	k := MaxCandidateStations * (1 - remainingRangeKm/(remainingRangeKm+remainingTripKm))
	if k < 0 {
		k = 0
	}

	return int(math.Round(k))
}

// FilterCandidateStations is C_k = Top_k(RTree(Buffer(Line(S,D), 5km))).
// Top_k = least-detour stations; a linear scan stands in for the R-tree.
func FilterCandidateStations(routeRequest *dto.RouteRequestDTO) []dto.StationDTO {
	start := routeRequest.StartCoords
	finish := routeRequest.FinishCoords

	type candidate struct {
		station      dto.StationDTO
		lineDistance float64
	}

	inBuffer := make([]candidate, 0, len(routeRequest.FilteredStations))
	for _, station := range routeRequest.FilteredStations {
		lineDistance := distancePointToSegmentKm(station.Coords, start, finish)
		if lineDistance <= CandidateBufferKm {
			inBuffer = append(inBuffer, candidate{station: station, lineDistance: lineDistance})
		}
	}

	sort.Slice(inBuffer, func(i, j int) bool {
		return inBuffer[i].lineDistance < inBuffer[j].lineDistance
	})

	k := candidateCount(routeRequest)
	if k > len(inBuffer) {
		k = len(inBuffer)
	}
	inBuffer = inBuffer[:k]

	candidates := make([]dto.StationDTO, len(inBuffer))
	for i, c := range inBuffer {
		candidates[i] = c.station
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Id < candidates[j].Id
	})

	return candidates
}

// NearestStationTariff is r_ij^near.
func NearestStationTariff(candidates []dto.StationDTO, a, b dto.CoordsDTO) float64 {
	if len(candidates) == 0 {
		return 0
	}

	nearest := candidates[0]
	nearestDist := distancePointToSegmentKm(nearest.Coords, a, b)
	for _, candidate := range candidates[1:] {
		dist := distancePointToSegmentKm(candidate.Coords, a, b)
		if dist < nearestDist {
			nearest = candidate
			nearestDist = dist
		}
	}

	return cheapestActiveSlotPrice(nearest.Slots)
}

// cheapestActiveSlotPrice is min(r_jm) over active slots.
func cheapestActiveSlotPrice(slots []dto.SlotDTO) float64 {
	price := 0.0
	found := false
	for _, slot := range slots {
		if !slot.IsActive {
			continue
		}
		if !found || slot.Price < price {
			price = slot.Price
			found = true
		}
	}

	return price
}

// selectBestSlot is the x_jm decision: the active slot minimizing
// Z_w+Z_ch for chargeAmountKWh (sum_m x_jm <= 1).
func selectBestSlot(
	slots []dto.SlotDTO,
	chargeAmountKWh, consumptionPer100Km, avgSpeedKmH float64,
) (slot dto.SlotDTO, waitCost, chargeCost float64, chargeDuration time.Duration, ok bool) {
	bestTotal := math.MaxFloat64

	for _, s := range slots {
		if !s.IsActive || s.Power <= 0 {
			continue
		}

		sChargeCost := s.Price * chargeAmountKWh
		sWaitCost := consumptionPer100Km / 100 * s.Price * avgSpeedKmH * s.WaitTime
		total := sChargeCost + sWaitCost

		if total < bestTotal {
			bestTotal = total
			slot = s
			waitCost = sWaitCost
			chargeCost = sChargeCost
			chargeDuration = HoursToDuration(chargeAmountKWh / s.Power)
			ok = true
		}
	}

	return slot, waitCost, chargeCost, chargeDuration, ok
}

// BuildDirectedEdge is edge i->j under a "charge to full" policy: prices
// Z_s+Z_w+Z_ch, enforcing x_ij=0 and C_i >= E(T)/100*S_ij; ok=false omits it.
func BuildDirectedEdge(
	i, j dto.StationDTO,
	distanceKm float64,
	tripDuration time.Duration,
	avgSpeedKmH float64,
	consumptionPer100Km float64,
	nearTariff float64,
	routeRequest *dto.RouteRequestDTO,
) (dto.Edge, bool) {
	if j.Id == StartStationId || i.Id == FinishStationId {
		return dto.Edge{}, false
	}

	departureCharge := routeRequest.AccMax
	if i.Id == StartStationId {
		departureCharge = routeRequest.AccLevel
	}

	energyNeeded := energyForDistanceKWh(distanceKm, consumptionPer100Km)
	if departureCharge < energyNeeded {
		return dto.Edge{}, false
	}

	arrivalCharge := departureCharge - energyNeeded
	driveCost := nearTariff * energyNeeded

	edge := dto.Edge{
		Distance:     distanceKm,
		TripDuration: tripDuration,
		Cost:         driveCost,
	}

	if j.Id == FinishStationId {
		return edge, true
	}

	chargeAmount := routeRequest.AccMax - arrivalCharge
	if chargeAmount <= 0 {
		return edge, true
	}

	slot, waitCost, chargeCost, chargeDuration, ok := selectBestSlot(j.Slots, chargeAmount, consumptionPer100Km, avgSpeedKmH)
	if !ok {
		return dto.Edge{}, false
	}

	edge.SlotId = slot.Id
	edge.ChargeDuration = chargeDuration
	edge.Cost += waitCost + chargeCost

	return edge, true
}
