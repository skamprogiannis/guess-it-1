package main

import (
	"math"
	"slices"
)

// stats contains rounded descriptive statistics for a set of integers.
type stats struct {
	average           int
	median            int
	variance          int
	standardDeviation int
}

// calculateStats returns rounded average, median, variance, and standard deviation.
func calculateStats(nums []int) stats {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	average := float64(sum) / float64(len(nums))

	sorted := slices.Clone(nums)
	slices.Sort(sorted)
	mid := len(sorted) / 2
	median := float64(sorted[mid])
	if len(sorted)%2 == 0 {
		median = float64(sorted[mid-1]+sorted[mid]) / 2
	}

	squaredDeviationSum := 0.0
	for _, n := range nums {
		deviation := float64(n) - average
		squaredDeviationSum += deviation * deviation
	}
	variance := squaredDeviationSum / float64(len(nums))

	return stats{
		average:           int(math.Round(average)),
		median:            int(math.Round(median)),
		variance:          int(math.Round(variance)),
		standardDeviation: int(math.Round(math.Sqrt(variance))),
	}
}

