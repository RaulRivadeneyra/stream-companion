package nodes

import (
	"fmt"

	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
)

type SharedVariable struct {
	name     string
	id       uuid.UUID
	luaType  lua.LValueType
	value    any
	luaValue lua.LValue
	nodeRef  INode[any]
}

func NewSharedVariable(node INode[any], name string, value any) *SharedVariable {

	sharedVar := SharedVariable{
		id:      uuid.New(),
		nodeRef: node,
		name:    name,
	}

	setValue(&sharedVar, value)

	return &sharedVar
}

func (sharedVar *SharedVariable) GetName() string {
	return sharedVar.name
}

func (sharedVar *SharedVariable) FullName() string {
	nodeLabel := sharedVar.nodeRef.GetLabel()
	return fmt.Sprintf("%s_%s", nodeLabel, sharedVar.name)
}

func (sharedVar *SharedVariable) GetValue() any {
	return sharedVar.value
}

func (sharedVar *SharedVariable) GetLuaValue() lua.LValue {
	return sharedVar.luaValue
}

func (sharedVar *SharedVariable) GetLuaType() lua.LValueType {
	return sharedVar.luaType
}

func coerceValue(value any) (lua.LValueType, lua.LValue, any, bool) {
	switch value := value.(type) {
	case string:
		return lua.LTString, lua.LString(value), value, true
	case float64:
		return lua.LTNumber, lua.LNumber(value), value, true
	case bool:
		return lua.LTBool, lua.LBool(value), value, true
	case lua.LString:
		return lua.LTString, lua.LString(value), lua.LVAsString(value), true
	case lua.LNumber:
		return lua.LTNumber, lua.LNumber(value), float64(lua.LVAsNumber(value)), true
	case lua.LBool:
		return lua.LTBool, lua.LBool(value), lua.LVAsBool(value), true
	default:
		return 0, nil, nil, false
	}
}

func setValue(sharedVar *SharedVariable, value any) {

	luaType, luaVal, rawVal, canCoerce := coerceValue(value)

	if !canCoerce {
		panic(fmt.Sprintf("Passed value doesn't match any valid type: %s ", value))
	}

	sharedVar.luaType = luaType
	sharedVar.luaValue = luaVal
	sharedVar.value = rawVal
}
