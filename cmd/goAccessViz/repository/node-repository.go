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

func ReadGraph(packagePath string) ([]node.TrackedEntity, error) {
	prog, pkgs, err := buildSSAProgramWithPackages(packagePath)
	if err != nil {
		return nil, err
	}

	// Build function call graph
	cg := cha.CallGraph(prog)
	nodeMap, childrenMap := buildNodeMaps(cg)

	// Add all functions from the package, not just those in call graph
	addAllPackageFunctions(prog, pkgs, nodeMap, childrenMap)

	// Analyze SQL strings and create DB table nodes
	sqlStrings := analyzePackageForSQL(pkgs)
	dbTableMap := createDBTableNodesMap(sqlStrings)

	// Establish function-to-table relationships
	establishFunctionTableRelationships(nodeMap, childrenMap, pkgs, dbTableMap)

	// Populate nodes with updated children (including SQL tables)
	populateNodes(nodeMap, childrenMap)

	// Return only function nodes (SQL table nodes are now children of functions)
	var allNodes []node.TrackedEntity
	for _, fnNode := range nodeMap {
		allNodes = append(allNodes, fnNode)
	}

	return allNodes, nil
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

func buildNodeMaps(cg *callgraph.Graph) (map[*ssa.Function]*node.FunctionTrackedEntity, map[*ssa.Function][]node.TrackedEntity) {
	nodeMap := make(map[*ssa.Function]*node.FunctionTrackedEntity)
	childrenMap := make(map[*ssa.Function][]node.TrackedEntity)
	callgraph.GraphVisitEdges(cg, createEdgeVisitor(nodeMap, childrenMap))
	return nodeMap, childrenMap
}

func createEdgeVisitor(nodeMap map[*ssa.Function]*node.FunctionTrackedEntity, childrenMap map[*ssa.Function][]node.TrackedEntity) func(*callgraph.Edge) error {
	return func(edge *callgraph.Edge) error {
		caller, callee := edge.Caller.Func, edge.Callee.Func
		ensureNodeExists(nodeMap, caller)
		ensureNodeExists(nodeMap, callee)
		childrenMap[caller] = append(childrenMap[caller], nodeMap[callee])
		return nil
	}
}

func ensureNodeExists(nodeMap map[*ssa.Function]*node.FunctionTrackedEntity, fn *ssa.Function) {
	if _, exists := nodeMap[fn]; !exists {
		nodeMap[fn] = &node.FunctionTrackedEntity{}
	}
}

func populateNodes(nodeMap map[*ssa.Function]*node.FunctionTrackedEntity, childrenMap map[*ssa.Function][]node.TrackedEntity) {
	for fn, fnNode := range nodeMap {
		children := childrenMap[fn]
		*fnNode = *node.NewFunctionTrackedEntity(fn.String(), children)
	}
}

func addAllPackageFunctions(prog *ssa.Program, pkgs []*packages.Package, nodeMap map[*ssa.Function]*node.FunctionTrackedEntity, childrenMap map[*ssa.Function][]node.TrackedEntity) {
	// Iterate through all SSA packages and their functions
	for _, ssaPkg := range prog.AllPackages() {
		// Check if this SSA package corresponds to one of our target packages
		for _, pkg := range pkgs {
			if ssaPkg.Pkg.Path() == pkg.PkgPath {
				// Add all functions from this package
				for _, member := range ssaPkg.Members {
					if fn, ok := member.(*ssa.Function); ok {
						// Only add if not already in nodeMap
						if _, exists := nodeMap[fn]; !exists {
							nodeMap[fn] = &node.FunctionTrackedEntity{}
							// Initialize empty children slice if not exists
							if _, exists := childrenMap[fn]; !exists {
								childrenMap[fn] = []node.TrackedEntity{}
							}
						}
					}
				}
			}
		}
	}
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

func createDBTableNodesMap(sqlStrings []string) map[string]*node.DatabaseTableTrackedEntity {
	tableMap := make(map[string]*node.DatabaseTableTrackedEntity)

	for _, sql := range sqlStrings {
		tables := extractTablesFromSQL(sql)
		for _, table := range tables {
			if _, exists := tableMap[table]; !exists {
				tableMap[table] = node.NewDatabaseTableTrackedEntity(table, []node.TrackedEntity{})
			}
		}
	}

	return tableMap
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

func establishFunctionTableRelationships(nodeMap map[*ssa.Function]*node.FunctionTrackedEntity, childrenMap map[*ssa.Function][]node.TrackedEntity, pkgs []*packages.Package, dbTableMap map[string]*node.DatabaseTableTrackedEntity) {
	// Map function names to their SSA functions for lookup
	funcNameToSSA := make(map[string]*ssa.Function)
	for ssaFunc := range nodeMap {
		if ssaFunc.Name() != "" {
			funcNameToSSA[ssaFunc.Name()] = ssaFunc
		}
	}

	// Analyze each package for SQL strings within functions
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				// Look for function declarations
				if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name != nil {
					funcName := funcDecl.Name.Name

					// Find the corresponding SSA function
					var targetSSAFunc *ssa.Function
					for ssaFunc := range nodeMap {
						if ssaFunc.Name() == funcName {
							targetSSAFunc = ssaFunc
							break
						}
					}

					if targetSSAFunc != nil {
						// Find SQL strings within this function
						sqlStringsInFunc := findSQLStringsInFunction(funcDecl)

						// For each SQL string, find referenced tables and add them as children
						for _, sqlStr := range sqlStringsInFunc {
							tables := extractTablesFromSQL(sqlStr)
							for _, tableName := range tables {
								if dbTableNode, exists := dbTableMap[tableName]; exists {
									// Add the table node as a child of this function
									childrenMap[targetSSAFunc] = append(childrenMap[targetSSAFunc], dbTableNode)
								}
							}
						}
					}
				}
				return true
			})
		}
	}
}

func findSQLStringsInFunction(funcDecl *ast.FuncDecl) []string {
	var sqlStrings []string

	if funcDecl.Body != nil {
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			if lit, ok := n.(*ast.BasicLit); ok && lit.Kind.String() == "STRING" {
				value := strings.Trim(lit.Value, `"'`+"`")
				if isSQLString(value) {
					sqlStrings = append(sqlStrings, value)
				}
			}
			return true
		})
	}

	return sqlStrings
}
