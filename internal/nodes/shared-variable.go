package nodes

import (
	"fmt"

	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
)

type SharedVariable struct {
	name    string
	id      uuid.UUID
	luaType lua.LValueType
	value   any
	nodeRef INode[any]
}

func NewSharedVariable(n INode[any]) *SharedVariable {
	return &SharedVariable{
		id:      uuid.New(),
		nodeRef: n,
	}
}

func (sv *SharedVariable) GetName() string {
	return sv.name
}

func (sv *SharedVariable) SetName(name string) {
	sv.name = name
}

func (sv *SharedVariable) SetBool(b bool) {
	sv.luaType = lua.LTBool
	sv.value = b
}
func (sv *SharedVariable) SetFloat(f float64) {
	sv.luaType = lua.LTNumber
	sv.value = f
}
func (sv *SharedVariable) SetString(s string) {
	sv.luaType = lua.LTString
	sv.value = s
}

func (sv *SharedVariable) FullName() string {
	nodeLabel := sv.nodeRef.GetLabel()
	return fmt.Sprintf("%s_%s", nodeLabel, sv.name)
}

func (sv *SharedVariable) ToLuaValue() lua.LValue {
	switch v := sv.value.(type) {
	case string:
		return lua.LString(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	default:
		panic(fmt.Sprintf("Mismatched Go and Lua value type: %s - %s", v, sv.luaType))
	}
}

func (sv *SharedVariable) FromLuaValue(lv lua.LValue) {
	switch v := lv.(type) {
	case lua.LString:
		sv.SetString(lua.LVAsString(v))
	case lua.LNumber:
		sv.SetFloat(float64(lua.LVAsNumber(v)))
	case lua.LBool:
		sv.SetBool(lua.LVAsBool(v))
	default:
		panic(fmt.Sprintf(`LValueType '%s' is not allowed`, v.Type()))
	}
}

func (sv *SharedVariable) GetValue() any {
	return sv.value
}

func (sv *SharedVariable) GetLuaType() lua.LValueType {
	return sv.luaType
}
