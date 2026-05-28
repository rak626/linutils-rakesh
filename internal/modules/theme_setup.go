package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
)

// IntegrateThemeSwitcher sets up the necessary hooks in various config files
// to allow the theme switcher to work correctly.
func IntegrateThemeSwitcher() error {
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Integrate Theme Switcher?").
				Description("This will add hooks to your config files (Nvim, Hyprland, i3, Waybar, Starship) to enable dynamic theme switching.").
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if !confirm {
		fmt.Println("Theme switcher integration skipped.")
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Neovim
	nvimConfig := filepath.Join(home, ".config/nvim/init.lua")
	if err := appendToFileIfMissing(nvimConfig, "require('active_theme')"); err != nil {
		fmt.Printf("Warning: Failed to update Neovim config: %v\n", err)
	}

	// Hyprland
	hyprConfig := filepath.Join(home, ".config/hypr/hyprland.conf")
	if err := appendToFileIfMissing(hyprConfig, "source = ~/.config/hypr/active_theme.conf"); err != nil {
		fmt.Printf("Warning: Failed to update Hyprland config: %v\n", err)
	}

	// i3
	i3Config := filepath.Join(home, ".config/i3/config")
	if err := appendToFileIfMissing(i3Config, "include ~/.config/i3/active_theme.i3"); err != nil {
		fmt.Printf("Warning: Failed to update i3 config: %v\n", err)
	}

	// Waybar
	waybarConfig := filepath.Join(home, ".config/waybar/style.css")
	if err := prependToFileIfMissing(waybarConfig, "@import \"active_theme.css\";"); err != nil {
		fmt.Printf("Warning: Failed to update Waybar config: %v\n", err)
	}

	// Starship
	starshipConfig := filepath.Join(home, ".config/starship.toml")
	if err := integrateStarship(starshipConfig); err != nil {
		fmt.Printf("Warning: Failed to update Starship config: %v\n", err)
	}

	fmt.Println("Theme switcher integration complete!")
	return nil
}

func integrateStarship(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create file if it doesn't exist
		return os.WriteFile(filePath, []byte("palette = \"default\"\n"), 0644)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), "palette =") {
		fmt.Printf("Starship already has a palette defined in %s\n", filePath)
		return nil
	}

	// Append palette definition
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("\npalette = \"default\"\n"); err != nil {
		return err
	}

	fmt.Printf("Added palette to %s\n", filePath)
	return nil
}
