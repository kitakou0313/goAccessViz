# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

goAccessViz is a Go application that visualizes function and database table relationships by analyzing Go source code. It generates DOT format graphs that can be rendered with Graphviz tools.

## Usage

```bash
go run ./cmd/goAccessViz <package-path>
```

The application outputs DOT format to stdout, which can be piped to Graphviz:
```bash
go run ./cmd/goAccessViz ./testpkg | dot -Tpng -o graph.png
```

## Commands

### Build and Run
- `go build ./cmd/goAccessViz` - Build the main application
- `go run ./cmd/goAccessViz <package-path>` - Run the application with target package
- `go mod tidy` - Clean up dependencies

### Testing
- `go test ./...` - Run all tests in the project
- `go test ./cmd/goAccessViz/domain/node` - Run tests for specific package
- `go test -v ./...` - Run tests with verbose output

### Code Quality
- `go fmt ./...` - Format all Go code
- `go vet ./...` - Run Go static analysis tool

## Development Guidelines

### Test-Driven Development (TDD)
When implementing new features or fixing bugs, follow Test-Driven Development practices:

1. **Red**: Write a failing test first that describes the desired behavior
2. **Green**: Write the minimal code necessary to make the test pass
3. **Refactor**: Improve the code while keeping tests green

Always write tests before implementing functionality to ensure:
- Clear specification of requirements
- Better code design and modularity
- Confidence in refactoring and changes
- Comprehensive test coverage

## Architecture

This project follows Clean Architecture and Domain-Driven Design (DDD) principles with clear separation of concerns:

### Core Structure
- **Domain Layer** (`cmd/goAccessViz/domain/node/`): Contains the core business logic
  - `TrackedEntity` interface: Contract for all node types with `GetChildren()` and `GetLabel()` methods
  - `FunctionTrackedEntity`: Represents Go functions in the graph
  - `DatabaseTableTrackedEntity`: Represents database tables accessed by functions
- **Application Layer** (`cmd/goAccessViz/application/`): Handles DOT graph creation and conversion
  - `NewDotGraph()`: Converts domain nodes into DOT format graphs for visualization
  - `ConvertDotGraphToString()`: Marshals graphs to DOT string format
- **Repository Layer** (`cmd/goAccessViz/repository/`): Responsible for code analysis and graph building
  - `ReadGraph()`: Main entry point that analyzes Go packages and returns TrackedEntity nodes
  - Static analysis using golang.org/x/tools for function call graph generation
  - SQL string detection and table extraction from both raw strings and sqlx method calls

### Key Features
1. **Function Call Graph Analysis**: Uses SSA (Single Static Assignment) and CHA (Class Hierarchy Analysis) to build complete function call graphs
2. **SQL Analysis**: Detects SQL strings in source code and extracts referenced table names
3. **Database Integration**: Recognizes sqlx method calls (`Get`, `Select`, `Exec`, `Query`, etc.) and extracts SQL from their arguments
4. **Graph Visualization**: Converts the analyzed relationships into DOT format for visualization

### Dependencies
- `gonum.org/v1/gonum` for graph data structures and DOT format marshalling
- `golang.org/x/tools` for Go source code analysis and call graph generation
- `github.com/jmoiron/sqlx` (recognized for SQL analysis)

### File Structure
```
cmd/goAccessViz/
├── main.go                    # Entry point
├── domain/node/               # Domain layer
│   ├── node.go               # TrackedEntity interface
│   ├── function-node.go      # Function node implementation
│   ├── db-table-node.go      # Database table node implementation
│   └── node_test.go          # Domain tests
├── application/              # Application layer
│   ├── dot-node.go          # DOT graph creation
│   └── dot-node_test.go     # Application tests
└── repository/              # Repository layer
    ├── node-repository.go   # Graph analysis and building
    └── readgraph_test.go    # Repository tests
```

The application generates DOT format output that can be visualized with Graphviz tools to show function call relationships and database table dependencies.