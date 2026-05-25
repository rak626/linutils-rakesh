package modules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

func SetupShell(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Shell Configuration ---")

	// 1. Ensure zoxide is installed
	if !pkgmanager.IsCommandAvailable("zoxide") {
		fmt.Println("zoxide is not installed. Installing it now...")
		if err := manager.Install("zoxide"); err != nil {
			return fmt.Errorf("failed to install zoxide: %v", err)
		}
	}

	// 2. Setup zoxide in shell rc files
	if err := setupZoxide(); err != nil {
		return fmt.Errorf("failed to setup zoxide: %v", err)
	}

	fmt.Println("Shell configuration complete!")
	return nil
}

func setupZoxide() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	shellPath := os.Getenv("SHELL")
	var rcFile string
	var shellName string

	if strings.Contains(shellPath, "zsh") {
		rcFile = filepath.Join(home, ".zshrc")
		shellName = "zsh"
	} else if strings.Contains(shellPath, "bash") {
		rcFile = filepath.Join(home, ".bashrc")
		shellName = "bash"
	} else {
		// Fallback to both if we can't detect, or just skip
		fmt.Println("Could not detect shell. Skipping automatic rc file update.")
		return nil
	}

	zoxideCmd := fmt.Sprintf("eval \"$(zoxide init %s)\"", shellName)
	
	return appendToFileIfMissing(rcFile, zoxideCmd)
}

func appendToFileIfMissing(filePath, line string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create file if it doesn't exist
		return os.WriteFile(filePath, []byte(line+"\n"), 0644)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), line) {
			fmt.Printf("Line already exists in %s: %s\n", filePath, line)
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Append line
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + line + "\n"); err != nil {
		return err
	}

	fmt.Printf("Added to %s: %s\n", filePath, line)
	return nil
}
