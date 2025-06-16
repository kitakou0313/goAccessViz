package repository

import (
	"go/ast"
	"goAccessViz/cmd/goAccessViz/domain/node"
	"regexp"
	"strings"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func ReadGraph(packagePath string) ([]node.Node, error) {
	prog, pkgs, err := buildSSAProgramWithPackages(packagePath)
	if err != nil {
		return nil, err
	}
	return createNodesFromCallGraphWithSQL(prog, pkgs), nil
}

func buildSSAProgram(packagePath string) (*ssa.Program, error) {
	cfg := createPackageConfig()
	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		return nil, err
	}
	return buildProgram(pkgs), nil
}

func buildSSAProgramWithPackages(packagePath string) (*ssa.Program, []*packages.Package, error) {
	cfg := createPackageConfig()
	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		return nil, nil, err
	}
	prog := buildProgram(pkgs)
	return prog, pkgs, nil
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

func createNodesFromCallGraphWithSQL(prog *ssa.Program, pkgs []*packages.Package) []node.Node {
	cg := cha.CallGraph(prog)
	nodeMap, childrenMap := buildNodeMaps(cg)
	populateNodes(nodeMap, childrenMap)

	// Analyze SQL strings and create DB table nodes
	sqlStrings := analyzePackageForSQL(pkgs)
	dbNodes := createDBTableNodes(sqlStrings)

	// Combine function nodes and DB table nodes
	var allNodes []node.Node
	for _, fnNode := range nodeMap {
		allNodes = append(allNodes, fnNode)
	}
	for _, dbNode := range dbNodes {
		allNodes = append(allNodes, dbNode)
	}

	return allNodes
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

// SQL analysis functions
func extractTablesFromSQL(sql string) []string {
	sql = strings.ToUpper(strings.TrimSpace(sql))
	var tables []string
	tableSet := make(map[string]bool)

	// Patterns for different SQL operations
	patterns := []string{
		`FROM\s+(\w+)`,
		`JOIN\s+(\w+)`,
		`INTO\s+(\w+)`,
		`UPDATE\s+(\w+)`,
		`DELETE\s+FROM\s+(\w+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(sql, -1)
		for _, match := range matches {
			if len(match) > 1 {
				tableName := strings.ToLower(match[1])
				if !tableSet[tableName] {
					tableSet[tableName] = true
					tables = append(tables, tableName)
				}
			}
		}
	}

	return tables
}

func detectSQLStrings(sourceCode string) []string {
	var sqlStrings []string

	// Simple regex to find string literals that look like SQL
	sqlPatterns := []string{
		`"[^"]*(?:SELECT|INSERT|UPDATE|DELETE|FROM|JOIN|INTO)[^"]*"`,
		`'[^']*(?:SELECT|INSERT|UPDATE|DELETE|FROM|JOIN|INTO)[^']*'`,
		"`[^`]*(?:SELECT|INSERT|UPDATE|DELETE|FROM|JOIN|INTO)[^`]*`",
	}

	for _, pattern := range sqlPatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindAllString(sourceCode, -1)
		for _, match := range matches {
			// Remove quotes
			cleaned := strings.Trim(match, `"'`+"`")
			sqlStrings = append(sqlStrings, cleaned)
		}
	}

	return sqlStrings
}

func createDBTableNodes(sqlStrings []string) []*node.DBTableNode {
	tableSet := make(map[string]bool)
	var dbNodes []*node.DBTableNode

	for _, sql := range sqlStrings {
		tables := extractTablesFromSQL(sql)
		for _, table := range tables {
			if !tableSet[table] {
				tableSet[table] = true
				dbNodes = append(dbNodes, node.NewDBTableNode(table, []node.Node{}))
			}
		}
	}

	return dbNodes
}

func analyzePackageForSQL(pkgs []*packages.Package) []string {
	var allSQLStrings []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				if lit, ok := n.(*ast.BasicLit); ok && lit.Kind.String() == "STRING" {
					value := strings.Trim(lit.Value, `"'`+"`")
					if isSQLString(value) {
						allSQLStrings = append(allSQLStrings, value)
					}
				}
				return true
			})
		}
	}

	return allSQLStrings
}

func isSQLString(s string) bool {
	s = strings.ToUpper(strings.TrimSpace(s))
	sqlKeywords := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "FROM", "JOIN", "INTO"}
	for _, keyword := range sqlKeywords {
		if strings.Contains(s, keyword) {
			return true
		}
	}
	return false
}
