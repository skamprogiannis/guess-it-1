# Strategy Notes

`guess-it-1` rewards narrow correct ranges much more than wide correct ranges.
The current strategy therefore makes small, high-value bets instead of trying to
cover most possible next values.

## Runtime Flow

For each number read from standard input:

1. `main` appends the number to the running history.
2. `PredictRange` scores the two ranges that were predicted before this new
   number arrived.
3. `PredictRange` asks two experts for the next range: `robustRange` and
   `quartileRange`.
4. The score selector chooses the expert with the stronger recent evidence.
5. The program prints that expert's range for the next value.

## Robust Statistics

`calculateStats` is the raw calculator: it returns average, first quartile,
median, third quartile, variance, and standard deviation for the exact slice it
receives.

`quartileRange` uses `recentSample` to keep the quartile statistics focused on
newer values without a fixed window size:

```text
sample size = sqrt(total history length)
```

That is a structural rule rather than a tuned window constant. It keeps more
history as the input grows, but still lets recent values dominate.

## Prediction Bands

After `calculateStats` returns, `PredictRange` uses the recent quartiles as
dynamic bands:

```text
lower band        -> value < first quartile
lower-middle band -> value < median
upper-middle band -> value < third quartile
upper band        -> value >= third quartile
```

Then the latest value is classified:

```text
latest value is in lower band        -> predict exact repeat
latest value is in lower-middle band -> predict one lower
latest value is in upper-middle band -> predict repeat or one higher
latest value is in upper band        -> learn a small downward pullback
```

Most of these predictions are one- or two-value ranges. That keeps the hit rate
low, but each hit is worth enough points to beat wider strategies.

## Robust Expert

`robustRange` uses a longer recent sample and filters values that are far from
the sample median before calculating its final statistics. It then:

1. treats very distant latest values as outliers and predicts near the filtered
   average;
2. predicts exact repeats in the lower and middle stable areas;
3. predicts a small downward pullback when the latest value is high.

This expert is less elegant than the quartile-only predictor, but it is a useful
counterweight because it handles broad outliers better.

## Score Selector

The selector does not know which data set is being tested. It simply scores both
experts after each actual next value arrives:

```text
raw score = 10000000 / range width
```

Only informative comparisons are kept. If neither expert would have scored on a
step, that step is ignored for selector memory.

For the next prediction, the selector compares the experts over the latest
informative score events. It keeps the robust expert unless the quartile expert
is ahead by a meaningful score margin. This avoids switching based on one lucky
hit while still allowing the quartile expert to take over when the current input
stream favors it.

## Learned High Pullback

The upper band is the only band with a learned range. Instead of using a fixed
`-5 -4` style rule, `learnedHighPullback` looks at recent transitions:

1. keep the same derived recent sample used by `PredictRange`;
2. for each adjacent pair in that sample, calculate the recent stats available
   at the previous value;
3. keep negative deltas that followed an upper-band previous value;
4. use the first-quartile-to-median range of those negative deltas as the next
   pullback.

If there is not enough recent upper-band history, the predictor falls back to an
exact repeat for that input.

## Constants

This branch still has constants, especially in the robust expert. The obvious
duplicate outlier knobs have been collapsed into one deviation scale: the
outlier distance, outlier-center range, score-switch margin, and high pullback
range all derive from that shared scale.

The robust expert's longer sample size is also derived from the deviation floor
and wide-side scale instead of being a standalone `128` window.

The selector no longer uses constants tied directly to the 12,500-line tester
files, such as an `8000`-line warmup or a `2000`-line lookback. Instead, it
works from scored events, so silent misses do not distort the expert comparison.

The remaining tradeoff is explainability. The predictor reaches the required
benchmark, but the robust expert is still more tuned than the quartile-only
version.
