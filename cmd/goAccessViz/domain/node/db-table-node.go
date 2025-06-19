package node

// DatabaseTableTrackedEntity はSQLテーブルに相当するNode
type DatabaseTableTrackedEntity struct {
	tableName string
	children  []TrackedEntity
}

func NewDatabaseTableTrackedEntity(tableName string, children []TrackedEntity) *DatabaseTableTrackedEntity {
	return &DatabaseTableTrackedEntity{
		tableName: tableName,
		children:  children,
	}
}

func (dbtb *DatabaseTableTrackedEntity) GetChildren() []TrackedEntity {
	return dbtb.children
}

func (dbtb *DatabaseTableTrackedEntity) GetLabel() string {
	return dbtb.tableName
}
