package algo

import (
	"log"
	"maps"
	"math"
	"math/rand"
	"sort"
	"time"

	"ev_routing/internal/dto"
	"ev_routing/internal/service/additional"
	"ev_routing/internal/service/geo"
)

const (
	vnsMaxIterations       = 1000
	vnsMaxNoImprove        = 200
	vnsKMax                = 3
	vnsLocalSearchMaxSteps = 50
)

// vnsSolution is one candidate path in the VNS search, encoded the same way
// as GeneticRouteService's DNA (node id -> included).
type vnsSolution struct {
	dna      map[int]int
	pathIds  []int
	cost     float64
	feasible bool
}

// VNSRouteService shakes the incumbent in neighborhood Nk, descends to a
// local optimum, and keeps it if better (resetting k=1) else grows k.
type VNSRouteService struct{}

// NewVNSRouteService builds a VNSRouteService.
func NewVNSRouteService() *VNSRouteService {
	return &VNSRouteService{}
}

// GetRouteWithVNS runs VNS until vnsMaxIterations shakes or vnsMaxNoImprove
// stalls, returning the cheapest feasible path found.
func (s *VNSRouteService) GetRouteWithVNS(
	adjacencyMatrix map[int]map[int]dto.Edge,
	routeRequest *dto.RouteRequestDTO,
) []dto.RouteNodeDTO {
	start := time.Now()
	defer func() { log.Printf("VNS duration: %v", time.Since(start)) }()

	nodeIds := make(map[int]struct{})
	for nodeId, neighbors := range adjacencyMatrix {
		nodeIds[nodeId] = struct{}{}
		for neighborId := range neighbors {
			nodeIds[neighborId] = struct{}{}
		}
	}

	geneIds := make([]int, 0, len(nodeIds))
	for id := range nodeIds {
		geneIds = append(geneIds, id)
	}
	sort.Ints(geneIds)

	baseDna := make(map[int]int, len(geneIds))
	for _, id := range geneIds {
		if id == additional.StartStationId || id == additional.FinishStationId {
			baseDna[id] = 1
		} else {
			baseDna[id] = 0
		}
	}

	best := vnsSolution{dna: baseDna, cost: math.MaxFloat64}

	iterations := 0
	noImprove := 0
	for iterations < vnsMaxIterations && noImprove < vnsMaxNoImprove {
		k := 1
		for k <= vnsKMax && iterations < vnsMaxIterations && noImprove < vnsMaxNoImprove {
			iterations++
			log.Printf(
				"Iteration %d out of %d; Neighborhood k=%d; Best cost: %v",
				iterations, vnsMaxIterations, k, best.cost,
			)

			shaken := shakeDna(best.dna, k, geneIds)
			candidate := localSearchDna(adjacencyMatrix, shaken, geneIds)

			if candidate.feasible && candidate.cost < best.cost {
				best = candidate
				noImprove = 0
				k = 1
			} else {
				noImprove++
				k++
			}
		}
	}

	log.Printf("Iterations run: %d; Best cost: %v; Feasible: %t", iterations, best.cost, best.feasible)

	if !best.feasible {
		return []dto.RouteNodeDTO{}
	}

	return geo.GetRouteNodesFromIds(adjacencyMatrix, best.pathIds, routeRequest)
}

// shakeDna copies dna and flips min(k, len(interior genes)) distinct,
// randomly-chosen non-endpoint genes, producing a point in Nk(dna).
func shakeDna(dna map[int]int, k int, geneIds []int) map[int]int {
	shaken := copyDna(dna)

	interiorIds := geneIds[1 : len(geneIds)-1]
	if len(interiorIds) == 0 {
		return shaken
	}

	flips := min(k, len(interiorIds))

	for _, pos := range rand.Perm(len(interiorIds))[:flips] {
		geneId := interiorIds[pos]
		shaken[geneId] = 1 - shaken[geneId]
	}

	return shaken
}

// localSearchDna descends to a local optimum under N1 (single-gene flips),
// best-improvement, capped at vnsLocalSearchMaxSteps rounds.
func localSearchDna(
	adjacencyMatrix map[int]map[int]dto.Edge,
	dna map[int]int,
	geneIds []int,
) vnsSolution {
	current := evaluateDna(adjacencyMatrix, dna, geneIds)
	interiorIds := geneIds[1 : len(geneIds)-1]

	for range vnsLocalSearchMaxSteps {
		// Each neighbor only reads current.dna, so all flips can be
		// evaluated concurrently; the best-improvement reduce below stays
		// sequential (in interiorIds order) to keep tie-breaks deterministic.
		neighbors := make([]vnsSolution, len(interiorIds))
		additional.ParallelFor(len(interiorIds), func(idx int) {
			neighborDna := copyDna(current.dna)
			geneId := interiorIds[idx]
			neighborDna[geneId] = 1 - neighborDna[geneId]

			neighbors[idx] = evaluateDna(adjacencyMatrix, neighborDna, geneIds)
		})

		best := current
		improved := false
		for _, neighbor := range neighbors {
			if neighbor.feasible && (!best.feasible || neighbor.cost < best.cost) {
				best = neighbor
				improved = true
			}
		}

		if !improved {
			break
		}
		current = best
	}

	return current
}

// evaluateDna sums edge costs along dna's included genes; a missing
// consecutive edge, or <2 nodes, means infeasible.
func evaluateDna(
	adjacencyMatrix map[int]map[int]dto.Edge,
	dna map[int]int,
	geneIds []int,
) vnsSolution {
	pathIds := make([]int, 0, len(dna))
	for _, id := range geneIds {
		if dna[id] == 1 {
			pathIds = append(pathIds, id)
		}
	}

	if len(pathIds) < 2 {
		return vnsSolution{dna: dna, pathIds: pathIds, feasible: false}
	}

	cost := 0.0
	for i := 0; i < len(pathIds)-1; i++ {
		edge, ok := adjacencyMatrix[pathIds[i]][pathIds[i+1]]
		if !ok {
			return vnsSolution{dna: dna, pathIds: pathIds, feasible: false}
		}
		cost += edge.Cost
	}

	return vnsSolution{dna: dna, pathIds: pathIds, cost: cost, feasible: true}
}

// copyDna returns a shallow copy of dna, so callers can mutate it without
// aliasing the source solution.
func copyDna(dna map[int]int) map[int]int {
	copied := make(map[int]int, len(dna))
	maps.Copy(copied, dna)

	return copied
}
