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

// HoursToDuration converts fractional hours (as returned by
// OpenRouteService) into a time.Duration.
func HoursToDuration(hours float64) time.Duration {
	return time.Duration(hours * float64(time.Hour))
}
