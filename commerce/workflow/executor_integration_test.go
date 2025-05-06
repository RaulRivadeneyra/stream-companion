package workflow_test

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/RaulRivadeneyra/stream-companion/internal/workflow"
	lua "github.com/yuin/gopher-lua"
)

func TestExecuteWorkflow_SimpleLinearChain(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a mock plugins table (empty for now)
	plugins := L.NewTable()

	// Build workflow
	workflow := workflow.Workflow{
		Nodes: map[string]workflow.Node{
			"start": &workflow.StartNode{
				NodeID: "start",
				Next:   "get_name",
			},
			"get_name": &workflow.LuaNode{
				NodeID: "get_name",
				Source: `return { name = "Yujiko" }`,
				Inputs: map[string]string{},
				Next:   "format_message",
			},
			"format_message": &workflow.LuaNode{
				NodeID: "format_message",
				Source: `return { result = "Hello " .. input.name }`,
				Inputs: map[string]string{
					"name": "input.get_name.name",
				},
				Next: "",
			},
		},
	}

	result := workflow.ExecuteWorkflow(workflow, plugins)

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

	workflow := workflow.Workflow{
		Nodes: map[string]workflow.Node{
			"start": &workflow.StartNode{
				NodeID: "start",
				Next:   "get_type",
			},
			"get_type": &workflow.LuaNode{
				NodeID: "get_type",
				Source: `return { type = "electric" }`,
				Inputs: map[string]string{},
				Next:   "check_type",
			},
			"check_type": &workflow.IfEqNode{
				NodeID: "check_type",
				Inputs: map[string]string{
					"a": "input.get_type.type",
					"b": "electric",
				},
				True:  "handle_electric",
				False: "handle_other",
			},
			"handle_electric": &workflow.LuaNode{
				NodeID: "handle_electric",
				Source: `return { result = "⚡ it’s electric!" }`,
				Inputs: map[string]string{},
				Next:   "",
			},
			"handle_other": &workflow.LuaNode{
				NodeID: "handle_other",
				Source: `return { result = "🔥 not electric." }`,
				Inputs: map[string]string{},
				Next:   "",
			},
		},
	}

	result := workflow.ExecuteWorkflow(workflow, plugins)

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
func TestWorkflowworkflow_ExecutesFromJSON(t *testing.T) {
	jsonData := `
	{
		"nodes": [
			{ "id": "start", "type": "start", "next": "hello" },
			{ "id": "hello", "type": "lua", "source": "return { result = 'Hi' }", "inputs": {}, "next": "goodbye" },
			{ "id": "goodbye", "type": "lua", "source": "return { result = input.msg }", "inputs": { "msg": "input.hello.result" } }
		]
	}`

	var raw struct {
		Nodes []workflow.NodeJSON `json:"nodes"`
	}
	if err := json.Unmarshal([]byte(jsonData), &raw); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	nodesMap := map[string]workflow.Node{}
	for _, nj := range raw.Nodes {
		node, err := workflow.FromJSON(nj)
		if err != nil {
			t.Fatalf("Invalid node: %v", err)
		}
		nodesMap[nj.ID] = node
	}

	workflow := workflow.Workflow{Nodes: nodesMap}

	L := lua.NewState()
	defer L.Close()

	plugins := L.NewTable()
	result := workflow.ExecuteWorkflow(workflow, plugins)

	if result.Error != nil {
		t.Fatalf("Unexpected error: %v", result.Error)
	}
	if result.FinalNodeID != "goodbye" {
		t.Errorf("Final node should be goodbye, got %s", result.FinalNodeID)
	}
	if result.FinalResult.RawGetString("result").String() != "Hi" {
		t.Errorf("Unexpected output: %s", result.FinalResult.RawGetString("result").String())
	}
}

func TestWorkflowworkflow_LogsExecutionPath(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	workflow := workflow.Workflow{
		Nodes: map[string]workflow.Node{
			"start": &workflow.StartNode{NodeID: "start", Next: "say"},
			"say":   &workflow.LuaNode{NodeID: "say", Source: `return { result = "yo" }`, Inputs: map[string]string{}},
		},
	}

	L := lua.NewState()
	plugins := L.NewTable()

	result := workflow.ExecuteWorkflow(workflow, plugins)

	if result.Error != nil {
		t.Fatalf("Execution failed: %v", result.Error)
	}
	if !strings.Contains(buf.String(), "[workflow] Executing node: say") {
		t.Errorf("Expected log output for node 'say' missing: %s", buf.String())
	}
}

func TestWorkflowworkflow_TerminatesGracefullyOnMissingNext(t *testing.T) {
	workflow := workflow.Workflow{
		Nodes: map[string]workflow.Node{
			"start": &workflow.StartNode{NodeID: "start", Next: "one"},
			"one":   &workflow.LuaNode{NodeID: "one", Source: `return { result = "done" }`, Inputs: map[string]string{}},
		},
	}

	L := lua.NewState()
	plugins := L.NewTable()

	result := workflow.ExecuteWorkflow(workflow, plugins)

	if result.Error != nil {
		t.Fatalf("Should have terminated cleanly, but got error: %v", result.Error)
	}
	if result.FinalNodeID != "one" {
		t.Errorf("Expected to end at 'one', got %s", result.FinalNodeID)
	}
}

func TestWorkflowworkflow_HandlesNodeFailure(t *testing.T) {
	workflow := workflow.Workflow{
		Nodes: map[string]workflow.Node{
			"start": &workflow.StartNode{NodeID: "start", Next: "fail"},
			"fail": &workflow.LuaNode{
				NodeID: "fail",
				Source: `return { result =  }`, // will error
				Inputs: map[string]string{},
			},
		},
	}

	L := lua.NewState()
	plugins := L.NewTable()

	result := workflow.ExecuteWorkflow(workflow, plugins)

	if result.Error == nil {
		t.Fatal("Expected error from failing node, got nil")
	}
	if result.FinalNodeID != "fail" {
		t.Errorf("Expected final node to be 'fail', got %s", result.FinalNodeID)
	}
}
