package workflow_test

import (
	"testing"

	"github.com/RaulRivadeneyra/stream-companion/commerce/workflow"
	lua "github.com/yuin/gopher-lua"
)

// Constructs a simple plugins table with plugins.hello(name)
func newPluginsTable(L *lua.LState) lua.LValue {
	tbl := L.NewTable()
	tbl.RawSetString("hello", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		L.Push(lua.LString("Hello, " + name))
		return 1
	}))
	return tbl
}

func TestRunLuaNode_ReturnsTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	plugins := newPluginsTable(L)
	inputs := map[string]lua.LValue{
		"name": lua.LString("Yujiko"),
	}

	code := `
		local msg = plugins.hello(input.name)
		return { result = msg, next = "done" }
	`

	ret, err := workflow.RunLuaNode(code, inputs, plugins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ret == nil {
		t.Fatalf("expected table, got nil")
	}

	result := ret.RawGetString("result")
	if result.Type() != lua.LTString {
		t.Errorf("expected result to be string, got %s", result.Type().String())
	}
}

func TestRunLuaNode_HandlesMissingNext(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	code := `return { result = 123 }`

	ret, err := workflow.RunLuaNode(code, nil, L.NewTable())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ret == nil {
		t.Fatalf("expected table, got nil")
	}

	next := ret.RawGetString("next")
	if next.Type() != lua.LTNil {
		t.Errorf("expected no 'next' field, got %s", next.Type().String())
	}
}

func TestRunLuaNode_ReturnsNonTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	code := `return "not a table"`

	ret, err := workflow.RunLuaNode(code, nil, L.NewTable())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ret != nil {
		t.Errorf("expected nil for non-table return, got type: %s", ret.Type().String())
	}
}
