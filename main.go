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
	history := make([]int, 0)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		n, err := strconv.Atoi(line)
		if err != nil {
			continue
		}

		history = append(history, n)
		low, high := PredictRange(history)
		fmt.Printf("%d %d\n", low, high)
	}
}
