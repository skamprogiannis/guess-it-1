package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
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

// readDataLines reads non-empty trimmed lines from path.
func readDataLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// parseInput converts text lines to integers and returns warnings for skipped lines.
func parseInput(lines []string) ([]int, []string, error) {
	nums := make([]int, 0, len(lines))
	parseWarnings := make([]string, 0)
	for _, line := range lines {
		val, err := strconv.Atoi(line)
		if err != nil {
			parseWarnings = append(parseWarnings, fmt.Sprintf("Skipping invalid number: %s", line))
			continue
		}
		nums = append(nums, val)
	}

	if len(nums) == 0 {
		return nil, parseWarnings, fmt.Errorf("no valid numbers provided")
	}

	return nums, parseWarnings, nil
}
