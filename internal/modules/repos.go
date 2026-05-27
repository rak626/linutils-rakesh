package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

var myRepos = map[string]string{
	"DSA Tracker":      "git@github.com:rak626/dsa-tracker.git",
	"Java Learning":    "git@github.com:rak626/java-learning.git",
	"Obsidian Vault":   "git@github.com:rak626/obsidian-vault.git",
	"Rakesh Portfolio": "git@github.com:rak626/rakesh-portfolio.git",
}

func CloneRepos(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- GitHub Repo Cloner ---")

	var selectedRepos []string
	var targetBaseDir string
	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, "workstation/projects")

	// 1. Select repos and target directory
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select Repositories to Clone").
				Options(
					huh.NewOption("DSA Tracker", myRepos["DSA Tracker"]),
					huh.NewOption("Java Learning", myRepos["Java Learning"]),
					huh.NewOption("Obsidian Vault", myRepos["Obsidian Vault"]),
					huh.NewOption("Rakesh Portfolio", myRepos["Rakesh Portfolio"]),
				).
				Value(&selectedRepos),

			huh.NewInput().
				Title("Target Directory").
				Placeholder(defaultDir).
				Value(&targetBaseDir),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if len(selectedRepos) == 0 {
		fmt.Println("No repositories selected.")
		return nil
	}

	if strings.TrimSpace(targetBaseDir) == "" {
		targetBaseDir = defaultDir
	}

	// Expand ~/ if present
	if strings.HasPrefix(targetBaseDir, "~/") {
		targetBaseDir = filepath.Join(home, targetBaseDir[2:])
	}

	// 2. Clone repos
	if err := os.MkdirAll(targetBaseDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %v", err)
	}

	for _, repoURL := range selectedRepos {
		// Extract repo name from URL
		parts := strings.Split(repoURL, "/")
		repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
		targetPath := filepath.Join(targetBaseDir, repoName)

		if _, err := os.Stat(targetPath); err == nil {
			fmt.Printf("Repository %s already exists at %s, skipping...\n", repoName, targetPath)
			continue
		}

		fmt.Printf("Cloning %s into %s...\n", repoName, targetPath)
		cmd := exec.Command("git", "clone", repoURL, targetPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: failed to clone %s: %v\n", repoName, err)
		}
	}

	fmt.Println("Repository cloning complete!")
	return nil
}
