# Tuning Results

This note records the constant-tuning experiment run against the local
`guess-it-dockerized` tester.

## Scope

The benchmark used the tester scoring formula from `server.js`:

```text
score += round(10000000 / (1 + high - low) / (data_length - 1))
```

Only data sets `1`, `2`, and `3` were included. Each group has five files, so
the required comparison set was:

```text
3 groups * 5 files * 3 required opponents = 45 comparisons
```

Required opponents:

- `big-range`
- `average`
- `median`

Bonus opponent:

- `nic`

The requested stop target was at least 90% wins against the required opponents.
That means at least `41/45` required wins, including strong performance against
each required opponent.

## Result

The best constants-only candidate found was:

```text
historyWindowSize         = 1
standardDeviationMultiple = 1
minimumRangeWidth         = 2
```

With a one-value window, the standard deviation is always `0`, and the momentum
adjustment does not run. The prediction is therefore effectively:

```text
last value - 2, last value + 2
```

That scored well on the local data, but it did not meet the 90% target:

```text
required wins: 35/45
win rate:      77.8%
```

Breakdown:

```text
big-range: 15/15
average:   14/15
median:     6/15
nic:        4/15  (bonus only)
```

## Interpretation

This result is useful as a benchmark finding, but it is not a satisfying
statistical solution. It mostly says that the target data often changes slowly
from one value to the next, so a very narrow range around the last value can
score highly.

It also shows the limit of tuning only constants in the current formula. The
best constants still lost most comparisons against `median`, so beating that
opponent consistently probably needs a formula change, not another pass over the
same three constants.

Because this heuristic does not pass the requested audit target and uses very
little of the intended statistics work, it should not be committed as the main
project strategy.

## Diagram

The generated visual report is included here:

[guess-it-tuning-results.html](guess-it-tuning-results.html)
