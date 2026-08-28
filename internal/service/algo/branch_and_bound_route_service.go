package algo

import (
	"log"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"ev_routing/internal/dto"
)

// bnbMaxExpansions caps the number of partial paths explored, guarding
// against combinatorial blow-up on large candidate pools.
const bnbMaxExpansions = 100000

// branchAndBoundCandidate is one feasible next hop considered while
// branching from a partial path.
type branchAndBoundCandidate struct {
	id   int
	edge dto.Edge
}

// BranchAndBoundRouteService does a DFS over partial paths, pruning any
// branch whose cost-so-far >= best found; exact since costs are non-negative.
type BranchAndBoundRouteService struct{}

// NewBranchAndBoundRouteService builds a BranchAndBoundRouteService.
func NewBranchAndBoundRouteService() *BranchAndBoundRouteService {
	return &BranchAndBoundRouteService{}
}

// GetRouteWithBranchAndBound explores partial paths depth-first, branching
// cheapest-edge-first, and returns the cheapest path to finish.
func (s *BranchAndBoundRouteService) GetRouteWithBranchAndBound(
	adjacencyMatrix map[int]map[int]dto.Edge,
	routeRequest *dto.RouteRequestDTO,
) []dto.RouteNodeDTO {
	start := time.Now()
	defer func() { log.Printf("Branch and bound duration: %v", time.Since(start)) }()

	bestCost := math.MaxFloat64
	var bestPath []int
	var mu sync.Mutex
	var expansions atomic.Int64

	rootVisited := map[int]bool{startStationId: true}
	rootCandidates := candidatesFrom(adjacencyMatrix, startStationId, rootVisited)

	// Each root candidate starts its own DFS branch with an independent
	// visited/path, so branches run concurrently; only the shared best and
	// expansion counter need synchronization.
	parallelFor(len(rootCandidates), func(idx int) {
		candidate := rootCandidates[idx]

		visited := make(map[int]bool, len(rootVisited)+1)
		for id, v := range rootVisited {
			visited[id] = v
		}
		visited[candidate.id] = true

		branchAndBound(
			adjacencyMatrix,
			candidate.id,
			visited,
			[]int{startStationId, candidate.id},
			candidate.edge.Cost,
			&mu,
			&bestCost,
			&bestPath,
			&expansions,
		)
	})

	log.Printf("Branch and bound explored %d nodes; best cost: %v", expansions.Load(), bestCost)

	if bestPath == nil {
		return []dto.RouteNodeDTO{}
	}

	return getRouteNodesFromIds(adjacencyMatrix, bestPath, routeRequest)
}

// branchAndBound extends path from currentId, updates best on a cheaper
// finish, and prunes branches whose cost so far is already >= *bestCost.
// bestCost/bestPath are shared across concurrently running root branches,
// so every access to them goes through mu; expansions is a plain atomic
// counter, allowed to briefly overshoot bnbMaxExpansions under
// concurrency since it's only a soft exploration cap.
func branchAndBound(
	adjacencyMatrix map[int]map[int]dto.Edge,
	currentId int,
	visited map[int]bool,
	path []int,
	costSoFar float64,
	mu *sync.Mutex,
	bestCost *float64,
	bestPath *[]int,
	expansions *atomic.Int64,
) {
	if expansions.Add(1) > bnbMaxExpansions {
		return
	}

	mu.Lock()
	prune := costSoFar >= *bestCost
	mu.Unlock()
	if prune {
		return
	}

	if currentId == finishStationId {
		mu.Lock()
		if costSoFar < *bestCost {
			*bestCost = costSoFar
			*bestPath = append([]int(nil), path...)
		}
		mu.Unlock()

		return
	}

	candidates := candidatesFrom(adjacencyMatrix, currentId, visited)

	for _, candidate := range candidates {
		visited[candidate.id] = true
		branchAndBound(
			adjacencyMatrix,
			candidate.id,
			visited,
			append(path, candidate.id),
			costSoFar+candidate.edge.Cost,
			mu,
			bestCost,
			bestPath,
			expansions,
		)
		visited[candidate.id] = false
	}
}

// candidatesFrom lists currentId's unvisited neighbors, cheapest-edge-first.
func candidatesFrom(
	adjacencyMatrix map[int]map[int]dto.Edge,
	currentId int,
	visited map[int]bool,
) []branchAndBoundCandidate {
	candidates := make([]branchAndBoundCandidate, 0, len(adjacencyMatrix[currentId]))
	for nextId, edge := range adjacencyMatrix[currentId] {
		if !visited[nextId] {
			candidates = append(candidates, branchAndBoundCandidate{id: nextId, edge: edge})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].edge.Cost < candidates[j].edge.Cost
	})

	return candidates
}
