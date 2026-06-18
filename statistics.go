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

// Summary contains rounded descriptive statistics for a set of integers.
type Summary struct {
	Average           int
	Median            int
	Variance          int
	StandardDeviation int
}

// CalculateSummary returns rounded average, median, variance, and standard deviation.
func CalculateSummary(nums []int) Summary {
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

	return Summary{
		Average:           int(math.Round(average)),
		Median:            int(math.Round(median)),
		Variance:          int(math.Round(variance)),
		StandardDeviation: int(math.Round(math.Sqrt(variance))),
	}
}

// CalculateAverage returns the rounded average of nums.
func CalculateAverage(nums []int) int {
	return CalculateSummary(nums).Average
}

// CalculateMedian returns the rounded median of nums.
func CalculateMedian(nums []int) int {
	return CalculateSummary(nums).Median
}

// CalculateVariance returns the rounded variance of nums.
func CalculateVariance(nums []int) int {
	return CalculateSummary(nums).Variance
}

// CalculateStandardDeviation returns the rounded standard deviation of nums.
func CalculateStandardDeviation(nums []int) int {
	return CalculateSummary(nums).StandardDeviation
}

// ReadDataLines reads non-empty trimmed lines from path.
func ReadDataLines(path string) ([]string, error) {
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

// ParseInput converts text lines to integers and returns warnings for skipped lines.
func ParseInput(lines []string) ([]int, []string, error) {
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
