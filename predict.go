package main

import "math"

const (
	deviationScale      = 4
	wideSideScale       = deviationScale + 1
	minimumDeviation    = wideSideScale * wideSideScale
	robustStatsWindow   = minimumDeviation*wideSideScale + wideSideScale
	rangeScoreNumerator = 10000000
	scoreEventWindow    = minimumDeviation * 2
	scoreSwitchMargin   = rangeScoreNumerator / (deviationScale*2 + 1)
)

// band identifies the quartile region containing a value.
type band int

const (
	lowerBand band = iota
	lowerMiddleBand
	upperMiddleBand
	upperBand
)

type predictionRange struct {
	low  int
	high int
}

type scoreState struct {
	pendingRobust   predictionRange
	pendingQuartile predictionRange
	hasPending      bool
	robustEvents    []int
	quartileEvents  []int
}

var predictorState scoreState

// PredictRange returns lower and upper bounds for the next likely input value.
func PredictRange(nums []int) (int, int) {
	last := nums[len(nums)-1]
	predictorState.record(last)

	robust := robustRange(nums)
	quartile := quartileRange(nums)
	predictorState.remember(robust, quartile)

	selected := predictorState.selectRange(robust, quartile)
	return selected.low, selected.high
}

// record scores the previous ranges now that the actual next value is known.
func (state *scoreState) record(actual int) {
	if !state.hasPending {
		return
	}

	robustScore := scoreRange(state.pendingRobust, actual)
	quartileScore := scoreRange(state.pendingQuartile, actual)
	if robustScore == 0 && quartileScore == 0 {
		return
	}

	state.robustEvents = append(state.robustEvents, robustScore)
	state.quartileEvents = append(state.quartileEvents, quartileScore)
}

// remember stores the current expert ranges so the next input can score them.
func (state *scoreState) remember(robust predictionRange, quartile predictionRange) {
	state.pendingRobust = robust
	state.pendingQuartile = quartile
	state.hasPending = true
}

// selectRange chooses the stronger recent expert, with a bias toward robustRange.
func (state *scoreState) selectRange(robust predictionRange, quartile predictionRange) predictionRange {
	robustScore := recentScore(state.robustEvents)
	quartileScore := recentScore(state.quartileEvents)
	if quartileScore > robustScore+scoreSwitchMargin {
		return quartile
	}
	return robust
}

// recentScore sums the latest informative score events.
func recentScore(scores []int) int {
	start := len(scores) - scoreEventWindow
	if start < 0 {
		start = 0
	}

	total := 0
	for _, score := range scores[start:] {
		total += score
	}
	return total
}

// scoreRange rewards a correct range more when the range is narrow.
func scoreRange(candidate predictionRange, actual int) int {
	if actual < candidate.low || actual > candidate.high {
		return 0
	}

	width := candidate.high - candidate.low + 1
	return rangeScoreNumerator / width
}

// robustRange predicts from stable statistics after filtering large outliers.
func robustRange(nums []int) predictionRange {
	last := nums[len(nums)-1]
	stats := stableStats(nums)
	spread := stats.standardDeviation
	if spread < 1 {
		spread = 1
	}

	outlierDistance := spread * deviationScale
	if outlierDistance < minimumDeviation {
		outlierDistance = minimumDeviation
	}
	if absolute(last-stats.median) > outlierDistance {
		return predictionRange{
			low:  stats.average - deviationScale,
			high: stats.average + deviationScale,
		}
	}

	lowerBand := stats.median - spread - spread/wideSideScale
	upperBand := stats.median + spread
	switch {
	case last < lowerBand:
		return predictionRange{low: last, high: last}
	case last < stats.median:
		return predictionRange{low: last - 1, high: last - 1}
	case last < upperBand:
		return predictionRange{low: last, high: last}
	default:
		return predictionRange{
			low:  last - wideSideScale,
			high: last - deviationScale,
		}
	}
}

// stableStats returns statistics from a recent sample after trimming outliers.
func stableStats(nums []int) stats {
	sample := recentValues(nums, robustStatsWindow)
	base := calculateStats(sample)

	deviations := make([]int, len(sample))
	for i, n := range sample {
		deviations[i] = absolute(n - base.median)
	}
	deviationStats := calculateStats(deviations)
	deviationLimit := deviationStats.median * deviationScale
	if deviationLimit < minimumDeviation {
		deviationLimit = minimumDeviation
	}

	core := make([]int, 0, len(sample))
	for _, n := range sample {
		if absolute(n-base.median) <= deviationLimit {
			core = append(core, n)
		}
	}
	if len(core) == 0 {
		core = sample
	}
	return calculateStats(core)
}

// recentValues returns at most limit values from the end of nums.
func recentValues(nums []int, limit int) []int {
	start := len(nums) - limit
	if start < 0 {
		start = 0
	}

	return nums[start:]
}

// quartileRange predicts from the current value's recent quartile band.
func quartileRange(nums []int) predictionRange {
	last := nums[len(nums)-1]
	stats := calculateStats(recentSample(nums))
	currentBand := classifyBand(last, stats)

	lowDelta, highDelta := selectDeltas(nums, currentBand)
	return predictionRange{low: last + lowDelta, high: last + highDelta}
}

// classifyBand maps a value to its quartile band using recent statistics.
func classifyBand(value int, stats stats) band {
	switch {
	case value < stats.firstQuartile:
		return lowerBand
	case value < stats.median:
		return lowerMiddleBand
	case value < stats.thirdQuartile:
		return upperMiddleBand
	default:
		return upperBand
	}
}

// selectDeltas returns the delta range to apply for a classified value.
func selectDeltas(nums []int, currentBand band) (int, int) {
	switch currentBand {
	case lowerBand:
		return 0, 0
	case lowerMiddleBand:
		return -1, -1
	case upperMiddleBand:
		return 0, 1
	case upperBand:
		return learnedHighPullback(nums)
	}
	return 0, 0
}

// learnedHighPullback derives a downward range from recent upper-band moves.
func learnedHighPullback(nums []int) (int, int) {
	sample := recentSample(nums)
	negativeDeltas := make([]int, 0, len(sample))

	for i := 1; i < len(sample); i++ {
		previous := sample[i-1]
		current := sample[i]
		previousStats := calculateStats(recentSample(sample[:i]))
		previousBand := classifyBand(previous, previousStats)
		if previousBand == upperBand && current < previous {
			negativeDeltas = append(negativeDeltas, current-previous)
		}
	}

	if len(negativeDeltas) == 0 {
		return 0, 0
	}

	stats := calculateStats(negativeDeltas)
	return stats.firstQuartile, stats.median
}

// recentSample returns a recency-weighted sample without a fixed window size.
func recentSample(nums []int) []int {
	sampleSize := int(math.Sqrt(float64(len(nums))))
	start := len(nums) - sampleSize
	return nums[start:]
}

// absolute returns the absolute value of n.
func absolute(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
