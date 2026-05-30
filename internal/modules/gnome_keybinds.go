package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

func SetupGnomeKeybinds() error {
	if !pkgmanager.IsCommandAvailable("gsettings") {
		return fmt.Errorf("gsettings command not found. This module only works on GNOME")
	}

	fmt.Println("\n--- Setting up GNOME Keybindings ---")

	// 1. Fixed Workspaces
	fmt.Println("Configuring 9 fixed workspaces...")
	runGsettings("set", "org.gnome.mutter", "dynamic-workspaces", "false")
	runGsettings("set", "org.gnome.desktop.wm.preferences", "num-workspaces", "9")

	// 2. Remove existing Super+Number app bindings
	// 3. Clear workspace bindings
	// 4. Super+1..9 => switch workspace
	// 5. Super+Shift+1..9 => move window
	for i := 1; i <= 9; i++ {
		si := strconv.Itoa(i)
		runGsettings("set", "org.gnome.shell.keybindings", "switch-to-application-"+si, "[]")
		runGsettings("set", "org.gnome.desktop.wm.keybindings", "switch-to-workspace-"+si, "[]")
		runGsettings("set", "org.gnome.desktop.wm.keybindings", "move-to-workspace-"+si, "[]")

		runGsettings("set", "org.gnome.desktop.wm.keybindings", "switch-to-workspace-"+si, "['<Super>"+si+"']")
		runGsettings("set", "org.gnome.desktop.wm.keybindings", "move-to-workspace-"+si, "['<Super><Shift>"+si+"']")
	}

	// 6. Window Management
	fmt.Println("Setting window management shortcuts...")
	runGsettings("set", "org.gnome.desktop.wm.keybindings", "close", "['<Super>q']")

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

	setupCustomBind(0, "Alacritty", "alacritty", "<Super>Return")
	setupCustomBind(1, "Chromium", "chromium-browser --new-window", "<Super><Shift>Return")
	setupCustomBind(2, "Files", "nautilus", "<Super>e")
	setupCustomBind(3, "Zed", "zed", "<Super><Shift>z")
	setupCustomBind(4, "Brave", "brave-browser --new-window", "<Super><Shift>b")
	setupCustomBind(5, "Ulauncher", "ulauncher-toggle", "<Super>d")
	setupCustomBind(6, "Toggle GNOME Panel", toggleScriptPath, "<Super>h")
	setupCustomBind(7, "Github Desktop", "github-desktop", "<Super><Shift>g")
	setupCustomBind(8, "Intellij Idea", "idea", "<Super><Shift>i")

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
