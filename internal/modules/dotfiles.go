package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/huh"
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

	// 3. List directories for stowing
	entries, err := os.ReadDir(dotfilesDir)
	if err != nil {
		return fmt.Errorf("failed to read dotfiles directory: %v", err)
	}

	var folders []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name()[0] != '.' && entry.Name() != "kitty" {
			folders = append(folders, entry.Name())
		}
	}

	if len(folders) == 0 {
		fmt.Println("No stowable folders found in ~/.dotfiles")
		return nil
	}

	// 4. Select folders to stow
	var selectedFolders []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select folders to stow").
				Description("These will be symlinked to your home directory.").
				Options(huh.NewOptions(folders...)...).
				Value(&selectedFolders),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if len(selectedFolders) == 0 {
		fmt.Println("No folders selected for stowing.")
		return nil
	}

	// 5. Run Stow
	fmt.Println("Stowing selected configurations...")
	for _, folder := range selectedFolders {
		fmt.Printf("Stowing %s...\n", folder)
		
		// Pre-stow cleanup to handle auto-created directories (like Hyprland)
		prepareForStow(home, dotfilesDir, folder)

		// stow -v -R -t ~ folder
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

	// Read directories inside ~/.dotfiles/<folder>/.config/
	entries, err := os.ReadDir(sourceConfigPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			targetDir := filepath.Join(home, ".config", entry.Name())
			// Check if target directory already exists and is not a symlink
			if info, err := os.Lstat(targetDir); err == nil {
				if info.Mode()&os.ModeSymlink == 0 && info.IsDir() {
					fmt.Printf("  Cleaning up existing directory to allow stowing: %s\n", targetDir)
					// Rename to .bak to be safe, removing old bak if needed
					bakDir := targetDir + ".bak"
					os.RemoveAll(bakDir)
					if err := os.Rename(targetDir, bakDir); err != nil {
						// If rename fails (e.g., cross-device or permission), try to remove
						os.RemoveAll(targetDir)
					}
					// Recreate empty directory so stow links files inside it
					os.MkdirAll(targetDir, 0755)
				}
			}
		}
	}
}
