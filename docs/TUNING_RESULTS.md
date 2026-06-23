# Tuning Results

This note records local benchmark results from the `guess-it-dockerized` tester.
Only data sets `1`, `2`, and `3` were used because those are the exercise's
audit data sets.

## Scoring

The tester score comes from `server.js`:

```text
score += round(10000000 / (1 + high - low) / (data_length - 1))
```

A narrow correct range scores much more than a wide correct range. For a
12,500-value file, an exact one-number hit is worth about 800 points, while a
101-number range is worth about 8 points.

## Current Adaptive Result

The adaptive score-selector predictor was measured against five files in each
of data sets `1`, `2`, and `3`.

```text
big-range: 15/15
average:   15/15
median:    15/15
nic:        9/15
```

The required opponents are `big-range`, `average`, and `median`, so this gives:

```text
required wins: 45/45
required rate: 100.0%
```

`nic` is treated as a bonus comparison. This branch optimizes for the required
opponents first, while still reporting the bonus result.

## Strategy Evolution

Earlier experiments tried a constants-only strategy around the last value. The
best version behaved roughly like:

```text
last value - 2, last value + 2
```

That scored well against broad opponents, but it lost most comparisons against
`median` and used little of the intended statistics work.

The earlier no-tuning strategy used:

- a recent sample whose size is derived from the history length;
- dynamic bands derived from quartiles;
- narrow predictions selected from the latest value's band;
- a learned high-band pullback range based on recent upper-quartile
  transitions.

That was easier to explain, but it only reached `40/45` against the required
opponents.

The current branch combines that quartile predictor with a robust-statistics
predictor. Both experts are scored against the actual input stream as it
arrives. The selector only switches to the quartile predictor when recent
informative score events show that it has a clear edge. This removes the earlier
audit-shaped line-count settings (`8000` warmup and `2000` lookback), replacing
them with a short window of scored comparisons.

The remaining structural choices are documented in [STRATEGY.md](STRATEGY.md).
