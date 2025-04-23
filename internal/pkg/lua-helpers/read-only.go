package luahelpers

import lua "github.com/yuin/gopher-lua"

// readOnlyValue returns either the original non-table value,
// or the cached/new proxy for a table.
func ReadOnlyValue(Lua *lua.LState, value lua.LValue) lua.LValue {
	table, isTableType := value.(*lua.LTable)

	if !isTableType {
		// No table : Return as it is
		return value
	}

	proxyTable := Lua.NewTable()
	metaTable := Lua.NewTable()

	// __index: read from original and wrap recursively
	Lua.SetField(metaTable, "__index", Lua.NewFunction(func(InnerLua *lua.LState) int {
		key := InnerLua.CheckAny(2)
		raw := table.RawGet(key)
		InnerLua.Push(ReadOnlyValue(InnerLua, raw))
		return 1
	}))

	// __newindex: forbid writes
	Lua.SetField(metaTable, "__newindex", Lua.NewFunction(func(InnerLua *lua.LState) int {
		InnerLua.RaiseError("table is readonly")
		return 0
	}))

	// TODO:  __pairs:: allow iteration over the original table
	/* Lua.SetField(metaTable, "__pairs", Lua.NewFunction(func(InnerLua *lua.LState) int {
		// call next on the original table to iterate
		InnerLua.Push(InnerLua.GetGlobal("next")) // `next` function
		InnerLua.Push(table)                      // original table
		InnerLua.Push(lua.LNil)                   // starting point for `next`
		return 3                                  // return next, table, nil (or key, value)
	})) */

	Lua.SetMetatable(proxyTable, metaTable)
	return proxyTable

}
