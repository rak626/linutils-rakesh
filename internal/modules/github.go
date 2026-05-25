package modules

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

func SetupGitHub(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- GitHub Setup ---")

	// 1. Install GitHub CLI (gh)
	ghPkg := "github-cli"
	if _, err := exec.LookPath("dnf"); err == nil {
		ghPkg = "gh"
	} else if _, err := exec.LookPath("pacman"); err == nil {
		ghPkg = "github-cli"
	} else if _, err := exec.LookPath("apt"); err == nil {
		ghPkg = "gh"
	}

	if !manager.IsInstalled(ghPkg) && !pkgmanager.IsCommandAvailable("gh") {
		fmt.Println("Installing GitHub CLI...")
		if err := manager.Install(ghPkg); err != nil {
			return err
		}
	} else {
		fmt.Println("GitHub CLI is already installed.")
	}

	// 2. Check SSH Key
	sshKeyPath := os.ExpandEnv("$HOME/.ssh/id_ed25519")
	if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
		var generateKey bool
		huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("No SSH key found. Generate one?").
					Description("Recommended for secure GitHub access (Ed25519)").
					Value(&generateKey),
			),
		).Run()

		if generateKey {
			fmt.Println("Generating SSH key...")
			cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", "linutils-rakesh", "-N", "", "-f", sshKeyPath)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to generate SSH key: %v", err)
			}
			fmt.Println("SSH key generated at", sshKeyPath)
		}
	}

	// 3. Check Auth Status
	err := exec.Command("gh", "auth", "status").Run()
	if err == nil {
		fmt.Println("Already authenticated with GitHub.")
		return nil
	}

	// 4. Ask to authenticate
	var authenticate bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Authenticate with GitHub?").
				Description("This will run 'gh auth login' using SSH protocol").
				Value(&authenticate),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if authenticate {
		// We use TTY for interactive login
		// --git-protocol ssh ensures git is configured to use ssh
		// --web allows browser-based auth which is often easier
		cmd := exec.Command("gh", "auth", "login", "--git-protocol", "ssh", "--web")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return nil
}
