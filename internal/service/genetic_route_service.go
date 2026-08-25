package service

import (
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"ev_routing/internal/dto"
)

const (
	genRep                 = 1000
	genSuccessRep          = 100
	minParentsForCrossover = 3
)

// GeneticRouteService searches an adjacency matrix (from
// RoutingService.GetAdjacencyMatrix) for the cheapest path via a GA.
type GeneticRouteService struct{}

// NewGeneticRouteService builds a GeneticRouteService.
func NewGeneticRouteService() *GeneticRouteService {
	return &GeneticRouteService{}
}

// GetRouteWithEvolution encodes each path as a DNA map (node id ->
// included) and evolves it via mutation/crossover, keeping the cheapest valid one.
func (s *GeneticRouteService) GetRouteWithEvolution(
	adjacencyMatrix map[int]map[int]dto.Edge,
	routeRequest *dto.RouteRequestDTO,
) []dto.RouteNodeDTO {
	start := time.Now()
	defer func() { log.Printf("Genetic duration: %v", time.Since(start)) }()

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

	// Each island runs the hill-climb independently (its own parent/child
	// lineage, unshared across goroutines), so they run in parallel; the
	// cheapest path across all islands wins.
	islands := parallelWorkers()
	islandPathIds := make([][]int, islands)
	islandCosts := make([]float64, islands)

	parallelFor(islands, func(i int) {
		islandPathIds[i], islandCosts[i] = s.runIsland(i, adjacencyMatrix, geneIds)
	})

	bestIsland := -1
	minCost := math.MaxFloat64
	for i, cost := range islandCosts {
		if islandPathIds[i] != nil && cost < minCost {
			minCost = cost
			bestIsland = i
		}
	}

	if bestIsland == -1 {
		return []dto.RouteNodeDTO{}
	}

	return getRouteNodesFromIds(adjacencyMatrix, islandPathIds[bestIsland], routeRequest)
}

// runIsland runs one independent instance of the mutation/crossover
// hill-climb described on GetRouteWithEvolution, returning the cheapest
// valid path it found (nil, +Inf if it never found one).
func (s *GeneticRouteService) runIsland(
	island int,
	adjacencyMatrix map[int]map[int]dto.Edge,
	geneIds []int,
) ([]int, float64) {
	parent := dto.GenerationDTO{
		Dna:     make(map[int]int, len(geneIds)),
		Parents: make([]map[int]int, 0),
	}
	for _, id := range geneIds {
		if id == startStationId || id == finishStationId {
			parent.Dna[id] = 1
		} else {
			parent.Dna[id] = 0
		}
	}

	rep := 0
	success := 0
	minCost := math.MaxFloat64
	var bestPathIds []int

	for success < genSuccessRep {
		rep++
		log.Printf("Island %d: generation %d out of %d; Successful DNAs: %d", island, rep, genRep, success)

		if rep > genRep && success == 0 {
			return nil, math.MaxFloat64
		}

		var child dto.GenerationDTO
		switch rand.Intn(4) {
		case 0:
			child = s.getMutation(parent, geneIds)
		case 1:
			child = s.getCrossover(parent, geneIds)
		case 2:
			child = s.getMutation(parent, geneIds)
			child = s.getCrossover(child, geneIds)
		case 3:
			child = s.getCrossover(parent, geneIds)
			child = s.getMutation(child, geneIds)
		}

		pathIds := make([]int, 0)
		for geneId, included := range child.Dna {
			if included == 1 {
				pathIds = append(pathIds, geneId)
			}
		}
		sort.Ints(pathIds)

		routeCost := 0.0
		failed := false
		for i := 0; i < len(pathIds)-1; i++ {
			currentNodeId := pathIds[i]
			nextNodeId := pathIds[i+1]

			neighbors, ok := adjacencyMatrix[currentNodeId]
			if !ok {
				failed = true
				break
			}

			edge, ok := neighbors[nextNodeId]
			if !ok {
				failed = true
				break
			}

			routeCost += edge.Cost
		}

		if !failed {
			success++

			if routeCost < minCost {
				minCost = routeCost
				bestPathIds = pathIds
			}
		}

		parent = child
	}

	log.Printf("Island %d: generation %d out of %d; Successful DNAs: %d", island, rep, genRep, success)

	return bestPathIds, minCost
}

// getMutation flips a single, randomly-chosen non-endpoint gene.
func (s *GeneticRouteService) getMutation(generation dto.GenerationDTO, geneIds []int) dto.GenerationDTO {
	randomGeneIndex := rand.Intn(len(geneIds)-2) + 1
	mutatedGeneId := geneIds[randomGeneIndex]
	generation.Dna[mutatedGeneId] = 1 - generation.Dna[mutatedGeneId]

	return generation
}

// getCrossover combines two prior generations' genes; falls back to a
// mutation until at least minParentsForCrossover prior generations exist.
func (s *GeneticRouteService) getCrossover(generation dto.GenerationDTO, geneIds []int) dto.GenerationDTO {
	parentDnas := generation.Parents
	if len(parentDnas) < minParentsForCrossover {
		mutant := s.getMutation(generation, geneIds)
		mutant.Parents = append(mutant.Parents, mutant.Dna)

		return mutant
	}

	parent1 := parentDnas[rand.Intn(len(parentDnas)-1)]
	parent2 := parentDnas[rand.Intn(len(parentDnas)-1)]

	for _, geneId := range geneIds[1 : len(geneIds)-1] {
		childGene := parent1[geneId] + parent2[geneId]
		if childGene == 2 {
			generation.Dna[geneId] = rand.Intn(2)
		} else {
			generation.Dna[geneId] = childGene
		}
	}

	generation.Parents = append(generation.Parents, generation.Dna)

	return generation
}
