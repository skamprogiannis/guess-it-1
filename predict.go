package main

const (
	historyWindowSize         = 12
	standardDeviationMultiple = 2
	minimumRangeWidth         = 25
)

func PredictRange(nums []int) (int, int) {
	if len(nums) == 0 {
		return 0, 100
	}

	window := lastValues(nums, historyWindowSize)
	summary := CalculateSummary(window)
	width := max(summary.StandardDeviation * standardDeviationMultiple, minimumRangeWidth)
	prediction := summary.Average

	if len(window) >= 2 {
		prediction += window[len(window)-1] - window[len(window)-2]
	}

	return prediction - width, prediction + width
}

func lastValues(nums []int, limit int) []int {
	if len(nums) <= limit {
		return nums
	}
	return nums[len(nums)-limit:]
}
