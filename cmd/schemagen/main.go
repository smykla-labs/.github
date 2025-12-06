// Package main provides a simple CLI to generate JSON Schema for sync configuration.
// This binary is not released; it's used via `go run` in workflows.
package main

import (
	"fmt"
	"os"

	"github.com/smykla-labs/.github/pkg/schema"
)

func main() {
	output, err := schema.GenerateSchema(
		"github.com/smykla-labs/.github",
		"./pkg/config",
	)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)

		os.Exit(1)
	}

	fmt.Print(string(output))
}
