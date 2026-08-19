package service

import (
	"log"
	"math"
	"math/rand"
	"sort"

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
	nodeIds := make(map[int]struct{})
	for nodeId, neighbors := range adjacencyMatrix {
		nodeIds[nodeId] = struct{}{}
		for neighborId := range neighbors {
			nodeIds[neighborId] = struct{}{}
		}
	}

	parent := dto.GenerationDTO{
		Dna:     make(map[int]int, len(nodeIds)),
		Parents: make([]map[int]int, 0),
	}
	for id := range nodeIds {
		if id == startStationId || id == finishStationId {
			parent.Dna[id] = 1
		} else {
			parent.Dna[id] = 0
		}
	}

	geneIds := make([]int, 0, len(parent.Dna))
	for id := range parent.Dna {
		geneIds = append(geneIds, id)
	}
	sort.Ints(geneIds)

	rep := 0
	success := 0
	minCost := math.MaxFloat64
	var bestPathIds []int

	for success < genSuccessRep {
		rep++
		log.Printf("Generation %d out of %d; Successful DNAs: %d", rep, genRep, success)

		if rep > genRep && success == 0 {
			return []dto.RouteNodeDTO{}
		}

		var child dto.GenerationDTO
		switch rand.Intn(4) {
		case 0:
			child = s.GetMutation(parent, geneIds)
		case 1:
			child = s.GetCrossover(parent, geneIds)
		case 2:
			child = s.GetMutation(parent, geneIds)
			child = s.GetCrossover(child, geneIds)
		case 3:
			child = s.GetCrossover(parent, geneIds)
			child = s.GetMutation(child, geneIds)
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

	log.Printf("Generation %d out of %d; Successful DNAs: %d", rep, genRep, success)

	return getRouteNodesFromIds(adjacencyMatrix, bestPathIds, routeRequest)
}

// GetMutation flips a single, randomly-chosen non-endpoint gene.
func (s *GeneticRouteService) GetMutation(generation dto.GenerationDTO, geneIds []int) dto.GenerationDTO {
	randomGeneIndex := rand.Intn(len(geneIds)-2) + 1
	mutatedGeneId := geneIds[randomGeneIndex]
	generation.Dna[mutatedGeneId] = 1 - generation.Dna[mutatedGeneId]

	return generation
}

// GetCrossover combines two prior generations' genes; falls back to a
// mutation until at least minParentsForCrossover prior generations exist.
func (s *GeneticRouteService) GetCrossover(generation dto.GenerationDTO, geneIds []int) dto.GenerationDTO {
	parentDnas := generation.Parents
	if len(parentDnas) < minParentsForCrossover {
		mutant := s.GetMutation(generation, geneIds)
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
