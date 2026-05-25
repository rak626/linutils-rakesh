package modules

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

type GitConfig struct {
	Username   string
	Email      string
	InitBranch string
}

func SetupGit(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Git Setup ---")

	// 1. Ensure Git is installed
	if !manager.IsInstalled("git") {
		fmt.Println("Installing git...")
		if err := manager.Install("git"); err != nil {
			return err
		}
	} else {
		fmt.Println("Git is already installed.")
	}

	// 2. Check current config and environment variables for defaults
	currentName := getGitConfig("user.name")
	currentEmail := getGitConfig("user.email")
	currentBranch := getGitConfig("init.defaultBranch")

	// Prioritize environment variables if set, otherwise fallback to current config
	defaultName := os.Getenv("GIT_USER_NAME")
	if defaultName == "" {
		defaultName = currentName
	}
	defaultEmail := os.Getenv("GIT_USER_EMAIL")
	if defaultEmail == "" {
		defaultEmail = currentEmail
	}

	var cfg GitConfig
	cfg.Username = defaultName
	cfg.Email = defaultEmail
	cfg.InitBranch = currentBranch
	if cfg.InitBranch == "" {
		cfg.InitBranch = "main"
	}

	// 3. Ask for updates if needed with validation
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Git Username").
				Placeholder("Enter your full name").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("username cannot be empty")
					}
					return nil
				}).
				Value(&cfg.Username),
			huh.NewInput().
				Title("Git Email").
				Placeholder("Enter your email address").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("email cannot be empty")
					}
					if !strings.Contains(s, "@") {
						return fmt.Errorf("invalid email address")
					}
					return nil
				}).
				Value(&cfg.Email),
			huh.NewInput().
				Title("Default Branch Name").
				Value(&cfg.InitBranch),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	// 4. Apply changes (Idempotent check and whitespace check)
	finalName := strings.TrimSpace(cfg.Username)
	finalEmail := strings.TrimSpace(cfg.Email)

	if finalName != "" && finalName != currentName {
		setGitConfig("user.name", finalName)
	}
	if finalEmail != "" && finalEmail != currentEmail {
		setGitConfig("user.email", finalEmail)
	}
	if cfg.InitBranch != "" && cfg.InitBranch != currentBranch {
		setGitConfig("init.defaultBranch", cfg.InitBranch)
	}

	fmt.Println("Git configuration updated.")
	return nil
}

func getGitConfig(key string) string {
	out, err := exec.Command("git", "config", "--global", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func setGitConfig(key, value string) {
	fmt.Printf("Setting git %s to %s\n", key, value)
	exec.Command("git", "config", "--global", key, value).Run()
}
