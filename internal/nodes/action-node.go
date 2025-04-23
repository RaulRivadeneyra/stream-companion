package nodes

import (
	luahelpers "github.com/RaulRivadeneyra/stream-companion/internal/pkg/lua-helpers"
	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
)

const DEFAULT_ACTION_NODE_LABEL = "actionNode"

type ActionNode struct {
	id         uuid.UUID
	code       string //TODO: Make it so the whole code is retrieved rather than stored here
	label      string
	next, prev INode[any]
	Executor
}

func NewActionNode() *ActionNode {
	return &ActionNode{
		id:    uuid.New(),
		label: DEFAULT_ACTION_NODE_LABEL,
	}
}

func (an *ActionNode) GetId() uuid.UUID {
	return an.id
}

func (an *ActionNode) GetValue() any {
	return an.code
}

func (n *ActionNode) SetValue(val any) {
	if s, ok := val.(string); ok {
		n.code = s
	}
}

func (an *ActionNode) GetLabel() string {
	return an.label
}

func (an *ActionNode) SetLabel(l string) {
	an.label = l
}

func (an *ActionNode) GetNext() INode[any] {
	return an.next
}
func (an *ActionNode) GetPrev() INode[any] {
	return an.prev
}

func (an *ActionNode) SetNext(n INode[any]) {
	an.next = n
}
func (an *ActionNode) SetPrev(n INode[any]) {
	an.prev = n
}

func setInputs(L *lua.LState, svc *SharedVariableCollection) {
	inputsTable := L.NewTable()

	for _, v := range svc.GetSharedVariableList() {
		L.SetField(inputsTable, v.FullName(), v.GetLuaValue())
	}

	readonlyInput := luahelpers.ReadOnlyValue(L, inputsTable)

	L.SetGlobal("input", readonlyInput)
}

func (an *ActionNode) Execute(svc *SharedVariableCollection) error {
	L := lua.NewState()
	defer L.Close()
	setInputs(L, svc)

	output := L.NewTable()
	L.SetGlobal("output", output)

	// Execute some Lua code that modifies the output variable
	if err := L.DoString(an.code); err != nil {
		return err
	}

	L.ForEach(output, func(k, v lua.LValue) {
		sv := NewSharedVariable(an, k.String(), v)
		svc.AddSharedVariable(sv)
	})

	return nil
}
