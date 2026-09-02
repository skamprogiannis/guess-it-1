# Guess It

Guess It is a streaming Go predictor built for the Zone01 statistics exercise.
It reads one integer at a time from standard input and prints an inclusive range
for the next value:

```text
lower upper
```

The scoring model rewards narrow correct ranges, so the predictor balances hit
rate against range width instead of trying to cover every possible value.

## Strategy

Two statistical experts produce a candidate range after every valid input:

- **Robust expert:** calculates recent statistics after filtering large
  outliers, then chooses an exact value or a small pullback range.
- **Quartile expert:** classifies the latest value into a recent quartile band
  and learns downward moves that followed earlier upper-band values.

When the next value arrives, both previous ranges are scored. A rolling selector
uses those results to choose the better expert while retaining a bias toward the
more outlier-resistant prediction. See [Strategy Notes](docs/STRATEGY.md) for
the formulas and selection flow.

## Architecture

```text
stdin
  -> parse valid integers
  -> append to history
  -> score the previous expert ranges
  -> robust range + quartile range
  -> rolling expert selector
  -> "lower upper" on stdout
```

- `main.go` owns the streaming input/output protocol.
- `predict.go` contains the experts, rolling score state, and range selection.
- `statistics.go` calculates mean, quartiles, median, variance, and standard
  deviation without mutating the input sample.

## Requirements

- Go 1.19 or newer

## Run

Use the audit-compatible script:

```sh
printf '189\n113\n121\n' | ./script.sh
```

Or build a standalone binary:

```sh
go build -o guess-it .
printf '189\n113\n121\n' | ./guess-it
```

Blank lines and non-integer input are ignored. Prediction ranges are the only
content written to standard output; scanner errors are written to standard
error and produce a non-zero exit status.

## Test

```sh
go test ./...
go test -race ./...
go vet ./...
```

The unit tests cover the streaming protocol, scanner failures, representative
statistics, stable sequences, and outlier filtering. The official Zone01
benchmark requires the separate `guess-it-dockerized` tester and its supplied
data files, which are not included in this repository.

## Recorded Results

The latest recorded local tester run for this branch won all 45 required
comparisons against `big-range`, `average`, and `median`, plus 9 of 15 bonus
comparisons against `nic`. These are historical measurements, not results
reproduced by CI; see [Tuning Results](docs/TUNING_RESULTS.md) for scope and
scoring details.

The [archived constants-only report](docs/guess-it-tuning-results.html)
documents an earlier unsuccessful experiment and is retained to show why the
strategy changed. It does not describe the current implementation.

## Status and Limitations

`tester-optimized` is the canonical portfolio branch. Alternative experiments
remain in Git history for comparison.

- The model was tuned for the Zone01 exercise datasets and should not be treated
  as a general forecasting system.
- Predictions are intentionally narrow, so many individual values fall outside
  the returned range even when the strategy scores well overall.
- Predictor state represents one input sequence per process.

## Author

Built by [Stefanos Kamprogiannis](https://github.com/skamprogiannis) during the
Zone01 Athens program.
