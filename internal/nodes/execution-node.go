package nodes

import "github.com/google/uuid"

type Variables = map[string]any

type ExecutionNode interface {
	Execute(input Variables) (Variables, error)
	NextNode() *ExecutionNode
	Id() uuid.UUID
	Label() string
}
