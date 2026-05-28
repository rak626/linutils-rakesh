package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/config"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

func RunStandaloneThemeSwitcher(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	var selectedThemeName string
	
	// Prepare options from internal config
	options := make([]huh.Option[string], len(config.GlobalThemes))
	for i, t := range config.GlobalThemes {
		options[i] = huh.NewOption(t.Name, t.Name)
	}

	// Add option for community themes
	options = append(options, huh.NewOption("--- Community Themes ---", "community_header"))
	options = append(options, huh.NewOption("Import Community Theme (GitHub URL)", "import_community"))

	// Load local community themes from ~/.config/linutils/themes/
	home, _ := os.UserHomeDir()
	commDir := filepath.Join(home, ".config", "linutils", "themes")
	if entries, err := os.ReadDir(commDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				options = append(options, huh.NewOption(entry.Name(), "comm:"+entry.Name()))
			}
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Global Theme Switcher").
				Description(fmt.Sprintf("Select a theme to apply for your current environment (%s).", sysInfo.DE)).
				Options(options...).
				Value(&selectedThemeName),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if selectedThemeName == "community_header" {
		return RunStandaloneThemeSwitcher(manager, sysInfo) // Refresh
	}

	if selectedThemeName == "import_community" {
		return handleImportCommunityTheme()
	}

	var selectedTheme config.ThemeConfig
	if len(selectedThemeName) > 5 && selectedThemeName[:5] == "comm:" {
		// Handle community theme from disk
		themePath := filepath.Join(commDir, selectedThemeName[5:], "theme.json")
		data, err := os.ReadFile(themePath)
		if err != nil {
			return fmt.Errorf("failed to read community theme: %v", err)
		}
		if err := json.Unmarshal(data, &selectedTheme); err != nil {
			return fmt.Errorf("failed to parse community theme: %v", err)
		}
	} else {
		// Use internal theme
		for _, t := range config.GlobalThemes {
			if t.Name == selectedThemeName {
				selectedTheme = t
				break
			}
		}
	}

	return ApplyGlobalTheme(selectedTheme, sysInfo)
}

func handleImportCommunityTheme() error {
	var url string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter GitHub Theme URL").
				Placeholder("https://github.com/user/omarchy-theme").
				Value(&url),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if url == "" {
		return nil
	}

	home, _ := os.UserHomeDir()
	commDir := filepath.Join(home, ".config", "linutils", "themes")
	os.MkdirAll(commDir, 0755)

	// Extract folder name from URL
	parts := regexp.MustCompile(`/`).Split(url, -1)
	folderName := parts[len(parts)-1]
	dest := filepath.Join(commDir, folderName)

	fmt.Printf("Cloning theme from %s to %s...\n", url, dest)
	cmd := exec.Command("git", "clone", "--depth", "1", url, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone theme: %v", err)
	}

	fmt.Println("Theme imported! You can now select it from the menu.")
	return nil
}

func ApplyGlobalTheme(theme config.ThemeConfig, sysInfo system.Info) error {
	fmt.Printf("Applying theme: %s\n", theme.Name)

	home, _ := os.UserHomeDir()

	// 1. Universal Apps
	updateAlacritty(home, theme.Alacritty)
	updateZed(home, theme.Zed)
	updateNeovim(home, theme.Neovim)
	updateVim(home, theme.Vim)
	updateGTK(theme.GTK)
	updateStarship(home, theme.Starship)
	updateVSCodium(home, theme.VSCodium)
	updateUlauncher(home, theme.Ulauncher)
	
	// New Universal Apps
	updateGhostty(home, theme.Ghostty)
	updateBtop(home, theme.Btop)
	updateKitty(home, theme.Kitty)
	updateIcons(theme.Icons, theme.Cursor)

	// 2. Environment Specific Logic
	switch sysInfo.DE {
	case "gnome":
		updateGNOME(theme.GTK, theme.GNOME.ShellTheme, theme.GNOME.Wallpaper)
	case "hyprland":
		updateHyprland(home, theme.Hyprland)
		updateWaybar(home, theme.Waybar)
		updateMako(home, theme.Mako)
		updateHyprlock(home, theme.Hyprlock)
		updateSwayOSD(home, theme.SwayOSD)
	case "i3":
		updateI3(home, theme.I3)
	default:
		fmt.Printf("Note: No specific UI tweaks for detected environment: %s\n", sysInfo.DE)
	}

	fmt.Println("Global theme applied successfully!")
	return nil
}

func updateAlacritty(home, themeFile string) {
	if themeFile == "" { return }
	path := filepath.Join(home, ".config", "alacritty", "alacritty.toml")
	content, err := os.ReadFile(path)
	if err != nil { return }

	re := regexp.MustCompile(`import\s*=\s*\[\s*".*/alacritty/.*\.toml"\s*\]`)
	newImport := fmt.Sprintf(`import = ["~/.config/alacritty/%s"]`, themeFile)
	newContent := re.ReplaceAllString(string(content), newImport)

	os.WriteFile(path, []byte(newContent), 0644)
	fmt.Println("  - Alacritty updated.")
}

func updateZed(home, themeName string) {
	if themeName == "" { return }
	path := filepath.Join(home, ".config", "zed", "settings.json")
	content, err := os.ReadFile(path)
	if err != nil { return }

	var settings map[string]interface{}
	if err := json.Unmarshal(content, &settings); err != nil { return }

	settings["theme"] = themeName
	newContent, _ := json.MarshalIndent(settings, "", "  ")
	os.WriteFile(path, newContent, 0644)
	fmt.Println("  - Zed updated.")
}

func updateNeovim(home, themeName string) {
	if themeName == "" { return }
	dir := filepath.Join(home, ".config", "nvim", "lua")
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "active_theme.lua")
	content := fmt.Sprintf("vim.cmd.colorscheme(\"%s\")\n", themeName)
	os.WriteFile(path, []byte(content), 0644)
	fmt.Println("  - Neovim active_theme.lua updated.")
}

func updateVim(home, themeName string) {
	if themeName == "" { return }
	path := filepath.Join(home, ".vim", "active_theme.vim")
	os.MkdirAll(filepath.Dir(path), 0755)
	content := fmt.Sprintf("colorscheme %s\n", themeName)
	os.WriteFile(path, []byte(content), 0644)
	fmt.Println("  - Vim active_theme.vim updated.")
}

func updateGTK(themeName string) {
	if themeName == "" { return }
	exec.Command("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", themeName).Run()
	exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", "prefer-dark").Run()
	fmt.Println("  - GTK theme set via gsettings.")
}

func updateGNOME(gtkTheme, shellTheme, wallpaper string) {
	if gtkTheme != "" {
		exec.Command("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", gtkTheme).Run()
	}
	if shellTheme != "" {
		exec.Command("gsettings", "set", "org.gnome.shell.extensions.user-theme", "name", shellTheme).Run()
	}
	if wallpaper != "" {
		home, _ := os.UserHomeDir()
		wallPath := filepath.Join(home, "Pictures", "Wallpapers", wallpaper)
		if _, err := os.Stat(wallPath); err == nil {
			uri := "file://" + wallPath
			exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", uri).Run()
			exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", uri).Run()
			fmt.Println("  - GNOME Wallpaper updated.")
		}
	}
}

func updateStarship(home, themeName string) {
	if themeName == "" { return }
	path := filepath.Join(home, ".config", "starship.toml")
	content, err := os.ReadFile(path)
	if err != nil { return }
	re := regexp.MustCompile(`palette\s*=\s*".*"`)
	newPalette := fmt.Sprintf(`palette = "%s"`, themeName)
	newContent := re.ReplaceAllString(string(content), newPalette)
	os.WriteFile(path, []byte(newContent), 0644)
	fmt.Println("  - Starship updated.")
}

func updateVSCodium(home, themeName string) {
	if themeName == "" { return }
	path := filepath.Join(home, ".config", "VSCodium", "User", "settings.json")
	content, err := os.ReadFile(path)
	if err != nil { return }
	var settings map[string]interface{}
	if err := json.Unmarshal(content, &settings); err != nil { return }
	settings["workbench.colorTheme"] = themeName
	newContent, _ := json.MarshalIndent(settings, "", "  ")
	os.WriteFile(path, newContent, 0644)
	fmt.Println("  - VSCodium updated.")
}

func updateGhostty(home, themeName string) {
	if themeName == "" { return }
	path := filepath.Join(home, ".config", "ghostty", "config")
	content, err := os.ReadFile(path)
	if err != nil { return }
	re := regexp.MustCompile(`theme\s*=\s*.*`)
	newTheme := fmt.Sprintf(`theme = %s`, themeName)
	newContent := re.ReplaceAllString(string(content), newTheme)
	os.WriteFile(path, []byte(newContent), 0644)
	fmt.Println("  - Ghostty updated.")
}

func updateBtop(home, themeName string) {
	if themeName == "" { return }
	path := filepath.Join(home, ".config", "btop", "btop.conf")
	content, err := os.ReadFile(path)
	if err != nil { return }
	re := regexp.MustCompile(`color_theme\s*=\s*".*"`)
	newTheme := fmt.Sprintf(`color_theme = "%s"`, themeName)
	newContent := re.ReplaceAllString(string(content), newTheme)
	os.WriteFile(path, []byte(newContent), 0644)
	fmt.Println("  - btop updated.")
}

func updateKitty(home, themeName string) {
	if themeName == "" { return }
	// Kitty themes are usually just a file we can overwrite with a 'include' if we have them
	// but for now let's assume we use kitty-themes or manual snippets.
	// A simpler way is to use the kitty +kitten themes command
	exec.Command("kitty", "+kitten", "themes", "--reload-in=all", themeName).Run()
	fmt.Println("  - Kitty updated via kitten.")
}

func updateIcons(iconTheme, cursorTheme string) {
	if iconTheme != "" {
		exec.Command("gsettings", "set", "org.gnome.desktop.interface", "icon-theme", iconTheme).Run()
		fmt.Println("  - Icon theme set via gsettings.")
	}
	if cursorTheme != "" {
		exec.Command("gsettings", "set", "org.gnome.desktop.interface", "cursor-theme", cursorTheme).Run()
		fmt.Println("  - Cursor theme set via gsettings.")
	}
}

func updateWaybar(home, cssVars string) {
	if cssVars == "" { return }
	path := filepath.Join(home, ".config", "waybar", "active_theme.css")
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(cssVars), 0644)
	// Waybar usually reloads on config change, but let's signal it
	exec.Command("pkill", "-USR2", "waybar").Run()
	fmt.Println("  - Waybar active_theme.css updated.")
}

func updateMako(home, configContent string) {
	if configContent == "" { return }
	// We'll generate a snippet and mako doesn't support 'source' easily in a single file
	// but we can append it or use a separate file if user has configured it.
	// For now, let's just write to a known location
	path := filepath.Join(home, ".config", "mako", "active_theme")
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(configContent), 0644)
	exec.Command("makoctl", "reload").Run()
	fmt.Println("  - Mako reloaded with new theme.")
}

func updateHyprlock(home, configContent string) {
	if configContent == "" { return }
	path := filepath.Join(home, ".config", "hypr", "active_theme_lock.conf")
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(configContent), 0644)
	fmt.Println("  - Hyprlock theme updated.")
}

func updateSwayOSD(home, cssContent string) {
	if cssContent == "" { return }
	path := filepath.Join(home, ".config", "swayosd", "style.css") // Assuming user uses this path
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(cssContent), 0644)
	fmt.Println("  - SwayOSD theme updated.")
}

func updateHyprland(home, configContent string) {
	if configContent == "" { return }
	path := filepath.Join(home, ".config", "hypr", "active_theme.conf")
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(configContent), 0644)
	fmt.Println("  - Hyprland active_theme.conf updated.")
}

func updateI3(home, configContent string) {
	if configContent == "" { return }
	path := filepath.Join(home, ".config", "i3", "active_theme.i3")
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(configContent), 0644)
	exec.Command("i3-msg", "reload").Run()
	fmt.Println("  - i3 active_theme.i3 updated.")
}

func updateUlauncher(home, themeName string) {
	if themeName == "" { return }
	path := filepath.Join(home, ".config", "ulauncher", "settings.json")
	content, err := os.ReadFile(path)
	if err != nil { return }
	var settings map[string]interface{}
	if err := json.Unmarshal(content, &settings); err != nil { return }
	settings["theme-name"] = themeName
	newContent, _ := json.MarshalIndent(settings, "", "  ")
	os.WriteFile(path, newContent, 0644)
	fmt.Println("  - Ulauncher updated.")
}

func InstallThemeSwitcher(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Installing Global Theme Switcher ---")

	home, _ := os.UserHomeDir()
	binDir := filepath.Join(home, ".local", "bin")
	os.MkdirAll(binDir, 0755)

	exe, err := os.Executable()
	if err != nil { return err }

	destExe := filepath.Join(binDir, "linutils-rakesh")
	if err := copyFile(exe, destExe); err != nil { return err }
	os.Chmod(destExe, 0755)

	scriptPath := filepath.Join(binDir, "theme-switcher")
	scriptContent := fmt.Sprintf("#!/bin/bash\n%s theme\n", destExe)
	os.WriteFile(scriptPath, []byte(scriptContent), 0755)

	fmt.Println("Theme switcher installed to ~/.local/bin/theme-switcher")

	sysInfo := system.GetSystemInfo()
	switch sysInfo.DE {
	case "hyprland":
		hyprConfig := filepath.Join(home, ".config", "hypr", "hyprland.conf")
		if _, err := os.Stat(hyprConfig); err == nil {
			keybind := "\nbind = $mainMod ALT, T, exec, kitty --class floating -e theme-switcher\n"
			appendToFileIfMissing(hyprConfig, keybind)
			fmt.Println("Added keybind to hyprland.conf: $mainMod + ALT + T")
		}
	case "i3":
		i3Config := filepath.Join(home, ".config", "i3", "config")
		if _, err := os.Stat(i3Config); err == nil {
			keybind := "\nbindsym $mod+Mod1+t exec kitty --class floating -e theme-switcher\n"
			appendToFileIfMissing(i3Config, keybind)
			fmt.Println("Added keybind to i3 config: $mod + Alt + T")
		}
	case "gnome":
		fmt.Println("For GNOME, please add a Custom Shortcut in Settings -> Keyboard:")
		fmt.Println("  Command: theme-switcher")
		fmt.Println("  Shortcut: Super+Alt+T")
	}

	return nil
}
