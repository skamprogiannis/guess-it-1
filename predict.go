package main

func PredictRange(nums []int) (int, int) {
	if len(nums) == 0 {
		return 0, 100
	}

	window := lastValues(nums, 12)
	summary := CalculateSummary(window)
	width := summary.StandardDeviation * 2

	if width < 25 {
		width = 25
	}

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
