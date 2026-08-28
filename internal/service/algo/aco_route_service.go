package algo

import (
	"log"
	"math"
	"math/rand"
	"time"

	"ev_routing/internal/dto"
	"ev_routing/internal/service/additional"
	"ev_routing/internal/service/geo"
)

const (
	acoIterations       = 200
	acoAntCount         = 30
	acoAlpha            = 1.0 // pheromone influence
	acoBeta             = 3.0 // heuristic (1/cost) influence
	acoEvaporationRate  = 0.3 // rho
	acoInitialPheromone = 1.0
	acoPheromoneDeposit = 100.0 // Q, numerator of the per-ant deposit
	acoCostEpsilon      = 1e-6  // avoids div-by-zero for zero-cost edges
)

// acoCandidate is one feasible next hop considered while an ant is at a
// given node, paired with the edge it would take.
type acoCandidate struct {
	id   int
	edge dto.Edge
}

// acoAntResult is the path one ant completed in an iteration, kept for the
// end-of-iteration pheromone deposit.
type acoAntResult struct {
	pathIds []int
	cost    float64
}

// ACORouteService walks ants start->finish biased by pheromone^alpha *
// (1/cost)^beta, then evaporates and redeposits pheromone by path cost.
type ACORouteService struct{}

// NewACORouteService builds an ACORouteService.
func NewACORouteService() *ACORouteService {
	return &ACORouteService{}
}

// GetRouteWithACO runs acoIterations rounds of acoAntCount ants, updating
// pheromone each round, and returns the cheapest path any ant found.
func (s *ACORouteService) GetRouteWithACO(
	adjacencyMatrix map[int]map[int]dto.Edge,
	routeRequest *dto.RouteRequestDTO,
) []dto.RouteNodeDTO {
	start := time.Now()
	defer func() { log.Printf("ACO duration: %v", time.Since(start)) }()

	pheromone := make(map[int]map[int]float64, len(adjacencyMatrix))
	for from, neighbors := range adjacencyMatrix {
		pheromone[from] = make(map[int]float64, len(neighbors))
		for to := range neighbors {
			pheromone[from][to] = acoInitialPheromone
		}
	}

	maxSteps := len(adjacencyMatrix) + 1

	var bestPath []int
	bestCost := math.MaxFloat64

	for iteration := range acoIterations {
		// Ants only read pheromone this round (it's only mutated below,
		// after every ant has finished), so they can walk concurrently.
		antResults := make([]*acoAntResult, acoAntCount)
		additional.ParallelFor(acoAntCount, func(i int) {
			pathIds, cost, reached := walkAnt(adjacencyMatrix, pheromone, maxSteps)
			if reached {
				antResults[i] = &acoAntResult{pathIds: pathIds, cost: cost}
			}
		})

		results := make([]acoAntResult, 0, acoAntCount)
		for _, result := range antResults {
			if result == nil {
				continue
			}

			results = append(results, *result)
			if result.cost < bestCost {
				bestCost = result.cost
				bestPath = result.pathIds
			}
		}

		evaporatePheromone(pheromone)
		depositPheromone(pheromone, results)

		log.Printf(
			"Iteration %d out of %d; Ants reached finish: %d out of %d; Best cost: %v",
			iteration+1, acoIterations, len(results), acoAntCount, bestCost,
		)
	}

	if bestPath == nil {
		return []dto.RouteNodeDTO{}
	}

	return geo.GetRouteNodesFromIds(adjacencyMatrix, bestPath, routeRequest)
}

// walkAnt builds one path from start, picking each hop by pheromone^alpha *
// (1/cost)^beta; fails on a dead end or if maxSteps is exceeded.
func walkAnt(
	adjacencyMatrix map[int]map[int]dto.Edge,
	pheromone map[int]map[int]float64,
	maxSteps int,
) ([]int, float64, bool) {
	visited := map[int]bool{additional.StartStationId: true}
	path := []int{additional.StartStationId}
	currentId := additional.StartStationId
	cost := 0.0

	for range maxSteps {
		if currentId == additional.FinishStationId {
			return path, cost, true
		}

		candidates := make([]acoCandidate, 0, len(adjacencyMatrix[currentId]))
		weights := make([]float64, 0, len(adjacencyMatrix[currentId]))
		totalWeight := 0.0

		for nextId, edge := range adjacencyMatrix[currentId] {
			if visited[nextId] {
				continue
			}

			desirability := 1 / (edge.Cost + acoCostEpsilon)
			weight := math.Pow(pheromone[currentId][nextId], acoAlpha) * math.Pow(desirability, acoBeta)

			candidates = append(candidates, acoCandidate{id: nextId, edge: edge})
			weights = append(weights, weight)
			totalWeight += weight
		}

		if len(candidates) == 0 || totalWeight <= 0 {
			return nil, 0, false
		}

		next := candidates[selectWeightedIndex(weights, totalWeight)]

		visited[next.id] = true
		path = append(path, next.id)
		cost += next.edge.Cost
		currentId = next.id
	}

	return nil, 0, false
}

// selectWeightedIndex does roulette-wheel selection over weights, which
// must sum to totalWeight.
func selectWeightedIndex(weights []float64, totalWeight float64) int {
	r := rand.Float64() * totalWeight

	cumulative := 0.0
	for i, weight := range weights {
		cumulative += weight
		if r <= cumulative {
			return i
		}
	}

	return len(weights) - 1
}

// evaporatePheromone scales every edge's pheromone down by (1 -
// acoEvaporationRate), so trails not reinforced this iteration fade out.
func evaporatePheromone(pheromone map[int]map[int]float64) {
	for _, neighbors := range pheromone {
		for to := range neighbors {
			neighbors[to] *= 1 - acoEvaporationRate
		}
	}
}

// depositPheromone adds acoPheromoneDeposit/cost to every edge on each
// result's path, so cheaper completed paths reinforce their edges more.
func depositPheromone(pheromone map[int]map[int]float64, results []acoAntResult) {
	for _, result := range results {
		deposit := acoPheromoneDeposit / result.cost

		for i := 0; i < len(result.pathIds)-1; i++ {
			from, to := result.pathIds[i], result.pathIds[i+1]
			pheromone[from][to] += deposit
		}
	}
}
