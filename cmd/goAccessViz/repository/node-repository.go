package repository

import (
	"goAccessViz/cmd/goAccessViz/domain/node"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func ReadGraph(packagePath string) ([]node.Node, error) {
	prog, err := buildSSAProgram(packagePath)
	if err != nil {
		return nil, err
	}
	return createNodesFromCallGraph(prog), nil
}

func buildSSAProgram(packagePath string) (*ssa.Program, error) {
	cfg := createPackageConfig()
	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		return nil, err
	}
	return buildProgram(pkgs), nil
}

func createPackageConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
	}
}

func buildProgram(pkgs []*packages.Package) *ssa.Program {
	prog, _ := ssautil.Packages(pkgs, ssa.InstantiateGenerics)
	prog.Build()
	return prog
}

func createNodesFromCallGraph(prog *ssa.Program) []node.Node {
	cg := cha.CallGraph(prog)
	nodeMap, childrenMap := buildNodeMaps(cg)
	populateNodes(nodeMap, childrenMap)
	return collectRootNodes(nodeMap)
}

func buildNodeMaps(cg *callgraph.Graph) (map[*ssa.Function]*node.FunctionNode, map[*ssa.Function][]node.Node) {
	nodeMap := make(map[*ssa.Function]*node.FunctionNode)
	childrenMap := make(map[*ssa.Function][]node.Node)
	callgraph.GraphVisitEdges(cg, createEdgeVisitor(nodeMap, childrenMap))
	return nodeMap, childrenMap
}

func createEdgeVisitor(nodeMap map[*ssa.Function]*node.FunctionNode, childrenMap map[*ssa.Function][]node.Node) func(*callgraph.Edge) error {
	return func(edge *callgraph.Edge) error {
		caller, callee := edge.Caller.Func, edge.Callee.Func
		ensureNodeExists(nodeMap, caller)
		ensureNodeExists(nodeMap, callee)
		childrenMap[caller] = append(childrenMap[caller], nodeMap[callee])
		return nil
	}
}

func ensureNodeExists(nodeMap map[*ssa.Function]*node.FunctionNode, fn *ssa.Function) {
	if _, exists := nodeMap[fn]; !exists {
		nodeMap[fn] = &node.FunctionNode{}
	}
}

func populateNodes(nodeMap map[*ssa.Function]*node.FunctionNode, childrenMap map[*ssa.Function][]node.Node) {
	for fn, fnNode := range nodeMap {
		children := childrenMap[fn]
		*fnNode = *node.NewFunctionNode(fn.String(), children)
	}
}

func collectRootNodes(nodeMap map[*ssa.Function]*node.FunctionNode) []node.Node {
	var rootNodes []node.Node
	for _, fnNode := range nodeMap {
		rootNodes = append(rootNodes, fnNode)
	}
	return rootNodes
}
