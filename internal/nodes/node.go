package nodes

import (
	"fmt"

	"github.com/google/uuid"
)

type INode[T any] interface {
	GetValue() T
	SetValue(T)
	GetId() uuid.UUID
	GetLabel() string
	SetLabel(string)
	GetNext() INode[any]
	GetPrev() INode[any]
	SetNext(INode[any])
	SetPrev(INode[any])
}

func IsHead(n INode[any]) bool {
	return n.GetPrev() == nil
}

func IsTail(n INode[any]) bool {
	return n.GetNext() == nil
}
func GetFullname(n INode[any]) string {
	return fmt.Sprintf("%s_%s", n.GetLabel(), n.GetId())
}
