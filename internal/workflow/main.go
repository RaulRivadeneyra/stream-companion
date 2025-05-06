package workflow

import (
	"fmt"
	"log"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type Workflow struct {
	Nodes map[string]Node
}

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

func ExecuteWorkflow(workflow Workflow, plugins lua.LValue) ExecutionResult {
	ctx := &ExecutionContext{Results: make(map[string]*lua.LTable)}
	result := ExecutionResult{
		Context:       ctx,
		ExecutionPath: []string{},
	}

	currentID := "start"

	for {
		node, ok := workflow.Nodes[currentID]
		if !ok {
			err := fmt.Errorf("node %s not found in workflow", currentID)
			result.Error = err
			result.FinalNodeID = currentID
			return result
		}

		log.Printf("[workflow] Executing node: %s (%s)", currentID, node.Type())
		result.ExecutionPath = append(result.ExecutionPath, currentID)

		nextID, err := node.Execute(ctx, plugins)
		if err != nil {
			result.Error = err
			result.FinalNodeID = currentID
			return result
		}

		// Save current as last successfully executed
		result.FinalNodeID = currentID
		result.FinalResult = ctx.Results[currentID]

		// Done if there's nowhere else to go
		if nextID == "" {
			return result
		}

		currentID = nextID
	}
}

func ResolveInputs(
	inputMap map[string]string, ctx *ExecutionContext,
) (map[string]lua.LValue, error) {
	resolved := make(map[string]lua.LValue)

	for key, val := range inputMap {
		if strings.HasPrefix(val, "input.") {
			parts := strings.Split(val, ".")
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid input reference: %s", val)
			}
			nodeID, field := parts[1], parts[2]
			resultTable, ok := ctx.Results[nodeID]
			if !ok {
				return nil, fmt.Errorf("node %s not executed yet", nodeID)
			}
			resolved[key] = resultTable.RawGetString(field)
		} else {
			resolved[key] = lua.LString(val) // assume static string for now
		}
	}

	return resolved, nil
}
