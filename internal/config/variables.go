package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FontSizeBaseKey is the global default/fallback size.
const FontSizeBaseKey = "$font_size"

// FontSizeApps is the stable per-app order used for prompts + persistence.
// App IDs are internal (no "$", no spaces); each maps to $font_size_<app>.
var FontSizeApps = []string{
	"gnome",
	"gnome_terminal",
	"alacritty",
	"kitty",
	"foot",
	"rofi",
	"wofi",
	"waybar",
	"mako",
	"qt",
	"i3",
	"i3status",
	"polybar",
	"zed",
	"vscode",
}

// FontSizeKey returns the variables.conf key for an app, e.g. "$font_size_alacritty".
func FontSizeKey(app string) string {
	return FontSizeBaseKey + "_" + app
}

// FontSizeKeys returns all per-app size keys in stable order.
func FontSizeKeys() []string {
	keys := make([]string, 0, len(FontSizeApps))
	for _, app := range FontSizeApps {
		keys = append(keys, FontSizeKey(app))
	}
	return keys
}

// FontSizeFor returns the configured size for an app, falling back to
// $font_size (base) and finally 11 when unset/unparseable.
func FontSizeFor(app string) int {
	if v, ok := UserVars[FontSizeKey(app)]; ok {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n > 0 {
			return n
		}
	}
	if v, ok := UserVars[FontSizeBaseKey]; ok {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n > 0 {
			return n
		}
	}
	return 11
}

var UserVars = map[string]string{
	"$browser":                  "chromium-browser",
	"$secondary_browser":        "zen",
	"$terminal":                 "alacritty",
	"$editor":                   "zed",
	"$filemanager":              "nautilus",
	"$launcher":                 "wofi --show drun",
	"$font_mono":                "JetBrainsMono Nerd Font Mono",
	"$font_size":                "11",
	"$font_size_gnome":          "11",
	"$font_size_gnome_terminal": "11",
	"$font_size_alacritty":      "11",
	"$font_size_kitty":          "11",
	"$font_size_foot":           "11",
	"$font_size_rofi":           "11",
	"$font_size_wofi":           "11",
	"$font_size_waybar":         "11",
	"$font_size_mako":           "11",
	"$font_size_qt":             "11",
	"$font_size_i3":             "11",
	"$font_size_i3status":       "11",
	"$font_size_polybar":        "11",
	"$font_size_zed":            "11",
	"$font_size_vscode":         "11",
}

var VarOptions = map[string][]string{
	"$browser":           {"chromium-browser", "zen", "brave-browser", "firefox", "google-chrome-stable"},
	"$secondary_browser": {"zen", "chromium-browser", "brave-browser", "firefox", "google-chrome-stable"},
	"$terminal":          {"alacritty", "kitty", "ghostty", "gnome-terminal"},
	"$editor":            {"zed", "code", "nvim", "vim", "micro"},
	"$filemanager":       {"nautilus", "thunar", "dolphin"},
	"$launcher":          {"rofi -show drun", "wofi --show drun"},
}

func LoadVariables() {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "linutils")
	configPath := filepath.Join(configDir, "variables.conf")

	// Ensure directory exists
	os.MkdirAll(configDir, 0755)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		saveDefaultVariables(configPath)
		return
	}

	file, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Warning: failed to open %s: %v\n", configPath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if !strings.HasPrefix(key, "$") {
				key = "$" + key
			}
			UserVars[key] = val
		}
	}

	// Backfill: old configs only have $font_size. Any missing per-app
	// size inherits the base so behavior stays identical until changed.
	base := strings.TrimSpace(UserVars[FontSizeBaseKey])
	if base == "" {
		base = "11"
		UserVars[FontSizeBaseKey] = base
	}
	for _, app := range FontSizeApps {
		k := FontSizeKey(app)
		if strings.TrimSpace(UserVars[k]) == "" {
			UserVars[k] = base
		}
	}
}

func saveDefaultVariables(path string) {
	content := "# Linutils Variable Configuration\n"
	content += "# Format: $variable = value\n\n"
	content += "$browser = chromium-browser\n"
	content += "$secondary_browser = zen\n"
	content += "$terminal = alacritty\n"
	content += "$editor = zed\n"
	content += "$filemanager = nautilus\n"
	content += "$launcher = wofi --show drun\n"
	content += "$font_mono = JetBrainsMono Nerd Font Mono\n"
	content += "$font_size = 11\n"
	for _, app := range FontSizeApps {
		content += fmt.Sprintf("%s = 11\n", FontSizeKey(app))
	}

	os.WriteFile(path, []byte(content), 0644)
	fmt.Printf("Created default variables config at %s\n", path)
}

func SaveVariables() {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "linutils", "variables.conf")

	content := "# Linutils Variable Configuration\n"
	content += "# Format: $variable = value\n\n"

	// Order keys for consistency
	keys := []string{"$browser", "$secondary_browser", "$terminal", "$editor", "$filemanager", "$launcher", "$font_mono", "$font_size"}
	keys = append(keys, FontSizeKeys()...)
	for _, k := range keys {
		content += fmt.Sprintf("%s = %s\n", k, UserVars[k])
	}

	os.WriteFile(configPath, []byte(content), 0644)
}

func ExpandVariables(input string) string {
	output := input
	for key, val := range UserVars {
		output = strings.ReplaceAll(output, key, val)
	}
	return output
}
