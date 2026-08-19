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

// RoutingService builds the station graph used for route planning, using
// OpenRouteService to compute the distance/duration of each edge.
type RoutingService struct {
	ORS *Service
}

// NewRoutingService builds a RoutingService that computes edges via ors.
func NewRoutingService(ors *Service) *RoutingService {
	return &RoutingService{ORS: ors}
}

// GetAdjacencyMatrix orders the request's start point, finish point and
// candidate charging stations (see filterCandidateStations) into a single
// list, then calls OpenRouteService once per unordered station pair to
// compute each direction's cost separately: energy cost (Z_s) priced at the
// nearest candidate's tariff, plus charging cost (Z_ch) at the destination
// station under a charge-to-full policy. Directions that violate the
// model's constraints (no edge into the start, none out of the finish) or
// that the departure charge can't cover (C_i >= S_ij/Q(T)) are omitted from
// the matrix entirely, so callers can treat a missing edge as "not
// traversable".
func (rs *RoutingService) GetAdjacencyMatrix(routeRequest *dto.RouteRequestDTO) (map[int]map[int]dto.Edge, error) {
	stations := buildOrderedStationList(routeRequest)
	candidateStations := stations[1 : len(stations)-1]
	efficiency := efficiencyKmPerKWh(routeRequest.Temperature)

	adjacencyMatrix := make(map[int]map[int]dto.Edge, len(stations))
	for _, station := range stations {
		adjacencyMatrix[station.Id] = make(map[int]dto.Edge)
	}

	for i := range stations {
		for j := i + 1; j < len(stations); j++ {
			from := stations[i]
			to := stations[j]

			route, err := rs.ORS.GetRoute(
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

			if edge, ok := buildDirectedEdge(from, to, route.Distance, tripDuration, efficiency, nearTariff, routeRequest); ok {
				adjacencyMatrix[from.Id][to.Id] = edge
			}
			if edge, ok := buildDirectedEdge(to, from, route.Distance, tripDuration, efficiency, nearTariff, routeRequest); ok {
				adjacencyMatrix[to.Id][from.Id] = edge
			}
		}
	}

	return adjacencyMatrix, nil
}

// buildOrderedStationList merges the request's start point and finish point
// with its candidate charging stations (see filterCandidateStations) into a
// single list ordered start-first, finish-last, with the candidates sorted
// by ID in between.
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

// hoursToDuration converts a duration expressed in fractional hours (as
// returned by OpenRouteService) into a time.Duration.
func hoursToDuration(hours float64) time.Duration {
	return time.Duration(hours * float64(time.Hour))
}

// getRouteNodesFromIds walks the request's ordered station list and, for
// every station on the winning path, accumulates distance/cost/duration from
// the adjacency matrix edge connecting it to the previous station on the
// path. Shared by GeneticRouteService and DijkstraRouteService.
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
