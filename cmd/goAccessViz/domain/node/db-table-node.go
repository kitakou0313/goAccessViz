package node

// DBTableNode はSQLテーブルに相当するNode
type DBTableNode struct {
	tableName string
	children  []Node
}

func NewDBTableNode(tableName string, children []Node) *DBTableNode {
	return &DBTableNode{
		tableName: tableName,
		children:  children,
	}
}

func (dbtb *DBTableNode) GetChildren() []Node {
	return dbtb.children
}

func (dbtb *DBTableNode) GetLabel() string {
	return dbtb.tableName
}
