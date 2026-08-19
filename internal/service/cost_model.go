package service

import (
	"math"
	"sort"
	"time"

	"ev_routing/internal/dto"
)

const (
	// candidateBufferKm is the width of the corridor around the direct
	// S-D line that a station must fall within to be considered a
	// charging candidate.
	candidateBufferKm = 5.0

	// maxCandidateStations is the upper bound on the number of candidate
	// stations (k when the vehicle's remaining range is negligible
	// compared to the trip distance).
	maxCandidateStations = 30
)

// efficiencyKmPerKWh implements Q(T): the EV's energy efficiency, in
// kilometers of range per kWh, as a function of ambient temperature (°C).
func efficiencyKmPerKWh(temperature float64) float64 {
	return -0.00254*temperature*temperature + 0.110*temperature + 7.20
}

// energyForDistanceKWh converts a distance into the energy (kWh) needed to
// cover it at the given efficiency (km/kWh).
func energyForDistanceKWh(distanceKm, efficiency float64) float64 {
	return distanceKm / efficiency
}

// candidateCount implements k = 30(1 - R_EV/(R_EV+D_remain)): the number of
// charging-station candidates shrinks as the vehicle's remaining range
// (R_EV, from its current charge level) grows relative to the remaining
// straight-line trip distance (D_remain).
func candidateCount(routeRequest *dto.RouteRequestDTO) int {
	efficiency := efficiencyKmPerKWh(routeRequest.Temperature)
	remainingRangeKm := routeRequest.AccLevel * efficiency
	remainingTripKm := haversineDistanceKm(routeRequest.StartCoords, routeRequest.FinishCoords)

	k := maxCandidateStations * (1 - remainingRangeKm/(remainingRangeKm+remainingTripKm))
	if k < 0 {
		k = 0
	}

	return int(math.Round(k))
}

// filterCandidateStations implements
// C_k = Top_k(RTree(Buffer(Line(S,D), 5km))): it keeps stations within
// candidateBufferKm of the straight S-D line, then takes the k closest to
// that line (Top_k here means "least detour from the direct route" - the
// model doesn't state a ranking criterion, and this is the most natural
// one absent a stated alternative). No real R-tree is built since the
// station counts involved don't warrant one; a linear scan does the same
// filtering. The result is re-sorted by station ID to match the ordering
// buildOrderedStationList expects.
func filterCandidateStations(routeRequest *dto.RouteRequestDTO) []dto.StationDTO {
	start := routeRequest.StartCoords
	finish := routeRequest.FinishCoords

	type candidate struct {
		station      dto.StationDTO
		lineDistance float64
	}

	inBuffer := make([]candidate, 0, len(routeRequest.FilteredStations))
	for _, station := range routeRequest.FilteredStations {
		lineDistance := distancePointToSegmentKm(station.Coords, start, finish)
		if lineDistance <= candidateBufferKm {
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

// nearestStationTariff returns the price (r_near) of the candidate station
// nearest to the segment a-b, or 0 if there are no candidates to price the
// segment against.
func nearestStationTariff(candidates []dto.StationDTO, a, b dto.CoordsDTO) float64 {
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

	return float64(nearest.Price)
}

// buildDirectedEdge computes the directed i->j edge's cost under the
// deterministic "charge to full at every station visited" policy: the
// departure charge at i is routeRequest.AccLevel if i is the start, or
// AccMax otherwise (since every prior stop tops up to full before leaving).
// It enforces the model's constraints: no edge into the start (x_ij=0 when
// j=S), no edge out of the finish (x_ij=0 when i=D), and no edge whose
// energy requirement the departure charge can't cover (C_i >= S_ij/Q(T)).
// ok is false when the edge is forbidden or infeasible, in which case it
// should be omitted from the adjacency matrix entirely.
func buildDirectedEdge(
	i, j dto.StationDTO,
	distanceKm float64,
	tripDuration time.Duration,
	efficiency float64,
	nearTariff float64,
	routeRequest *dto.RouteRequestDTO,
) (dto.Edge, bool) {
	if j.Id == startStationId || i.Id == finishStationId {
		return dto.Edge{}, false
	}

	departureCharge := routeRequest.AccMax
	if i.Id == startStationId {
		departureCharge = routeRequest.AccLevel
	}

	energyNeeded := energyForDistanceKWh(distanceKm, efficiency)
	if departureCharge < energyNeeded {
		return dto.Edge{}, false
	}

	arrivalCharge := departureCharge - energyNeeded
	driveCost := nearTariff * energyNeeded

	var chargeCost float64
	var chargeDuration time.Duration
	if j.Id != finishStationId {
		chargeAmount := routeRequest.AccMax - arrivalCharge
		if chargeAmount < 0 {
			chargeAmount = 0
		}

		chargeCost = float64(j.Price) * chargeAmount
		if j.Power > 0 {
			chargeDuration = hoursToDuration(chargeAmount / float64(j.Power))
		}
	}

	return dto.Edge{
		Distance:       distanceKm,
		TripDuration:   tripDuration,
		ChargeDuration: chargeDuration,
		Cost:           driveCost + chargeCost,
	}, true
}
