package app

import (
	"testing"

	"github.com/eugenioenko/ttt/internal/command"
	"github.com/eugenioenko/ttt/internal/plugin"
	lua "github.com/yuin/gopher-lua"
)

func registerPanelPlugin(t *testing.T, p *plugin.Plugin) *App {
	t.Helper()
	a := appWithMenuPlugins(t, p)
	a.Reg = command.NewRegistry()
	a.registerPluginCommandsAndKeys(p)
	return a
}

func TestPluginPanelCommandsAreRegisteredForEachPanel(t *testing.T) {
	a := registerPanelPlugin(t, &plugin.Plugin{
		Name:             "find-type",
		SidebarTitle:     "Types",
		RenderFunc:       &lua.LFunction{},
		BottomTitle:      "Results",
		BottomRenderFunc: &lua.LFunction{},
	})

	for id, title := range map[string]string{
		"sidebar.plugin.find-type": "Show Types",
		"panel.plugin.find-type":   "Show Results",
	} {
		cmd, ok := a.Reg.Get(id)
		if !ok || cmd.Title != title {
			t.Errorf("Get(%q) = (%q, %v), want (%q, true)", id, cmd.Title, ok, title)
		}
	}
}

func TestPluginPanelCommandIsNotRegisteredForATitleWithoutRender(t *testing.T) {
	a := registerPanelPlugin(t, &plugin.Plugin{Name: "find-type", SidebarTitle: "Types", BottomTitle: "Results"})

	for _, id := range []string{"sidebar.plugin.find-type", "panel.plugin.find-type"} {
		if _, ok := a.Reg.Get(id); ok {
			t.Errorf("%s registered without a render function", id)
		}
	}
}

func TestPluginPanelCommandIsSkippedWhenThePluginDeclaresTheSameID(t *testing.T) {
	a := registerPanelPlugin(t, &plugin.Plugin{
		Name:         "find-type",
		SidebarTitle: "Types",
		RenderFunc:   &lua.LFunction{},
		Commands: []plugin.PluginCommand{
			{ID: "sidebar.plugin.find-type", Title: "Find Type: Own", Handler: func() error { return nil }},
		},
	})

	if cmd, _ := a.Reg.Get("sidebar.plugin.find-type"); cmd.Title != "Find Type: Own" {
		t.Errorf("title = %q, want the plugin's own command", cmd.Title)
	}
}
