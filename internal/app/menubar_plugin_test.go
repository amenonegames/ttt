package app

import (
	"path/filepath"
	"testing"

	"github.com/eugenioenko/ttt/internal/plugin"
	"github.com/eugenioenko/ttt/internal/ui"
)

func TestMenuIndexForName(t *testing.T) {
	cases := []struct {
		name  string
		index int
		ok    bool
	}{
		{"file", 0, true},
		{"Edit", 1, true},
		{"SELECTION", 2, true},
		{"view", 3, true},
		{"help", 5, true},
		{"options", 0, false},
		{"tools", 0, false},
		{"", 0, false},
	}
	for _, c := range cases {
		index, ok := menuIndexForName(c.name)
		if index != c.index || ok != c.ok {
			t.Errorf("menuIndexForName(%q) = (%d, %v), want (%d, %v)", c.name, index, ok, c.index, c.ok)
		}
	}
}

func appWithMenuPlugins(t *testing.T, plugins ...*plugin.Plugin) *App {
	t.Helper()
	dir := t.TempDir()
	m := plugin.NewManager(filepath.Join(dir, "plugins"), filepath.Join(dir, "plugins.ttt.json"))
	for _, p := range plugins {
		m.RegisterDebugPlugin(p)
	}
	return &App{PluginManager: m}
}

func menuLabels(items []ui.ContextMenuItem) []string {
	labels := make([]string, len(items))
	for i, item := range items {
		if item.IsSep {
			labels[i] = "---"
			continue
		}
		labels[i] = item.Label
	}
	return labels
}

func TestPluginMenuItemsGroupsEachEnabledPluginAfterASeparator(t *testing.T) {
	a := appWithMenuPlugins(t,
		&plugin.Plugin{Name: "alpha", Enabled: true, MenuItems: []plugin.PluginMenuItem{
			{Menu: "file", Label: "Alpha One", Command: "alpha.one"},
			{Menu: "view", Label: "Alpha View", Command: "alpha.view"},
			{Menu: "file", Label: "Alpha Two", Command: "alpha.two"},
		}},
		&plugin.Plugin{Name: "disabled", Enabled: false, MenuItems: []plugin.PluginMenuItem{
			{Menu: "file", Label: "Disabled", Command: "disabled.run"},
		}},
		&plugin.Plugin{Name: "beta", Enabled: true, MenuItems: []plugin.PluginMenuItem{
			{Menu: "file", Label: "Beta", Command: "beta.run"},
		}},
		&plugin.Plugin{Name: "elsewhere", Enabled: true, MenuItems: []plugin.PluginMenuItem{
			{Menu: "help", Label: "Elsewhere", Command: "elsewhere.run"},
		}},
	)

	got := menuLabels(pluginMenuItems(a, 0))
	want := []string{"---", "Alpha One", "Alpha Two", "---", "Beta"}
	if len(got) != len(want) {
		t.Fatalf("file menu plugin items = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("file menu plugin items = %q, want %q", got, want)
		}
	}
}

func TestPluginMenuItemsSkipsDisabledPlugin(t *testing.T) {
	a := appWithMenuPlugins(t, &plugin.Plugin{Name: "disabled", Enabled: false, MenuItems: []plugin.PluginMenuItem{
		{Menu: "file", Label: "Disabled", Command: "disabled.run"},
	}})

	if items := pluginMenuItems(a, 0); len(items) != 0 {
		t.Errorf("disabled plugin contributed %q", menuLabels(items))
	}
}

func TestWithPluginMenuItemsLeavesBuiltInMenuUntouched(t *testing.T) {
	a := appWithMenuPlugins(t, &plugin.Plugin{Name: "alpha", Enabled: true, MenuItems: []plugin.PluginMenuItem{
		{Menu: "file", Label: "Alpha", Command: "alpha.run"},
	}})
	builtIn := menuBarMenus[0]
	builtInLen := len(builtIn)
	builtInLast := builtIn[builtInLen-1].Label

	combined := withPluginMenuItems(a, 0, builtIn)

	if len(menuBarMenus[0]) != builtInLen || menuBarMenus[0][builtInLen-1].Label != builtInLast {
		t.Fatalf("built-in File menu was mutated: %q", menuLabels(menuBarMenus[0]))
	}
	if len(combined) != builtInLen+2 {
		t.Fatalf("combined menu = %q, want built-in items plus a separator and one plugin item", menuLabels(combined))
	}
	if !combined[builtInLen].IsSep || combined[builtInLen+1].Label != "Alpha" {
		t.Errorf("plugin block = %q", menuLabels(combined[builtInLen:]))
	}
}
