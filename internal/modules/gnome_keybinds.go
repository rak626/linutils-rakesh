package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rakesh/linutils-rakesh/internal/config"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/tui"
)

func RunInteractiveGnomeKeybinds() error {
	items := []tui.ListItem{
		{Key: "default", Name: "Default Settings", Description: "Apply variables from your ~/.config/linutils/variables.conf file."},
		{Key: "customize", Name: "Customize Settings", Description: "Interactively choose preferred apps for terminal, browser, editor, etc."},
	}

	action, results, err := tui.RunListUIWithDesc("GNOME Keybindings", "Choose how to configure your shortcuts.", items)
	if err != nil {
		return err
	}

	if action == "" || action == "back" { // User quit or pressed esc
		return nil
	}

	var choice string
	for _, item := range results {
		if item.Selected {
			choice = item.Key
			break
		}
	}

	switch choice {
	case "customize":
		return runCustomizationFlow()
	case "default":
		return SetupGnomeKeybinds()
	}
	return nil
}

func runCustomizationFlow() error {
	vars := []string{"$terminal", "$browser", "$filemanager", "$editor", "$launcher"}
	varNames := map[string]string{
		"$terminal":    "Terminal Emulator",
		"$browser":     "Web Browser",
		"$filemanager": "File Manager",
		"$editor":      "Code/Text Editor",
		"$launcher":    "App Launcher",
	}

	for _, v := range vars {
		var items []tui.ListItem
		for _, opt := range config.VarOptions[v] {
			items = append(items, tui.ListItem{
				Key:         opt,
				Name:        opt,
				Description: "Select " + opt + " as your " + varNames[v],
			})
		}

		action, results, err := tui.RunListUIWithDesc("Select "+varNames[v], "Choose your preferred application.", items)
		if err != nil {
			return err
		}
		if action == "" || action == "back" {
			return nil // User cancelled
		}

		for _, item := range results {
			if item.Selected {
				config.UserVars[v] = item.Key
				break
			}
		}
	}

	config.SaveVariables()
	fmt.Println("\nVariables saved! Applying GNOME keybindings...")
	return SetupGnomeKeybinds()
}

func SetupGnomeKeybinds() error {
	if !pkgmanager.IsCommandAvailable("gsettings") {
		return fmt.Errorf("gsettings command not found. This module only works on GNOME")
	}

	fmt.Println("\n--- Setting up GNOME Keybindings ---")

	// 1. Fixed Workspaces
	fmt.Println("Configuring 9 fixed workspaces...")
	runGsettings("set", "org.gnome.mutter", "dynamic-workspaces", "false")
	runGsettings("set", "org.gnome.desktop.wm.preferences", "num-workspaces", "9")

	// 2. Remove existing Super+Number app bindings first to avoid conflict resolution races
	for i := 1; i <= 9; i++ {
		si := strconv.Itoa(i)
		runGsettings("set", "org.gnome.shell.keybindings", "switch-to-application-"+si, "@as []")
	}
	// 3. Clear existing workspace bindings
	for i := 1; i <= 9; i++ {
		si := strconv.Itoa(i)
		runGsettings("set", "org.gnome.desktop.wm.keybindings", "switch-to-workspace-"+si, "@as []")
		runGsettings("set", "org.gnome.desktop.wm.keybindings", "move-to-workspace-"+si, "@as []")
	}

	// 4. Set the new workspace bindings (Super+1..9 for switch, Super+Shift+1..9 for move window)
	for i := 1; i <= 9; i++ {
		si := strconv.Itoa(i)
		runGsettings("set", "org.gnome.desktop.wm.keybindings", "switch-to-workspace-"+si, "['<Super>"+si+"']")
		runGsettings("set", "org.gnome.desktop.wm.keybindings", "move-to-workspace-"+si, "['<Super><Shift>"+si+"']")
	}

	// 6. Window Management
	fmt.Println("Setting window management shortcuts...")
	runGsettings("set", "org.gnome.desktop.wm.keybindings", "close", "['<Super>q']")
	
	// Disable GNOME's default Super+D (Show Desktop) to allow it for Launcher (wofi/rofi)
	runGsettings("set", "org.gnome.desktop.wm.keybindings", "show-desktop", "@as []")

	// Remove Super+H keybind (which is to hide/minimize the window)
	runGsettings("set", "org.gnome.desktop.wm.keybindings", "minimize", "@as []")

	// Configure Alt+Tab to switch windows on current workspace only
	fmt.Println("Configuring Alt+Tab to switch windows on current workspace only...")
	runGsettings("set", "org.gnome.desktop.wm.keybindings", "switch-applications", "['<Super>Tab']")
	runGsettings("set", "org.gnome.desktop.wm.keybindings", "switch-applications-backward", "['<Shift><Super>Tab']")
	runGsettings("set", "org.gnome.desktop.wm.keybindings", "switch-windows", "['<Alt>Tab']")
	runGsettings("set", "org.gnome.desktop.wm.keybindings", "switch-windows-backward", "['<Shift><Alt>Tab']")
	runGsettings("set", "org.gnome.shell.window-switcher", "current-workspace-only", "true")

	// 7. Custom Shortcuts
	fmt.Println("Configuring custom app shortcuts...")

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}
	toggleScriptPath := filepath.Join(home, ".dotfiles", "scripts", "gnome", "toggle-panel.sh")

	customBinds := []string{
		"'/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/'",
		"'/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom1/'",
		"'/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom2/'",
		"'/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom3/'",
		"'/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom4/'",
		"'/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom5/'",
		"'/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom6/'",
		"'/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom7/'",
		"'/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom8/'",
	}

	runGsettings("set", "org.gnome.settings-daemon.plugins.media-keys", "custom-keybindings", "["+strings.Join(customBinds, ", ")+"]")

	setupCustomBind(0, "Terminal", config.ExpandVariables("$terminal"), "<Super>Return")
	setupCustomBind(1, "Browser", config.ExpandVariables("$browser"), "<Super><Shift>Return")
	setupCustomBind(2, "Files", config.ExpandVariables("$filemanager"), "<Super>e")
	setupCustomBind(3, "Editor", config.ExpandVariables("$editor"), "<Super><Shift>z")
	setupCustomBind(4, "Launcher", config.ExpandVariables("$launcher"), "<Super>d")
	setupCustomBind(5, "Toggle GNOME Panel", toggleScriptPath, "<Super>h")
	setupCustomBind(6, "Github Desktop", "github-desktop", "<Super><Shift>g")
	setupCustomBind(7, "Intellij Idea", "idea", "<Super><Shift>i")

	fmt.Println("GNOME keybindings setup complete.")
	return nil
}

func setupCustomBind(index int, name, command, binding string) {
	path := "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom" + strconv.Itoa(index) + "/"
	schema := "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:" + path

	runGsettings("set", schema, "name", "'"+name+"'")
	runGsettings("set", schema, "command", "'"+command+"'")
	runGsettings("set", schema, "binding", "'"+binding+"'")
}

func runGsettings(args ...string) {
	pkgmanager.RunCommand("gsettings", args...)
}
