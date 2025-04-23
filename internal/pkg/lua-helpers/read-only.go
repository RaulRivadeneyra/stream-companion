package luahelpers

import lua "github.com/yuin/gopher-lua"

// proxies remembers, for each original table, its readonly proxy.
// (Go maps aren’t weak, so these entries will stick around
// until the Lua state is closed—but it matches the Lua version.)
var proxies = map[*lua.LTable]*lua.LTable{}

// Loader registers the module so in Lua you can do:
//
//	local ro = require("readonly")
//	local t = ro.readOnly(myTable)
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"readOnly": readOnly,
	})
	L.Push(mod)
	return 1
}

// readOnly is the exposed Lua function.
func readOnly(L *lua.LState) int {
	val := L.CheckAny(1)
	L.Push(ReadOnlyValue(L, val))
	return 1
}

// readOnlyValue returns either the original non-table value,
// or the cached/new proxy for a table.
func ReadOnlyValue(L *lua.LState, val lua.LValue) lua.LValue {
	if tbl, ok := val.(*lua.LTable); ok {
		// if we've already made a proxy, reuse it
		if proxy, found := proxies[tbl]; found {
			return proxy
		}

		// create the proxy table
		proxy := L.NewTable()
		proxies[tbl] = proxy

		mt := L.NewTable()

		// __index: read from original and wrap recursively
		L.SetField(mt, "__index", L.NewFunction(func(L2 *lua.LState) int {
			key := L2.CheckAny(2)
			raw := tbl.RawGet(key)
			L2.Push(ReadOnlyValue(L2, raw))
			return 1
		}))

		// __newindex: forbid writes
		L.SetField(mt, "__newindex", L.NewFunction(func(L2 *lua.LState) int {
			L2.RaiseError("table is readonly")
			return 0
		}))

		// __pairs: allow iteration over the original table
		L.SetField(mt, "__pairs", L.NewFunction(func(L2 *lua.LState) int {
			// call next on the original table to iterate
			L2.Push(L2.GetGlobal("next")) // `next` function
			L2.Push(tbl)                  // original table
			L2.Push(lua.LNil)             // starting point for `next`
			return 3                      // return next, table, nil (or key, value)
		}))

		L.SetMetatable(proxy, mt)
		return proxy
	}
	// non-table: return as-is
	return val
}
