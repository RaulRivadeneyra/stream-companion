package plugins_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RaulRivadeneyra/stream-companion/commerce/plugins"
	lua "github.com/yuin/gopher-lua"
)

func writePlugin(t *testing.T, root string, relPath string, content string) {
	fullPath := filepath.Join(root, relPath)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create plugin dir for %s: %v", relPath, err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write plugin file %s: %v", relPath, err)
	}
}

func loadPluginsOrFail(t *testing.T, root string) *lua.LTable {
	L := lua.NewState()
	t.Cleanup(L.Close)

	pluginsTable, err := plugins.LoadPlugins(root, L)
	if err != nil {
		t.Fatalf("LoadPlugins failed: %v", err)
	}

	return pluginsTable.(*lua.LTable)
}

func TestLoadsFlatPlugin(t *testing.T) {
	root := t.TempDir()
	writePlugin(t, root, "hello.lua", `return function() return "hi" end`)

	pluginsTable := loadPluginsOrFail(t, root)

	val := pluginsTable.RawGetString("hello")
	if val.Type() != lua.LTFunction {
		t.Errorf("Expected plugins.hello to be a function, got %s", val.Type().String())
	}
}

func TestLoadsNamespacedPlugin(t *testing.T) {
	root := t.TempDir()
	writePlugin(t, root, "twitch/send_message.lua", `return function(msg) return "sent: " .. msg end`)

	pluginsTable := loadPluginsOrFail(t, root)

	twitch := pluginsTable.RawGetString("twitch")
	if twitch.Type() != lua.LTTable {
		t.Fatalf("Expected plugins.twitch to be a table, got %s", twitch.Type().String())
	}

	send := twitch.(*lua.LTable).RawGetString("send_message")
	if send.Type() != lua.LTFunction {
		t.Errorf("Expected plugins.twitch.send_message to be a function, got %s", send.Type().String())
	}
}

func TestSkipsBrokenLuaFile(t *testing.T) {
	root := t.TempDir()
	writePlugin(t, root, "broken.lua", `return function(`) // bad syntax

	pluginsTable := loadPluginsOrFail(t, root)

	if val := pluginsTable.RawGetString("broken"); val.Type() != lua.LTNil {
		t.Errorf("Expected plugins.broken to be nil, got %s", val.Type().String())
	}
}

func TestSkipsPluginThatReturnsNonFunction(t *testing.T) {
	root := t.TempDir()
	writePlugin(t, root, "not_callable.lua", `return { foo = "bar" }`)

	pluginsTable := loadPluginsOrFail(t, root)

	if val := pluginsTable.RawGetString("not_callable"); val.Type() != lua.LTNil {
		t.Errorf("Expected plugins.not_callable to be nil, got %s", val.Type().String())
	}
}
