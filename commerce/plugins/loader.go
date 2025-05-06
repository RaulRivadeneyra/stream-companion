package plugins

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func LoadPlugins(root string, luaState *lua.LState) (lua.LValue, error) {
	pluginsTable := luaState.NewTable()

	err := filepath.WalkDir(root, func(path string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error walking path %s: %v", path, err)
		}
		if dirEntry.IsDir() || !strings.HasSuffix(dirEntry.Name(), ".lua") {
			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			log.Printf("Failed to get realtive path for %s: %v", path, err)
			return nil
		}

		compiledFn, err := luaState.LoadFile(path)

		if err != nil {
			log.Printf("Syntax error in %s: %v", path, err)
			return nil
		}

		luaState.Push(compiledFn)
		err = luaState.PCall(0, 1, nil) // execute file and capture returned value
		if err != nil {
			log.Printf("Runtime error in %s: %v", path, err)
			return nil
		}

		plugin := luaState.Get(-1)
		luaState.Pop(1)
		if plugin.Type() != lua.LTFunction {
			log.Printf("Plugin %s does not return a function", path)
			return nil
		}
		parts := strings.Split(relativePath, string(os.PathSeparator))
		current := pluginsTable

		for i := range len(parts) - 1 {
			part := parts[i]
			if sub := current.RawGetString(part); sub.Type() == lua.LTTable {
				current = sub.(*lua.LTable)
			} else {
				newTbl := luaState.NewTable()
				current.RawSetString(part, newTbl)
				current = newTbl
			}
		}

		pluginName := strings.TrimSuffix(parts[len(parts)-1], ".lua")
		current.RawSetString(pluginName, plugin)

		log.Printf("Validated & registered plugin: plugins.%s", strings.Join(parts, "."))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return pluginsTable, nil
}
