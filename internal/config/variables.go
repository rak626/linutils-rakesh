package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var UserVars = map[string]string{
	"$browser":     "chromium-browser --new-window",
	"$terminal":    "alacritty",
	"$editor":      "zed",
	"$filemanager": "nautilus",
	"$launcher":    "wofi --show drun",
}

var VarOptions = map[string][]string{
	"$browser":     {"zen", "chromium-browser --new-window", "brave-browser --new-window", "firefox", "google-chrome-stable"},
	"$terminal":    {"alacritty", "kitty", "ghostty", "gnome-terminal"},
	"$editor":      {"zed", "code", "nvim", "vim", "micro"},
	"$filemanager": {"nautilus", "thunar", "dolphin"},
	"$launcher":    {"rofi -show drun", "wofi --show drun"},
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
}

func saveDefaultVariables(path string) {
	content := "# Linutils Variable Configuration\n"
	content += "# Format: $variable = value\n\n"
	content += "$browser = chromium-browser --new-window\n"
	content += "$terminal = alacritty\n"
	content += "$editor = zed\n"
	content += "$filemanager = nautilus\n"
	content += "$launcher = wofi --show drun\n"

	os.WriteFile(path, []byte(content), 0644)
	fmt.Printf("Created default variables config at %s\n", path)
}

func SaveVariables() {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "linutils", "variables.conf")

	content := "# Linutils Variable Configuration\n"
	content += "# Format: $variable = value\n\n"
	
	// Order keys for consistency
	keys := []string{"$browser", "$terminal", "$editor", "$filemanager", "$launcher"}
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
