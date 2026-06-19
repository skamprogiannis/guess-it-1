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

// calculateAverage returns the rounded average of nums.
func calculateAverage(nums []int) int {
	return calculateStats(nums).average
}

// calculateMedian returns the rounded median of nums.
func calculateMedian(nums []int) int {
	return calculateStats(nums).median
}

// calculateVariance returns the rounded variance of nums.
func calculateVariance(nums []int) int {
	return calculateStats(nums).variance
}

// calculateStandardDeviation returns the rounded standard deviation of nums.
func calculateStandardDeviation(nums []int) int {
	return calculateStats(nums).standardDeviation
}

