package modules

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rak626/linutils-rakesh/internal/pkgmanager"
)

// repos.conf format (one per line): Display Name=git-url
// Example: My Project=git@github.com:myuser/my-project.git
// Copy examples/repos.example.conf to ~/.config/linutils/repos.conf to use.
func loadRepos() map[string]string {
	repos := map[string]string{}
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".config", "linutils", "repos.conf")
	f, err := os.Open(path)
	if err != nil {
		return repos
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		url := strings.TrimSpace(parts[1])
		if name != "" && url != "" {
			repos[name] = url
		}
	}
	return repos
}

func CloneRepos(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- GitHub Repo Cloner ---")

	myRepos := loadRepos()
	if len(myRepos) == 0 {
		home, _ := os.UserHomeDir()
		fmt.Printf("No repos configured. Copy examples/repos.example.conf to %s and add your own.\n",
			filepath.Join(home, ".config", "linutils", "repos.conf"))
		return nil
	}

	var options []huh.Option[string]
	for name, url := range myRepos {
		options = append(options, huh.NewOption(name, url))
	}

	var selectedRepos []string
	var targetBaseDir string
	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, "workstation/projects")

	// 1. Select repos and target directory
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select Repositories to Clone").
				Options(options...).
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
