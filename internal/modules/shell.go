package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

func SetupShell(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Shell Configuration ---")

	// 1. Ensure tools are installed
	tools := []string{"zoxide", "starship", "fastfetch"}
	for _, tool := range tools {
		if !pkgmanager.IsCommandAvailable(tool) {
			fmt.Printf("%s is not installed. Installing it now...\n", tool)
			if err := manager.Install(tool); err != nil {
				return fmt.Errorf("failed to install %s: %v", tool, err)
			}
		}
	}

	// 2. Setup zoxide, starship in shell rc files
	if err := setupShellRC(); err != nil {
		return fmt.Errorf("failed to setup shell RC: %v", err)
	}

	// 3. Optional configurations from dotfiles
	home, _ := os.UserHomeDir()
	dotfilesDir := filepath.Join(home, ".dotfiles")

	// Starship config
	starshipDot := filepath.Join(dotfilesDir, "starship", "starship.toml")
	if _, err := os.Stat(starshipDot); err == nil {
		var confirm bool
		huh.NewConfirm().
			Title("Detected Starship config in ~/.dotfiles. Symlink it?").
			Value(&confirm).
			Run()
		if confirm {
			target := filepath.Join(home, ".config", "starship.toml")
			if err := symlink(starshipDot, target); err != nil {
				fmt.Printf("Warning: failed to symlink starship config: %v\n", err)
			}
		}
	}

	// Fastfetch config
	fastfetchDot := filepath.Join(dotfilesDir, "fastfetch")
	if _, err := os.Stat(fastfetchDot); err == nil {
		var confirm bool
		huh.NewConfirm().
			Title("Detected Fastfetch config in ~/.dotfiles. Symlink it?").
			Value(&confirm).
			Run()
		if confirm {
			target := filepath.Join(home, ".config", "fastfetch")
			if err := symlink(fastfetchDot, target); err != nil {
				fmt.Printf("Warning: failed to symlink fastfetch config: %v\n", err)
			}
		}
	}

	fmt.Println("Shell configuration complete!")
	return nil
}

func setupShellRC() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	shellPath := os.Getenv("SHELL")
	var rcFiles []string
	var shellName string

	if strings.Contains(shellPath, "zsh") {
		rcFiles = append(rcFiles, filepath.Join(home, ".zshrc"))
		shellName = "zsh"
	} else if strings.Contains(shellPath, "bash") {
		rcFiles = append(rcFiles, filepath.Join(home, ".bashrc"))
		shellName = "bash"
	} else {
		fmt.Println("Could not detect shell. Defaulting to .bashrc and .zshrc if they exist.")
		rcFiles = []string{filepath.Join(home, ".bashrc"), filepath.Join(home, ".zshrc")}
		shellName = "bash" // fallback for init command
	}

	zoxideCmd := fmt.Sprintf("eval \"$(zoxide init %s)\"", shellName)
	starshipCmd := fmt.Sprintf("eval \"$(starship init %s)\"", shellName)

	for _, rcFile := range rcFiles {
		if _, err := os.Stat(rcFile); err == nil {
			appendToFileIfMissing(rcFile, zoxideCmd)
			appendToFileIfMissing(rcFile, starshipCmd)
		}
	}
	
	return nil
}

func symlink(src, dest string) error {
	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	// Remove existing if it's a file or symlink
	if _, err := os.Lstat(dest); err == nil {
		if err := os.RemoveAll(dest); err != nil {
			return err
		}
	}

	fmt.Printf("Symlinking %s -> %s\n", src, dest)
	return os.Symlink(src, dest)
}
