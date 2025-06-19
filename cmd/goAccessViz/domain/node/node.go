package node

type TrackedEntity interface {
	GetChildren() []TrackedEntity
	GetLabel() string
}
