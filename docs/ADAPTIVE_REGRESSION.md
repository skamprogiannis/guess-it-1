# Adaptive Regression Notes

> **Design exploration.** These notes capture an approach considered during
> development. The canonical `tester-optimized` branch uses robust statistics,
> quartile bands, and a rolling expert selector instead of linear regression.

`guess-it-1` asks the program to read one number at a time and print a range
where the next number is expected to fall. The input can be treated as a graph:
the line number is `x`, and the value read from standard input is `y`.

Adaptive regression is a way to estimate the next `y` value by looking at the
recent shape of that graph, then widening or narrowing the prediction range
based on how noisy the recent data has been.

## Core Idea

Instead of using every previous value equally, adaptive regression gives more
importance to recent values. This matters because the data may change direction
or volatility over time. A method that reacts to recent behavior can adjust
faster than one that only uses the full historical average.

The basic process is:

1. Keep the numbers received so far.
2. Focus on a recent window of values.
3. Estimate the current direction of movement from that window.
4. Predict the next value from that direction.
5. Measure recent prediction errors or spread.
6. Print a lower and upper bound around the prediction.

The prediction answers: "If the current local trend continues, where should the
next value probably be?"

The range width answers: "How much uncertainty should be allowed based on recent
behavior?"

## Trend

Regression estimates a line through recent points. In this project, the points
are simple:

- `x` is the input index: `0`, `1`, `2`, `3`, ...
- `y` is the number read from standard input.

If recent values are rising, the fitted line slopes upward. If recent values are
falling, it slopes downward. If recent values bounce around without a clear
direction, the slope stays close to flat.

The next prediction comes from extending that recent line one step forward.

A useful mental model:

- Recent average tells you the current level.
- Recent slope tells you the current direction.
- Recent error tells you how much trust to put in the prediction.

## Range Width

A correct but very wide range scores poorly. A narrow range scores better, but
only if the next value lands inside it. The important tuning problem is finding a
range that is usually wide enough without being wasteful.

The range should usually grow when recent data is unstable and shrink when recent
data is predictable.

Signals that the range should grow:

- recent values jump sharply up and down;
- recent predictions would have missed by a lot;
- the trend changes direction often;
- the newest value is far from the recent average.

Signals that the range can shrink:

- recent values move smoothly;
- recent prediction errors are small;
- values stay close to the fitted trend;
- volatility has been low for several inputs.

This is where the `math-skills` ideas fit naturally. Average, variance, and
standard deviation help describe the recent center and spread of the data.

## Recent Windows

The window size controls how much history affects the next prediction.

A short window reacts quickly. It can follow sudden changes, but it may overreact
to one unusual value.

A long window is smoother. It ignores random noise better, but it can lag behind
when the data starts moving in a new direction.

For this exercise, it is often useful to think in terms of several phases:

- With very little data, use a simple fallback range.
- With a small amount of data, use recent average and recent spread.
- With enough data, add trend estimation.
- With more data, tune the window size based on which recent behavior has worked
  best.

## Adaptation

The "adaptive" part means the algorithm changes its confidence as it learns.

A fixed range may work on one dataset and fail on another. Adaptive behavior lets
the program respond to the current dataset while it is running.

Examples of adaptive choices:

- Increase the range after recent misses would have happened.
- Decrease the range after many stable inputs.
- Use a shorter window when the trend changes quickly.
- Use a longer window when the data is noisy but directionally stable.
- Set a minimum range so early predictions are not too narrow.
- Set a maximum range so the program does not give away too much score.

The goal is not to perfectly predict every next value. The goal is to balance hit
rate against range size.

## Tuning Approach

Start simple and measure behavior.

Good things to track while testing:

- how often the next value lands inside the range;
- the average width of the printed ranges;
- which inputs cause misses;
- whether misses happen after trend changes, spikes, or quiet periods;
- whether the range stays too wide after volatility calms down.

When tuning, change one idea at a time. For example, first tune the recent window
size, then tune the range multiplier, then tune minimum and maximum widths.

Avoid tuning only for one dataset. The exercise mentions multiple datasets, so a
strategy that is slightly less perfect on one file but more stable across all
files is usually better.

## Practical Shape

A reasonable adaptive regression strategy can be described as:

- maintain a history of numbers;
- choose a recent window from that history;
- estimate the local trend from the window;
- predict the next number from that trend;
- estimate uncertainty from recent spread or recent prediction error;
- print prediction minus uncertainty and prediction plus uncertainty.

The exact constants are the part to experiment with. Window size, minimum range,
maximum range, and uncertainty multiplier all affect the score.

The best implementation should be fast, because the tester may send many input
values. The calculations should update quickly and avoid unnecessary work inside
the input loop.
