package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunPrintsOnePredictionPerValidInput(t *testing.T) {
	predictorState = scoreState{}
	input := strings.NewReader("100\n\nnot-a-number\n100\n")
	var output bytes.Buffer

	err := run(input, &output)

	if err != nil {
		t.Fatalf("run returned an unexpected error: %v", err)
	}
	const want = "100 100\n100 100\n"
	if output.String() != want {
		t.Fatalf("run output = %q, want %q", output.String(), want)
	}
}

func TestRunReturnsScannerErrors(t *testing.T) {
	input := strings.NewReader(strings.Repeat("1", 70*1024))
	var output bytes.Buffer

	err := run(input, &output)

	if err == nil {
		t.Fatal("run returned nil, want scanner error")
	}
	if output.Len() != 0 {
		t.Fatalf("run output = %q, want no output", output.String())
	}
}
