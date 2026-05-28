package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

// SetupEditors offers to install and configure Neovim, Vim, and IdeaVim.
func SetupEditors(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Editor Configuration ---")

	var selectedEditors []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select Editors to Setup").
				Description("Install editors and link configurations from ~/.dotfiles").
				Options(
					huh.NewOption("Neovim", "neovim"),
					huh.NewOption("Vim", "vim"),
					huh.NewOption("IdeaVim (Config Only)", "ideavim"),
				).
				Value(&selectedEditors),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if len(selectedEditors) == 0 {
		fmt.Println("No editors selected.")
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}
	dotfilesBase := filepath.Join(home, ".dotfiles")

	for _, editor := range selectedEditors {
		switch editor {
		case "neovim":
			fmt.Println("\n[Neovim Setup]")
			fmt.Println("Installing Neovim...")
			if err := manager.Install("neovim"); err != nil {
				fmt.Printf("Error installing neovim: %v\n", err)
			}

			src := filepath.Join(dotfilesBase, "nvim")
			dst := filepath.Join(home, ".config", "nvim")
			if _, err := os.Stat(src); err == nil {
				fmt.Printf("Linking %s -> %s\n", src, dst)
				if err := safeSymlink(src, dst); err != nil {
					fmt.Printf("Error symlinking neovim config: %v\n", err)
				}
			} else {
				fmt.Printf("No neovim config found at %s. Skipping configuration.\n", src)
			}

		case "vim":
			fmt.Println("\n[Vim Setup]")
			fmt.Println("Installing Vim...")
			if err := manager.Install("vim"); err != nil {
				fmt.Printf("Error installing vim: %v\n", err)
			}

			src := filepath.Join(dotfilesBase, "vim", ".vimrc")
			dst := filepath.Join(home, ".vimrc")
			if _, err := os.Stat(src); err == nil {
				fmt.Printf("Linking %s -> %s\n", src, dst)
				if err := safeSymlink(src, dst); err != nil {
					fmt.Printf("Error symlinking vim config: %v\n", err)
				}
			} else {
				fmt.Printf("No vim config found at %s. Skipping configuration.\n", src)
			}

		case "ideavim":
			fmt.Println("\n[IdeaVim Setup]")
			src := filepath.Join(dotfilesBase, "ideavim", ".ideavimrc")
			dst := filepath.Join(home, ".ideavimrc")
			if _, err := os.Stat(src); err == nil {
				fmt.Printf("Linking %s -> %s\n", src, dst)
				if err := safeSymlink(src, dst); err != nil {
					fmt.Printf("Error symlinking ideavim config: %v\n", err)
				}
			} else {
				fmt.Printf("No ideavim config found at %s. Skipping configuration.\n", src)
			}
		}
	}

	fmt.Println("\nEditor setup complete!")
	return nil
}

// safeSymlink creates a symlink at target pointing to source.
// It ensures the parent directory of target exists and removes any existing file at target.
func safeSymlink(source, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	if _, err := os.Lstat(target); err == nil {
		if err := os.RemoveAll(target); err != nil {
			return fmt.Errorf("failed to remove existing target %s: %v", target, err)
		}
	}

	return os.Symlink(source, target)
}
