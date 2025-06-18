package main

import (
	"fmt"
	"os"

	"goAccessViz/cmd/goAccessViz/application"
	"goAccessViz/cmd/goAccessViz/repository"
)

func main() {
	// Get package path from command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: goAccessViz <package-path>")
		os.Exit(1)
	}

	packagePath := os.Args[1]

	// Read the graph with SQL analysis
	nodes, err := repository.ReadGraph(packagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading graph: %v\n", err)
		os.Exit(1)
	}

	// Convert to DOT graph
	dotGraph := application.NewDotGraph(nodes)

	// Convert to string and print
	convertedDotGraph, err := application.ConvertDotGraphToString(dotGraph)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting to DOT format: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(convertedDotGraph)
}
