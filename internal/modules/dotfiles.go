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
