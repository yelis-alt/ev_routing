package service

import (
	"fmt"
	"math"
	"time"

	"ev_routing/internal/dto"
)

const (
	startStationId  = math.MinInt32
	finishStationId = math.MaxInt32
)

// RoutingService builds the routing graph over V = {S,D} ∪ C_k, using
// OpenRouteService for each edge's distance/duration.
type RoutingService struct {
	ors *Service
}

// NewRoutingService builds a RoutingService that computes edges via ors.
func NewRoutingService(ors *Service) *RoutingService {
	return &RoutingService{ors: ors}
}

// GetAdjacencyMatrix builds V's edges (see buildDirectedEdge); a missing
// edge means that direction is forbidden or infeasible.
func (rs *RoutingService) GetAdjacencyMatrix(routeRequest *dto.RouteRequestDTO) (map[int]map[int]dto.Edge, error) {
	stations := buildOrderedStationList(routeRequest)
	candidateStations := stations[1 : len(stations)-1]
	consumption := specificConsumptionKWhPer100Km(routeRequest.Temperature, routeRequest.SpendOpt)

	adjacencyMatrix := make(map[int]map[int]dto.Edge, len(stations))
	for _, station := range stations {
		adjacencyMatrix[station.Id] = make(map[int]dto.Edge)
	}

	for i := range stations {
		for j := i + 1; j < len(stations); j++ {
			from := stations[i]
			to := stations[j]

			route, err := rs.ors.GetRoute(
				from.Coords.Latitude,
				from.Coords.Longitude,
				to.Coords.Latitude,
				to.Coords.Longitude,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"get route between station %d and %d: %w",
					from.Id, to.Id, err,
				)
			}

			tripDuration := hoursToDuration(route.TripDuration)
			nearTariff := nearestStationTariff(candidateStations, from.Coords, to.Coords)

			var avgSpeedKmH float64
			if route.TripDuration > 0 {
				avgSpeedKmH = route.Distance / route.TripDuration
			}

			if edge, ok := buildDirectedEdge(from, to, route.Distance, tripDuration, avgSpeedKmH, consumption, nearTariff, routeRequest); ok {
				adjacencyMatrix[from.Id][to.Id] = edge
			}
			if edge, ok := buildDirectedEdge(to, from, route.Distance, tripDuration, avgSpeedKmH, consumption, nearTariff, routeRequest); ok {
				adjacencyMatrix[to.Id][from.Id] = edge
			}
		}
	}

	return adjacencyMatrix, nil
}

// buildOrderedStationList is V ordered S, C_k (by id), D.
func buildOrderedStationList(routeRequest *dto.RouteRequestDTO) []dto.StationDTO {
	candidateStations := filterCandidateStations(routeRequest)

	startPoint := dto.StationDTO{Id: startStationId, Coords: routeRequest.StartCoords}
	finishPoint := dto.StationDTO{Id: finishStationId, Coords: routeRequest.FinishCoords}

	ordered := make([]dto.StationDTO, 0, len(candidateStations)+2)
	ordered = append(ordered, startPoint)
	ordered = append(ordered, candidateStations...)
	ordered = append(ordered, finishPoint)

	return ordered
}

// hoursToDuration converts fractional hours (as returned by
// OpenRouteService) into a time.Duration.
func hoursToDuration(hours float64) time.Duration {
	return time.Duration(hours * float64(time.Hour))
}

// getRouteNodesFromIds accumulates Z and distance/duration along pathIds.
// Shared by GeneticRouteService and DijkstraRouteService.
func getRouteNodesFromIds(
	adjacencyMatrix map[int]map[int]dto.Edge,
	pathIds []int,
	routeRequest *dto.RouteRequestDTO,
) []dto.RouteNodeDTO {
	pathIdsSet := make(map[int]struct{}, len(pathIds))
	for _, id := range pathIds {
		pathIdsSet[id] = struct{}{}
	}

	var reachDurationCum time.Duration
	distanceCum := 0.0
	costCum := 0.0
	var chargeDurationCumm time.Duration
	previousId := startStationId

	routeNodesList := make([]dto.RouteNodeDTO, 0, len(pathIds))
	for _, station := range buildOrderedStationList(routeRequest) {
		stationId := station.Id
		if _, ok := pathIdsSet[stationId]; !ok {
			continue
		}

		routeNode := dto.RouteNodeDTO{RouteNode: station}

		if stationId != startStationId {
			edge := adjacencyMatrix[previousId][stationId]

			chargeDuration := edge.ChargeDuration
			distanceCum += edge.Distance
			costCum += edge.Cost
			reachDurationCum += edge.TripDuration
			reachDurationCum += chargeDurationCumm
			chargeDurationCumm += chargeDuration

			distanceCum = roundToTwoDecimals(distanceCum)
			costCum = roundToTwoDecimals(costCum)

			routeNode.Distance = distanceCum
			routeNode.Cost = costCum
			routeNode.ChargeDuration = chargeDuration
			routeNode.ReachDuration = reachDurationCum
			routeNode.SlotId = edge.SlotId
		} else {
			routeNode.Distance = 0.0
			routeNode.Cost = 0.0
			routeNode.ChargeDuration = 0
			routeNode.ReachDuration = 0
		}

		previousId = stationId
		routeNodesList = append(routeNodesList, routeNode)
	}

	return routeNodesList
}
