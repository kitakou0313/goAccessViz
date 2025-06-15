# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Run
- `go build ./cmd/goAccessViz` - Build the main application
- `go run ./cmd/goAccessViz` - Run the application directly
- `go mod tidy` - Clean up dependencies

### Testing
- `go test ./...` - Run all tests in the project
- `go test ./cmd/goAccessViz/domain/node` - Run tests for specific package
- `go test -v ./...` - Run tests with verbose output

### Code Quality
- `go fmt ./...` - Format all Go code
- `go vet ./...` - Run Go static analysis tool

## Architecture

This is a Go application that visualizes function and database table relationships using Domain-Driven Design (DDD) principles. The project follows Clean Architecture with clear separation of concerns:

### Core Structure
- **Domain Layer** (`domain/node/`): Contains the core business logic with `Node` interface and implementations (`FunctionNode`, `DBTableNode`)
- **Application Layer** (`application/`): Handles DOT graph creation and conversion using the gonum graph library
- **Repository Layer** (`repository/`): Responsible for reading Go source code and building call graphs using `golang.org/x/tools`

### Key Components
- `Node` interface defines the contract for all node types with `GetChildren()` and `GetLabel()` methods
- `NewDotGraph()` converts domain nodes into DOT format graphs for visualization
- `ReadGraph()` analyzes Go packages to extract function call relationships using static analysis

### Dependencies
- `gonum.org/v1/gonum` for graph data structures and DOT format marshalling
- `golang.org/x/tools` for Go source code analysis and call graph generation

The application generates DOT format output that can be visualized with Graphviz tools.