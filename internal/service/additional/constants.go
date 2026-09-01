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

// MaxGraphNodes = MaxCandidateStations candidates + start/finish vertices.
// Largest possible V = {S, D} ∪ C_k.
const MaxGraphNodes = MaxCandidateStations + 2

// Shared iteration cap for all five search strategies (relaxation rounds, expansions, generations, shakes, ants).
// Sized off MaxGraphNodes so none is tuned more generously than another.
const MaxIterations = MaxGraphNodes * MaxGraphNodes

const (
	// CandidateBufferKm is the Buffer() width, km, used by
	// FilterCandidateStations.
	CandidateBufferKm = 5.0

	// Bounds concurrent OpenRouteService calls while building the adjacency matrix.
	// Higher than ParallelWorkers() since these are network-bound, not CPU-bound.
	AdjacencyMatrixConcurrency = 4

	// Queue size above which Dijkstra's minimum scan is worth splitting across goroutines.
	// Below it, goroutine-spawn overhead outweighs the win.
	DijkstraParallelThreshold = 64

	// GeneticSuccessTarget is how many valid DNAs a genetic-algorithm
	// island collects before it stops and returns its cheapest one.
	GeneticSuccessTarget = MaxGraphNodes * 3

	// Minimum generation history a GA island needs before crossover can run.
	// Below it, crossover degrades to mutation.
	GeneticMinParentsForCrossover = 3

	// VNSMaxNoImprove stops VNS after this many shake rounds in a row
	// without a better solution.
	VNSMaxNoImprove = MaxGraphNodes * 6

	// VNSKMax is the largest VNS shake-neighborhood size.
	VNSKMax = 3

	// Caps VNS's local-descent rounds (each flips one interior gene).
	// Past MaxGraphNodes rounds without converging, it isn't going to.
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
