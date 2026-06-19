package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// main reads numbers from standard input and prints a range for each next value.
func main() {
	scanner := bufio.NewScanner(os.Stdin)
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
		fmt.Printf("%d %d\n", low, high)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
