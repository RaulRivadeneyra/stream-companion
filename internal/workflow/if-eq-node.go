package workflow

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

type IfEqNode struct {
	NodeID   string
	NodeType string
	Inputs   map[string]string
	True     string
	False    string
}

func (n *IfEqNode) ID() string {
	return n.NodeID
}

func (n *IfEqNode) Type() string {
	return "if_eq"
}

func (n *IfEqNode) Execute(ctx *ExecutionContext, plugins lua.LValue) (string, error) {
	// 1. Resolve inputs
	inputs, err := ResolveInputs(n.Inputs, ctx)
	if err != nil {
		return "", fmt.Errorf("if_eq node %s failed input resolution: %w", n.NodeID, err)
	}

	a, ok := inputs["a"]
	if !ok {
		return "", fmt.Errorf("if_eq node %s missing input 'a'", n.NodeID)
	}
	b, ok := inputs["b"]
	if !ok {
		return "", fmt.Errorf("if_eq node %s missing input 'b'", n.NodeID)
	}

	// 2. Compare using Lua equality
	L := lua.NewState()
	defer L.Close()
	equal := lua.LVAsBool(lua.LBool(a == b)) // Go comparison

	log.Printf("[if_eq:%s] Comparing: %v (%s) == %v (%s) → %v",
		n.NodeID,
		a, a.Type().String(),
		b, b.Type().String(),
		equal,
	)
	// 3. Store optional result (true/false) in context
	resultTable := L.NewTable()
	resultTable.RawSetString("result", lua.LBool(equal))
	ctx.Results[n.NodeID] = resultTable

	// 4. Return next node based on condition
	if equal {
		log.Printf("[if_eq:%s] Branch taken: TRUE → %s", n.NodeID, n.True)

		return n.True, nil
	}
	log.Printf("[if_eq:%s] Branch taken: FALSE → %s", n.NodeID, n.False)

	return n.False, nil
}
func (n *IfEqNode) ToJSON() NodeJSON {
	return NodeJSON{
		ID:     n.NodeID,
		Type:   "lua",
		Inputs: n.Inputs,
		True:   n.True,
		False:  n.False,
	}
}
