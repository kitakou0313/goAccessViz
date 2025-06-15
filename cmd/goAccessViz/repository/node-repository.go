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
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
	}
	
	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		return nil, err
	}

	prog, _ := ssautil.Packages(pkgs, ssa.InstantiateGenerics)
	prog.Build()

	cg := cha.CallGraph(prog)
	
	nodeMap := make(map[*ssa.Function]*node.FunctionNode)
	childrenMap := make(map[*ssa.Function][]node.Node)
	
	callgraph.GraphVisitEdges(cg, func(edge *callgraph.Edge) error {
		caller := edge.Caller.Func
		callee := edge.Callee.Func
		
		if _, exists := nodeMap[caller]; !exists {
			nodeMap[caller] = &node.FunctionNode{}
		}
		if _, exists := nodeMap[callee]; !exists {
			nodeMap[callee] = &node.FunctionNode{}
		}
		
		childrenMap[caller] = append(childrenMap[caller], nodeMap[callee])
		
		return nil
	})
	
	for fn, fnNode := range nodeMap {
		children := childrenMap[fn]
		*fnNode = *node.NewFunctionNode(fn.String(), children)
	}
	
	var rootNodes []node.Node
	for _, fnNode := range nodeMap {
		rootNodes = append(rootNodes, fnNode)
	}
	
	return rootNodes, nil
}
