package nodes

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
)

type code struct {
	inputVars  []string
	outputVars []string
	body       string
}

func (c *code) BodyFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	c.body = string(content)
	return nil
}

type actionNode struct {
	code     string
	label    string
	id       uuid.UUID
	nextNode *ExecutionNode
}

const DEFAULT_ACTION_NODE_LABEL = "actionNode"

func NewActionNode() *actionNode {
	return &actionNode{
		label: DEFAULT_ACTION_NODE_LABEL,
		id:    uuid.New(),
	}
}

func (an *actionNode) Execute(input Variables) (Variables, error) {
	L := lua.NewState()
	defer L.Close()

	inputsTable := L.NewTable()

	for k, v := range input {
		L.SetField(inputsTable, k, convertToLuaValue(v))
	}

	L.SetGlobal("input", inputsTable)

	output := L.NewTable()
	L.SetGlobal("output", output)

	// Execute some Lua code that modifies the output variable
	if err := L.DoString(an.code); err != nil {
		return nil, err
	}

	// Retrieve the modified output table
	L.GetGlobal("output")

	// Convert the Lua table to a map[string]any
	goMap := make(map[string]any)
	L.ForEach(output, func(k, v lua.LValue) {
		key := k.String()
		value := convertFromLuaValue(v)
		goMap[key] = value
	})

	return goMap, nil
}

func (an *actionNode) SetCode(s string) {
	an.code = s
}

func (an *actionNode) SetNextNode(nn *ExecutionNode) {
	an.nextNode = nn
}

func (an *actionNode) NextNode() *ExecutionNode {
	return an.nextNode
}

func convertToLuaValue(value any) lua.LValue {
	switch v := value.(type) {
	case string:
		return lua.LString(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	default:
		return lua.LNil
	}
}

func convertFromLuaValue(value lua.LValue) any {
	fmt.Println(value.Type())
	switch v := value.(type) {
	case lua.LString:
		return lua.LVAsString(v)
	case lua.LNumber:
		return float64(lua.LVAsNumber(v))
	case lua.LBool:
		return lua.LVAsBool(v)
	default:
		return nil
	}
}
