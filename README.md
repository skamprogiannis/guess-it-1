# guess-it-1

`guess-it-1` reads numbers from standard input and prints a range for the next
expected number after each input line.

The tester treats the input as a graph:

- `x` is the input line index: `0`, `1`, `2`, ...
- `y` is the number received on that line.

The program must respond with:

```text
lower upper
```

Only prediction ranges should be written to standard output.

## Run

```sh
printf '189\n113\n121\n' | ./script.sh
```

`script.sh` sets `GOCACHE` to a writable directory under `/tmp` when one is not
already configured, changes into the script directory, and runs the Go program.

## Test

```sh
env GOCACHE=/tmp/guess-it-go-build-cache go test ./...
```

To use the Zone01 tester, copy the required project files into the tester's
`student/` directory and run the tester from its root.

## Notes

- [Adaptive regression notes](docs/ADAPTIVE_REGRESSION.md)
- [Tuning results](docs/TUNING_RESULTS.md)
- [Tuning results diagram](docs/guess-it-tuning-results.html)
