package plugin

import "testing"

func TestMenuItemRegistration(t *testing.T) {
	p, cleanup := newTestPluginBase(PermissionSet{Commands: true})
	defer cleanup()

	err := p.State.DoString(`
		local ttt = require("ttt")
		ttt.register({
			menu_items = {
				{ menu = "File", label = "Find Type", command = "find-type.open" },
				{ menu = "file", label = 42, command = "find-type.open" },
				{ menu = "file", label = "", command = "find-type.open" },
				{ menu = "file", label = "No Command", command = {} },
				{ menu = "file", label = "Empty Command", command = "" },
				{ label = "No Menu", command = "find-type.open" },
				{ menu = "help", label = "Find Type Help", command = "find-type.help" },
			},
		})
	`)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	want := []PluginMenuItem{
		{Menu: "file", Label: "Find Type", Command: "find-type.open"},
		{Menu: "help", Label: "Find Type Help", Command: "find-type.help"},
	}
	if len(p.MenuItems) != len(want) {
		t.Fatalf("expected %d menu items, got %+v", len(want), p.MenuItems)
	}
	for i, item := range want {
		if p.MenuItems[i] != item {
			t.Errorf("menu item %d = %+v, want %+v", i, p.MenuItems[i], item)
		}
	}
}

func TestMenuItemRegistrationWithoutPermission(t *testing.T) {
	p, cleanup := newTestPluginBase(PermissionSet{})
	defer cleanup()

	err := p.State.DoString(`
		local ttt = require("ttt")
		ttt.register({
			menu_items = {
				{ menu = "file", label = "Find Type", command = "find-type.open" },
			},
		})
	`)
	if err == nil {
		t.Fatal("expected error when commands permission not granted")
	}
	if len(p.MenuItems) != 0 {
		t.Errorf("expected 0 menu items, got %d", len(p.MenuItems))
	}
}
