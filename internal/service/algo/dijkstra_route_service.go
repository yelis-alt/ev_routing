package algo

import (
	"log"
	"math"
	"sync"
	"time"

	"ev_routing/internal/dto"
)

// dijkstraParallelThreshold is the queue size above which the per-round
// minimum scan is worth splitting across goroutines; below it, the
// goroutine-spawn overhead outweighs the win.
const dijkstraParallelThreshold = 64

// DijkstraRouteService searches an adjacency matrix (from
// RoutingService.GetAdjacencyMatrix) for the cheapest path via Dijkstra.
type DijkstraRouteService struct{}

// NewDijkstraRouteService builds a DijkstraRouteService.
func NewDijkstraRouteService() *DijkstraRouteService {
	return &DijkstraRouteService{}
}

// GetRouteWithDijkstra relaxes the true minimum-cost unvisited node each
// round, then walks the predecessor chain back from finish to recover the path.
func (s *DijkstraRouteService) GetRouteWithDijkstra(
	adjacencyMatrix map[int]map[int]dto.Edge,
	routeRequest *dto.RouteRequestDTO,
) []dto.RouteNodeDTO {
	start := time.Now()
	defer func() { log.Printf("Dijkstra duration: %v", time.Since(start)) }()

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

		minPos := findMinPos(queue, routeMap)
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

	// Backtrack always seeds finish; unreachable means it never traces back to start.
	if len(pathIds) == 0 || pathIds[0] != startStationId {
		return []dto.RouteNodeDTO{}
	}

	return getRouteNodesFromIds(adjacencyMatrix, pathIds, routeRequest)
}

// findMinPos returns the index in queue whose routeMap value is smallest,
// preferring the first occurrence on ties (matching a plain left-to-right
// scan). Above dijkstraParallelThreshold, the scan is split into chunks
// evaluated concurrently and combined in order, which preserves that same
// tie-break.
func findMinPos(queue []int, routeMap map[int]float64) int {
	n := len(queue)
	if n < dijkstraParallelThreshold {
		return scanMinPos(queue, routeMap, 0, n)
	}

	workers := parallelWorkers()
	if workers > n {
		workers = n
	}
	chunkSize := (n + workers - 1) / workers

	localMins := make([]int, workers)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		lo := w * chunkSize
		hi := min(lo+chunkSize, n)
		if lo >= hi {
			localMins[w] = -1
			continue
		}

		wg.Add(1)
		go func(w, lo, hi int) {
			defer wg.Done()
			localMins[w] = scanMinPos(queue, routeMap, lo, hi)
		}(w, lo, hi)
	}
	wg.Wait()

	minPos := -1
	for _, pos := range localMins {
		if pos == -1 {
			continue
		}
		if minPos == -1 || routeMap[queue[pos]] < routeMap[queue[minPos]] {
			minPos = pos
		}
	}

	return minPos
}

// scanMinPos is a plain left-to-right minimum scan over queue[lo:hi].
func scanMinPos(queue []int, routeMap map[int]float64, lo, hi int) int {
	minPos := lo
	for pos := lo + 1; pos < hi; pos++ {
		if routeMap[queue[pos]] < routeMap[queue[minPos]] {
			minPos = pos
		}
	}

	return minPos
}

// reverseInts reverses ids in place.
func reverseInts(ids []int) {
	for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
		ids[i], ids[j] = ids[j], ids[i]
	}
}
