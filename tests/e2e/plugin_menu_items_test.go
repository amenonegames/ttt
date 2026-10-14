package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eugenioenko/ttt/internal/app"
	"github.com/eugenioenko/ttt/internal/plugin"
)

func TestPluginMenuItemAppearsInFileMenuWithShortcut(t *testing.T) {
	h := newTestHarness(t, 100, 30)
	defer h.stop()
	cfg := filepath.Join(h.dir, "config")
	h.app.PluginManager = plugin.NewManager(filepath.Join(cfg, "plugins"), filepath.Join(cfg, "plugins.ttt.json"))
	src := filepath.Join(h.dir, "findtype.lua")
	if err := os.WriteFile(src, []byte(`
		local ttt = require("ttt")
		ran = false
		ttt.register({
			commands = { { id = "find-type.open", title = "Find Type: Open", handler = function() ran = true end } },
			keybindings = { { key = "ctrl+k g", command = "find-type.open" } },
			menu_items = {
				{ menu = "File", label = "Find Type", command = "find-type.open" },
				{ menu = "options", label = "Ignored Options Item", command = "find-type.open" },
				{ menu = "view", command = "find-type.open" },
			},
		})
	`), 0644); err != nil {
		t.Fatal(err)
	}
	app.LoadPluginFromFile(h.app, src)
	p := h.app.PluginManager.FindPlugin("findtype")
	if p == nil {
		t.Fatal("plugin not loaded")
	}
	defer p.Destroy()

	h.exec("menu.file")
	h.assertContains("Find Type")
	h.assertContains("Ctrl+K G")
	h.assertContains("Quit")

	h.exec("menu.options")
	h.assertNotContains("Ignored Options Item")
}
