package executor

import lua "github.com/yuin/gopher-lua"

type LuaNode struct {
	NodeID   string
	NodeType string
	Source   string
	Inputs   map[string]string
	Next     string
}

func (n *LuaNode) ID() string {
	return n.NodeID
}

func (n *LuaNode) Type() string {
	return "lua"
}

func (n *LuaNode) Execute(ctx *ExecutionContext, plugins lua.LValue) (string, error) {
	inputs, err := ResolveInputs(n.Inputs, ctx)
	if err != nil {
		return "", err
	}

	result, err := RunLuaNode(n.Source, inputs, plugins)
	if err != nil {
		return "", err
	}

	ctx.Results[n.NodeID] = result
	return n.Next, nil
}

func (n *LuaNode) ToJSON() NodeJSON {
	return NodeJSON{
		ID:     n.NodeID,
		Type:   "lua",
		Source: n.Source,
		Inputs: n.Inputs,
		Next:   n.Next,
	}
}
