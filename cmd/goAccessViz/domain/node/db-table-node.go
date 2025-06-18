package node

// DatabaseTableGraphNode はSQLテーブルに相当するNode
type DatabaseTableGraphNode struct {
	tableName string
	children  []GraphNode
}

func NewDatabaseTableGraphNode(tableName string, children []GraphNode) *DatabaseTableGraphNode {
	return &DatabaseTableGraphNode{
		tableName: tableName,
		children:  children,
	}
}

func (dbtb *DatabaseTableGraphNode) GetChildren() []GraphNode {
	return dbtb.children
}

func (dbtb *DatabaseTableGraphNode) GetLabel() string {
	return dbtb.tableName
}
