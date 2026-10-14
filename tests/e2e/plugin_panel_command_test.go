package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eugenioenko/ttt/internal/app"
	"github.com/eugenioenko/ttt/internal/plugin"
	"github.com/eugenioenko/ttt/internal/ui"
)

const findTypePluginSource = `
	local ttt = require("ttt")
	typed = ""
	ttt.register({
		sidebar = {
			title = "Types",
			render = function(panel)
				panel:input({ placeholder = "Type name", on_change = function(text) typed = text end })
			end,
		},
		bottom = {
			title = "Results",
			render = function(panel)
				panel:input({ placeholder = "Filter", on_change = function(text) typed = text end })
			end,
		},
		keybindings = { { key = "ctrl+k t", command = "sidebar.plugin.find-type" } },
		menu_items = { { menu = "view", label = "Find Type", command = "sidebar.plugin.find-type" } },
	})
`

func loadFindTypePlugin(t *testing.T, h *testHarness) string {
	t.Helper()
	cfg := filepath.Join(h.dir, "config")
	h.app.PluginManager = plugin.NewManager(filepath.Join(cfg, "plugins"), filepath.Join(cfg, "plugins.ttt.json"))
	src := filepath.Join(h.dir, "find-type.lua")
	if err := os.WriteFile(src, []byte(findTypePluginSource), 0644); err != nil {
		t.Fatal(err)
	}
	app.LoadPluginFromFile(h.app, src)
	if h.app.PluginManager.FindPlugin("find-type") == nil {
		t.Fatal("plugin not loaded")
	}
	t.Cleanup(func() { h.app.PluginManager.FindPlugin("find-type").Destroy() })
	return src
}

func (h *testHarness) assertTypedIntoPlugin(text string) {
	h.t.Helper()
	h.typeText(text)
	p := h.app.PluginManager.FindPlugin("find-type")
	if got := p.State.GetGlobal("typed").String(); got != text {
		h.t.Errorf("plugin input text = %q, want %q", got, text)
	}
}

func TestPluginSidebarCommandShowsPanelAndFocusesItsInput(t *testing.T) {
	h := newTestHarness(t, 100, 30)
	defer h.stop()
	src := loadFindTypePlugin(t, h)

	cmd, ok := h.reg.Get("sidebar.plugin.find-type")
	if !ok {
		t.Fatal("sidebar.plugin.find-type not registered")
	}
	if cmd.Title != "Show Types" || cmd.Shortcut != "Ctrl+K T" {
		t.Errorf("command = (%q, %q), want (%q, %q)", cmd.Title, cmd.Shortcut, "Show Types", "Ctrl+K T")
	}

	h.exec("sidebar.explorer")
	h.app.HideSidebar()
	h.app.FocusEditor()
	h.redraw()

	h.exec("sidebar.plugin.find-type")
	if !h.app.Sidebar.Visible {
		t.Fatal("sidebar is hidden")
	}
	if h.app.Sidebar.ActivePanel != "plugin.find-type" {
		t.Fatalf("active sidebar panel = %q, want plugin.find-type", h.app.Sidebar.ActivePanel)
	}
	if h.app.Root.Focused != h.app.Sidebar.PanelWidget("plugin.find-type") {
		t.Fatalf("focus = %T, want the plugin panel", h.app.Root.Focused)
	}
	h.assertTypedIntoPlugin("Foo")

	stale := h.app.Sidebar.PanelWidget("plugin.find-type")
	h.app.FocusEditor()
	app.LoadPluginFromFile(h.app, src)
	h.redraw()
	live := h.app.Sidebar.PanelWidget("plugin.find-type")
	if live == stale {
		t.Fatal("reload did not replace the panel widget")
	}

	h.exec("sidebar.plugin.find-type")
	if h.app.Root.Focused != live {
		t.Fatalf("focus after reload = %T, want the replaced plugin panel", h.app.Root.Focused)
	}
	h.assertTypedIntoPlugin("Bar")
}

func TestPluginBottomCommandShowsPanelAndFocusesItsInput(t *testing.T) {
	h := newTestHarness(t, 100, 30)
	defer h.stop()
	loadFindTypePlugin(t, h)

	if cmd, _ := h.reg.Get("panel.plugin.find-type"); cmd.Title != "Show Results" {
		t.Errorf("title = %q, want %q", cmd.Title, "Show Results")
	}

	h.app.HideBottomPanel()
	h.app.ContentSplit.BottomH = 0
	h.redraw()

	h.exec("panel.plugin.find-type")
	if !h.app.ContentSplit.ShowBottom {
		t.Fatal("bottom panel is hidden")
	}
	if h.app.ContentSplit.BottomH < ui.MinPanelHeight {
		t.Errorf("bottom panel height = %d, want at least %d", h.app.ContentSplit.BottomH, ui.MinPanelHeight)
	}
	if h.app.BottomPanel.ActivePanel != "plugin.find-type" {
		t.Fatalf("active bottom panel = %q, want plugin.find-type", h.app.BottomPanel.ActivePanel)
	}
	if h.app.Root.Focused != h.app.BottomPanel.PanelWidget("plugin.find-type") {
		t.Fatalf("focus = %T, want the plugin panel", h.app.Root.Focused)
	}
	h.assertTypedIntoPlugin("Baz")
}
