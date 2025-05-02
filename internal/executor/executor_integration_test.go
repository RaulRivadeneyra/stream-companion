package executor_test

import (
	"testing"

	"github.com/RaulRivadeneyra/stream-companion/internal/executor"
	lua "github.com/yuin/gopher-lua"
)

func TestExecuteWorkflow_SimpleLinearChain(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a mock plugins table (empty for now)
	plugins := L.NewTable()

	// Build workflow
	workflow := executor.Workflow{
		Nodes: map[string]executor.Node{
			"start": &executor.StartNode{
				NodeID: "start",
				Next:   "get_name",
			},
			"get_name": &executor.LuaNode{
				NodeID: "get_name",
				Source: `return { name = "Yujiko" }`,
				Inputs: map[string]string{},
				Next:   "format_message",
			},
			"format_message": &executor.LuaNode{
				NodeID: "format_message",
				Source: `return { result = "Hello " .. input.name }`,
				Inputs: map[string]string{
					"name": "input.get_name.name",
				},
				Next: "",
			},
		},
	}

	result := executor.ExecuteWorkflow(workflow, plugins)

	if result.Error != nil {
		t.Fatalf("Workflow failed at node %s: %v", result.FinalNodeID, result.Error)
	}

	if result.FinalNodeID != "format_message" {
		t.Errorf("Expected final node to be format_message, got %s", result.FinalNodeID)
	}

	if result.FinalResult == nil {
		t.Fatal("Final result is nil")
	}

	output := result.FinalResult.RawGetString("result")
	if output.Type() != lua.LTString {
		t.Errorf("Expected result to be string, got %s", output.Type().String())
	}

	if output.String() != "Hello Yujiko" {
		t.Errorf("Unexpected output: got %s, want Hello Yujiko", output.String())
	}
}

func TestExecuteWorkflow_BranchingWithIfEq(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	plugins := L.NewTable()

	workflow := executor.Workflow{
		Nodes: map[string]executor.Node{
			"start": &executor.StartNode{
				NodeID: "start",
				Next:   "get_type",
			},
			"get_type": &executor.LuaNode{
				NodeID: "get_type",
				Source: `return { type = "electric" }`,
				Inputs: map[string]string{},
				Next:   "check_type",
			},
			"check_type": &executor.IfEqNode{
				NodeID: "check_type",
				Inputs: map[string]string{
					"a": "input.get_type.type",
					"b": "electric",
				},
				True:  "handle_electric",
				False: "handle_other",
			},
			"handle_electric": &executor.LuaNode{
				NodeID: "handle_electric",
				Source: `return { result = "⚡ it’s electric!" }`,
				Inputs: map[string]string{},
				Next:   "",
			},
			"handle_other": &executor.LuaNode{
				NodeID: "handle_other",
				Source: `return { result = "🔥 not electric." }`,
				Inputs: map[string]string{},
				Next:   "",
			},
		},
	}

	result := executor.ExecuteWorkflow(workflow, plugins)

	if result.Error != nil {
		t.Fatalf("Workflow failed at node %s: %v", result.FinalNodeID, result.Error)
	}

	switch result.FinalNodeID {
	case "handle_electric":
		val := result.FinalResult.RawGetString("result").String()
		if val != "⚡ it’s electric!" {
			t.Errorf("Unexpected result: %s", val)
		}
	case "handle_other":
		t.Errorf("Wrong branch taken — expected electric")
	default:
		t.Errorf("Unexpected final node: %s", result.FinalNodeID)
	}
}
