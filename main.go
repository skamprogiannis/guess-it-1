package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func run(input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	var history []int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		num, err := strconv.Atoi(line)
		if err != nil {
			continue
		}

		history = append(history, num)
		low, high := PredictRange(history)
		fmt.Fprintf(output, "%d %d\n", low, high)
	}

	return scanner.Err()
}

// main reads numbers from standard input and prints a range for each next value.
func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
