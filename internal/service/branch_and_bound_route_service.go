package service

import (
	"log"
	"math"
	"sort"
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
	expansions := 0

	visited := map[int]bool{startStationId: true}
	branchAndBound(
		adjacencyMatrix,
		startStationId,
		visited,
		[]int{startStationId},
		0,
		&bestCost,
		&bestPath,
		&expansions,
	)

	log.Printf("Branch and bound explored %d nodes; best cost: %v", expansions, bestCost)

	if bestPath == nil {
		return []dto.RouteNodeDTO{}
	}

	return getRouteNodesFromIds(adjacencyMatrix, bestPath, routeRequest)
}

// branchAndBound extends path from currentId, updates best on a cheaper
// finish, and prunes branches whose cost so far is already >= *bestCost.
func branchAndBound(
	adjacencyMatrix map[int]map[int]dto.Edge,
	currentId int,
	visited map[int]bool,
	path []int,
	costSoFar float64,
	bestCost *float64,
	bestPath *[]int,
	expansions *int,
) {
	if *expansions >= bnbMaxExpansions {
		return
	}
	*expansions++

	if costSoFar >= *bestCost {
		return
	}

	if currentId == finishStationId {
		*bestCost = costSoFar
		*bestPath = append([]int(nil), path...)

		return
	}

	candidates := make([]branchAndBoundCandidate, 0, len(adjacencyMatrix[currentId]))
	for nextId, edge := range adjacencyMatrix[currentId] {
		if !visited[nextId] {
			candidates = append(candidates, branchAndBoundCandidate{id: nextId, edge: edge})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].edge.Cost < candidates[j].edge.Cost
	})

	for _, candidate := range candidates {
		visited[candidate.id] = true
		branchAndBound(
			adjacencyMatrix,
			candidate.id,
			visited,
			append(path, candidate.id),
			costSoFar+candidate.edge.Cost,
			bestCost,
			bestPath,
			expansions,
		)
		visited[candidate.id] = false
	}
}
