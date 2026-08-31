package additional

import (
	"math"
	"time"
)

const (
	// StartStationId is the synthetic vertex id for the trip's start point.
	StartStationId = math.MinInt32
	// FinishStationId is the synthetic vertex id for the trip's finish point.
	FinishStationId = math.MaxInt32
)

// MaxCandidateStations is k's coefficient (30 in k = 30(1-...)); also the
// upper bound on |C_k|, the candidate-station pool built by
// FilterCandidateStations.
const MaxCandidateStations = 30

// MaxGraphNodes is the largest the routing graph V = {S, D} ∪ C_k can ever
// be: MaxCandidateStations candidates plus the two synthetic start/finish
// vertices.
const MaxGraphNodes = MaxCandidateStations + 2

// MaxIterations is the shared iteration cap used by every search strategy in
// internal/service/algo (Dijkstra's relaxation rounds, Branch & Bound's
// expansion count, the genetic algorithm's per-island generation attempts,
// VNS's shake rounds, and ACO's iteration count) — sized off MaxGraphNodes
// so it scales correctly if that cap ever changes, and kept identical across
// all five so none is tuned more generously than another.
const MaxIterations = MaxGraphNodes * MaxGraphNodes

const (
	// CandidateBufferKm is the Buffer() width, km, used by
	// FilterCandidateStations.
	CandidateBufferKm = 5.0

	// AdjacencyMatrixConcurrency bounds concurrent OpenRouteService calls
	// while building the adjacency matrix; higher than ParallelWorkers()
	// since these are network-bound, not CPU-bound.
	AdjacencyMatrixConcurrency = 4

	// DijkstraParallelThreshold is the queue size above which Dijkstra's
	// per-round minimum scan is worth splitting across goroutines; below
	// it, the goroutine-spawn overhead outweighs the win.
	DijkstraParallelThreshold = 64

	// GeneticSuccessTarget is how many valid DNAs a genetic-algorithm
	// island collects before it stops and returns its cheapest one.
	GeneticSuccessTarget = MaxGraphNodes * 3

	// GeneticMinParentsForCrossover is the minimum generation history a
	// genetic-algorithm island needs before crossover can run; below it,
	// crossover degrades to mutation.
	GeneticMinParentsForCrossover = 3

	// VNSMaxNoImprove stops VNS after this many shake rounds in a row
	// without a better solution.
	VNSMaxNoImprove = MaxGraphNodes * 6

	// VNSKMax is the largest VNS shake-neighborhood size.
	VNSKMax = 3

	// VNSLocalSearchMaxSteps caps VNS's local-descent rounds; each round
	// flips one interior gene, so more than MaxGraphNodes rounds without
	// converging means the descent isn't going to converge.
	VNSLocalSearchMaxSteps = MaxGraphNodes

	// ACOAntCount is the number of ants walked per ACO iteration; one per
	// possible candidate station.
	ACOAntCount = MaxGraphNodes - 2

	// ACOAlpha is ACO's pheromone-influence exponent.
	ACOAlpha = 1.0
	// ACOBeta is ACO's heuristic (1/cost) influence exponent.
	ACOBeta = 3.0
	// ACOEvaporationRate is ACO's pheromone evaporation rate, rho.
	ACOEvaporationRate = 0.3
	// ACOInitialPheromone is every edge's pheromone level before the first
	// iteration.
	ACOInitialPheromone = 1.0
	// ACOPheromoneDeposit is Q, the numerator of each ant's per-edge
	// pheromone deposit.
	ACOPheromoneDeposit = 100.0
	// ACOCostEpsilon avoids division by zero when weighing zero-cost edges.
	ACOCostEpsilon = 1e-6
)

// HoursToDuration converts fractional hours (as returned by
// OpenRouteService) into a time.Duration.
func HoursToDuration(hours float64) time.Duration {
	return time.Duration(hours * float64(time.Hour))
}
