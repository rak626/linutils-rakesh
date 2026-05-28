package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

func InstallCustomScripts(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Custom Scripts & Utilities ---")

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dotfilesDir := filepath.Join(home, ".dotfiles")
	binDir := filepath.Join(home, ".local", "bin")

	// Ensure ~/.local/bin exists
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	allScripts := []struct {
		Src  string
		Name string
	}{
		{"hyprland/.config/hypr/scripts/screenshot.sh", "screenshot"},
		{"waybar/.config/waybar/scripts/bt-devices.sh", "bt-devices"},
		{"waybar/.config/waybar/scripts/power-profile.sh", "power-profile"},
		{"rofi/.config/rofi/scripts/power-menu.sh", "power-menu"},
	}

	var availableScripts []string
	scriptMap := make(map[string]string)

	for _, s := range allScripts {
		srcPath := filepath.Join(dotfilesDir, s.Src)
		if _, err := os.Stat(srcPath); err == nil {
			availableScripts = append(availableScripts, s.Name)
			scriptMap[s.Name] = srcPath
		}
	}

	if len(availableScripts) == 0 {
		fmt.Println("No custom scripts found in ~/.dotfiles.")
		return nil
	}

	var selectedScripts []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select scripts to install to ~/.local/bin").
				Options(huh.NewOptions(availableScripts...)...).
				Value(&selectedScripts),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if len(selectedScripts) == 0 {
		fmt.Println("No scripts selected.")
		return nil
	}

	for _, name := range selectedScripts {
		src := scriptMap[name]
		dest := filepath.Join(binDir, name)

		fmt.Printf("Installing %s...\n", name)
		if err := copyFile(src, dest); err != nil {
			fmt.Printf("Error copying %s: %v\n", name, err)
			continue
		}

		if err := os.Chmod(dest, 0755); err != nil {
			fmt.Printf("Error setting executable permission for %s: %v\n", name, err)
		}
	}

	// Ensure ~/.local/bin is in PATH
	if err := ensureBinInPath(home); err != nil {
		fmt.Printf("Warning: Failed to ensure ~/.local/bin in PATH: %v\n", err)
	}

	fmt.Println("Custom scripts installation complete!")
	return nil
}

func ensureBinInPath(home string) error {
	pathLine := `export PATH="$HOME/.local/bin:$PATH"`
	
	configs := []string{
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".zshrc"),
	}

	for _, config := range configs {
		if _, err := os.Stat(config); err == nil {
			if err := appendToFileIfMissing(config, pathLine); err != nil {
				return err
			}
		}
	}
	return nil
}
