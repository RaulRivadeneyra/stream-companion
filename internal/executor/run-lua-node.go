package executor

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type Inputs = map[string]lua.LValue

const VM_CALL_STACK_SIZE = 256
const VM_REGISTRY_SIZE = 256

func RunLuaNode(source string, inputs Inputs, pluginsTable lua.LValue) (*lua.LTable, error) {
	luaState := lua.NewState(lua.Options{
		CallStackSize: VM_CALL_STACK_SIZE,
		RegistrySize:  VM_REGISTRY_SIZE,
	})
	defer luaState.Close()

	// Inject inputs
	inputTable := luaState.NewTable()
	for k, v := range inputs {
		inputTable.RawSetString(k, v)
	}
	luaState.SetGlobal("input", inputTable)

	// Inject plugins
	luaState.SetGlobal("plugins", pluginsTable)

	err := luaState.DoString(source)
	resultTable := luaState.NewTable()
	if err != nil {
		return resultTable, fmt.Errorf("node execution error: %w", err)
	}

	ret := luaState.Get(-1)

	tbl, ok := ret.(*lua.LTable)
	if !ok {
		return tbl, nil
	}

	return tbl, nil
}
