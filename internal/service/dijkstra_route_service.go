package service

import (
	"log"
	"math"

	"ev_routing/internal/dto"
)

// DijkstraRouteService searches an adjacency matrix (as built by
// RoutingService.GetAdjacencyMatrix) for the cheapest start-to-finish path
// using Dijkstra's algorithm.
type DijkstraRouteService struct{}

// NewDijkstraRouteService builds a DijkstraRouteService.
func NewDijkstraRouteService() *DijkstraRouteService {
	return &DijkstraRouteService{}
}

// GetRouteWithDijkstra finalizes, each round, the true minimum-cost
// unvisited node and relaxes its outgoing edges, until every reachable node
// has been processed, then walks the predecessor chain back from the finish
// node to recover the cheapest path.
func (s *DijkstraRouteService) GetRouteWithDijkstra(
	adjacencyMatrix map[int]map[int]dto.Edge,
	routeRequest *dto.RouteRequestDTO,
) []dto.RouteNodeDTO {
	routeMap := make(map[int]float64, len(adjacencyMatrix))
	connectMap := make(map[int]*int, len(adjacencyMatrix))
	queue := make([]int, 0, len(adjacencyMatrix))

	for nodeId := range adjacencyMatrix {
		if nodeId == startStationId {
			routeMap[nodeId] = 0
		} else {
			routeMap[nodeId] = math.MaxFloat64
		}
		connectMap[nodeId] = nil
		queue = append(queue, nodeId)
	}

	for len(queue) > 0 {
		log.Printf("Nodes left: %d", len(queue))

		minPos := 0
		for pos := 1; pos < len(queue); pos++ {
			if routeMap[queue[pos]] < routeMap[queue[minPos]] {
				minPos = pos
			}
		}

		currentNodeId := queue[minPos]
		queue = append(queue[:minPos], queue[minPos+1:]...)

		for nextNodeId, edge := range adjacencyMatrix[currentNodeId] {
			routeCost := edge.Cost + routeMap[currentNodeId]
			if routeMap[nextNodeId] > routeCost {
				routeMap[nextNodeId] = routeCost
				predecessorId := currentNodeId
				connectMap[nextNodeId] = &predecessorId
			}
		}
	}

	pathIds := []int{finishStationId}
	for {
		predecessor := connectMap[pathIds[len(pathIds)-1]]
		if predecessor == nil {
			break
		}
		pathIds = append(pathIds, *predecessor)
	}
	reverseInts(pathIds)

	return getRouteNodesFromIds(adjacencyMatrix, pathIds, routeRequest)
}

// reverseInts reverses ids in place.
func reverseInts(ids []int) {
	for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
		ids[i], ids[j] = ids[j], ids[i]
	}
}
