package core

import lua "github.com/yuin/gopher-lua"

type ExecutionContext struct {
	Results map[string]*lua.LTable
}

type ExecutionResult struct {
	Context       *ExecutionContext // All node results
	FinalNodeID   string            // ID of the last executed node
	FinalResult   *lua.LTable       // The last node’s result (if any)
	Error         error             // Error that stopped execution, if any
	ExecutionPath []string          // Ordered list of node IDs executed
}
