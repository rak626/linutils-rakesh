package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

const fontsRepo = "git@github.com:rak626/fonts.git"

func SetupFonts(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Fonts Setup ---")

	home, _ := os.UserHomeDir()
	fontsDir := filepath.Join(home, ".local/share/fonts")
	tempCloneDir := filepath.Join(os.TempDir(), "rak626-fonts")

	// 1. Ensure fonts directory exists
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		return fmt.Errorf("failed to create fonts directory: %v", err)
	}

	// 2. Clone fonts repo
	fmt.Printf("Cloning fonts from %s...\n", fontsRepo)
	if _, err := os.Stat(tempCloneDir); err == nil {
		os.RemoveAll(tempCloneDir)
	}

	cmd := exec.Command("git", "clone", "--depth", "1", fontsRepo, tempCloneDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone fonts repo: %v", err)
	}

	// 3. Copy fonts to ~/.local/share/fonts
	fmt.Println("Installing fonts...")
	err := filepath.Walk(tempCloneDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext == ".ttf" || ext == ".otf" || ext == ".woff" || ext == ".woff2" {
			dest := filepath.Join(fontsDir, filepath.Base(path))
			fmt.Printf("Copying %s to %s\n", filepath.Base(path), fontsDir)
			
			input, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(dest, input, 0644)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to install fonts: %v", err)
	}

	// 4. Update font cache
	fmt.Println("Updating font cache...")
	pkgmanager.RunCommand("fc-cache", "-fv")

	// 5. Cleanup
	os.RemoveAll(tempCloneDir)

	fmt.Println("Fonts setup complete!")
	return nil
}
