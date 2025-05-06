package executor

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

type StartNode struct {
	NodeID string
	Next   string
}

func (n *StartNode) ID() string   { return n.NodeID }
func (n *StartNode) Type() string { return "start" }

func (n *StartNode) Execute(ctx *ExecutionContext, plugins lua.LValue) (string, error) {
	if n.Next == "" {
		return "", fmt.Errorf("start node must define a `next` field")
	}
	log.Printf("[start:%s] Workflow begins → %s", n.NodeID, n.Next)
	return n.Next, nil
}

func (n *StartNode) ToJSON() NodeJSON {
	return NodeJSON{
		ID:   n.NodeID,
		Type: "start",
		Next: n.Next,
	}
}
