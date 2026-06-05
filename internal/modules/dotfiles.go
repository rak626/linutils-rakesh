package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

const dotfilesRepo = "https://github.com/rak626/dotfiles.git"

func SetupDotfiles(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Dotfiles Sync (GNU Stow) ---")

	// 1. Ensure stow is installed
	if !pkgmanager.IsCommandAvailable("stow") {
		fmt.Println("Installing GNU Stow...")
		if err := manager.Install("stow"); err != nil {
			return fmt.Errorf("failed to install stow: %v", err)
		}
	}

	home, _ := os.UserHomeDir()
	dotfilesDir := filepath.Join(home, ".dotfiles")

	// 2. Clone or Pull Dotfiles
	if _, err := os.Stat(dotfilesDir); os.IsNotExist(err) {
		fmt.Printf("Cloning dotfiles to %s...\n", dotfilesDir)
		cmd := exec.Command("git", "clone", dotfilesRepo, dotfilesDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone dotfiles: %v", err)
		}
	} else {
		fmt.Println("Dotfiles directory exists, pulling latest changes...")
		cmd := exec.Command("git", "-C", dotfilesDir, "pull")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	// 3. Filter directories for stowing
	entries, err := os.ReadDir(dotfilesDir)
	if err != nil {
		return fmt.Errorf("failed to read dotfiles directory: %v", err)
	}

	// Programs to stow (excluding WMs like i3, miracle-wm)
	stowWhitelist := map[string]bool{
		"alacritty": true, "bashrc": true, "btop": true, "fastfetch": true,
		"gtk": true, "ideavim": true, "mako": true, "nvim": true,
		"picom": true, "qt": true, "rofi": true, "scripts": true,
		"starship": true, "swayosd": true, "ulauncher": true, "uwsm": true,
		"vim": true, "waybar": true, "wofi": true, "desktop": true,
	}

	hyprVersion := getHyprlandVersion()
	isLuaVersion := isVersionGreaterOrEqual(hyprVersion, "0.55")

	var folders []string
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() || name[0] == '.' || name == "kitty" {
			continue
		}

		// Handle Hyprland versioning
		if name == "hyprland-lua" {
			if isLuaVersion {
				folders = append(folders, name)
			}
			continue
		}
		if name == "hyprland" {
			if !isLuaVersion {
				folders = append(folders, name)
			}
			continue
		}

		// Check against whitelist
		if stowWhitelist[name] {
			folders = append(folders, name)
		}
	}

	if len(folders) == 0 {
		fmt.Println("No stowable folders found in ~/.dotfiles")
		return nil
	}

	// 4. Automated Stow (stow all found folders)
	fmt.Printf("Automating stowing of folders: %v\n", folders)
	for _, folder := range folders {
		fmt.Printf("Stowing %s...\n", folder)
		
		// Pre-stow cleanup
		prepareForStow(home, dotfilesDir, folder)

		cmd := exec.Command("stow", "-v", "-R", "-t", home, folder)
		cmd.Dir = dotfilesDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: failed to stow %s: %v\n", folder, err)
		}
	}

	fmt.Println("Dotfiles sync complete!")
	return nil
}

func getHyprlandVersion() string {
	cmd := exec.Command("hyprland", "--version")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	// Typical output: "Hyprland, version 0.41.2 (34608c0...)"
	return string(out)
}

// isVersionGreaterOrEqual compares a version string (like "0.41.2") against a target (like "0.55")
func isVersionGreaterOrEqual(versionStr, target string) bool {
	if versionStr == "" {
		return true // Default to latest if not found/detectable during fresh install
	}

	// Simple extraction: look for digits and dots
	var v string
	for _, part := range strings.Fields(versionStr) {
		if strings.Contains(part, ".") {
			v = part
			break
		}
	}

	if v == "" {
		return true
	}

	// Just a simple string comparison for "0.55" vs "0.xx" might work if we assume 0.x format
	// But let's do a slightly better check
	parts := strings.Split(v, ".")
	targetParts := strings.Split(target, ".")

	for i := 0; i < len(parts) && i < len(targetParts); i++ {
		var pv, tv int
		fmt.Sscanf(parts[i], "%d", &pv)
		fmt.Sscanf(targetParts[i], "%d", &tv)
		if pv > tv {
			return true
		}
		if pv < tv {
			return false
		}
	}
	return len(parts) >= len(targetParts)
}


// prepareForStow checks the structure inside ~/.dotfiles/<folder>
// If it maps to a directory in ~/.config (e.g., ~/.config/hypr) that already exists,
// it deletes the contents of the target directory so stow can link individual files
// without complaining about existing directories or trying to link the parent.
func prepareForStow(home, dotfilesDir, folder string) {
	// Most dotfiles are stowed to ~/.config. We check if .config exists in the stow package.
	sourceConfigPath := filepath.Join(dotfilesDir, folder, ".config")
	if _, err := os.Stat(sourceConfigPath); os.IsNotExist(err) {
		return // Not a standard .config stow package, skip cleanup
	}

	// Read directories/files inside ~/.dotfiles/<folder>/.config/
	entries, err := os.ReadDir(sourceConfigPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		targetPath := filepath.Join(home, ".config", entry.Name())
		// Check if target exists
		if _, err := os.Lstat(targetPath); err == nil {
			fmt.Printf("  Removing existing entry to allow stowing: %s\n", targetPath)
			os.RemoveAll(targetPath)
		}
		
		// If it's a directory in the source, we might want to recreate the parent 
		// if we were stowing deep, but stow -t ~ usually handles the link creation.
		// However, the previous logic recreated the directory. 
		// If the user wants a clean stow, removing and letting stow handle it is usually better.
	}
}
