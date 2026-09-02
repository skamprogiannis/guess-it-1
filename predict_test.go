package main

import "testing"

func TestPredictRangeKeepsStableSeriesExact(t *testing.T) {
	predictorState = scoreState{}
	var history []int

	for i := 0; i < 5; i++ {
		history = append(history, 42)
		low, high := PredictRange(history)
		if low != 42 || high != 42 {
			t.Fatalf("prediction %d = %d %d, want 42 42", i+1, low, high)
		}
	}
}

func TestPredictRangeFiltersOneLargeOutlier(t *testing.T) {
	predictorState = scoreState{}
	history := make([]int, robustStatsWindow)
	for i := range history {
		history[i] = 100
	}
	history[len(history)-1] = 1000

	low, high := PredictRange(history)

	if low != 96 || high != 104 {
		t.Fatalf("PredictRange(outlier history) = %d %d, want 96 104", low, high)
	}
}
