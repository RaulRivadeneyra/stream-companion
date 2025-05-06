package workflow

import (
	"fmt"

	"github.com/RaulRivadeneyra/stream-companion/core"
	lua "github.com/yuin/gopher-lua"
)

type Node interface {
	ID() string
	Type() string
	Execute(ctx *core.ExecutionContext, plugins lua.LValue) (next string, err error)
	ToJSON() NodeJSON
}

type NodeJSON struct {
	ID     string            `json:"id"`
	Type   string            `json:"type"`
	Inputs map[string]string `json:"inputs,omitempty"`
	Source string            `json:"source,omitempty"`
	Next   string            `json:"next,omitempty"`
	True   string            `json:"true,omitempty"`
	False  string            `json:"false,omitempty"`
}

func FromJSON(n NodeJSON) (Node, error) {
	switch n.Type {
	case "start":
		return &StartNode{
			NodeID: n.ID,
			Next:   n.Next,
		}, nil
	case "lua":
		return &LuaNode{
			NodeID: n.ID,
			Source: n.Source,
			Inputs: n.Inputs,
			Next:   n.Next,
		}, nil
	case "if_eq":
		return &IfEqNode{
			NodeID: n.ID,
			Inputs: n.Inputs,
			True:   n.True,
			False:  n.False,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported node type: %s", n.Type)
	}
}
