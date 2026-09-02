package main

import (
	"reflect"
	"testing"
)

func TestCalculateStats(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want stats
	}{
		{
			name: "even positive sample",
			nums: []int{1, 2, 3, 4},
			want: stats{
				average:           3,
				median:            3,
				variance:          1,
				standardDeviation: 1,
				firstQuartile:     2,
				thirdQuartile:     4,
			},
		},
		{
			name: "symmetric signed sample",
			nums: []int{-4, -2, 0, 2, 4},
			want: stats{
				average:           0,
				median:            0,
				variance:          8,
				standardDeviation: 3,
				firstQuartile:     -2,
				thirdQuartile:     2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := append([]int(nil), tt.nums...)

			got := calculateStats(tt.nums)

			if got != tt.want {
				t.Fatalf("calculateStats(%v) = %+v, want %+v", tt.nums, got, tt.want)
			}
			if !reflect.DeepEqual(tt.nums, original) {
				t.Fatalf("calculateStats mutated input: got %v, want %v", tt.nums, original)
			}
		})
	}
}
