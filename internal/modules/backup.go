package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rak626/linutils-rakesh/internal/pkgmanager"
	"github.com/rak626/linutils-rakesh/internal/system"
)

// BackupToGit dumps dconf + package list into the dotfiles repo (git-only)
// and commits. No local timestamped backups — history is git log.
func BackupToGit(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	_ = manager
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dotfilesDir := filepath.Join(home, ".dotfiles")
	if _, err := os.Stat(filepath.Join(dotfilesDir, ".git")); err != nil {
		return fmt.Errorf("dotfiles repo not found at %s — run Fresh Install or Restore first", dotfilesDir)
	}

	fmt.Println("\n--- Backup Now to Git (git-only) ---")

	// 1. dconf dump (GNOME only)
	if pkgmanager.IsCommandAvailable("dconf") {
		out, err := exec.Command("dconf", "dump", "/").Output()
		if err == nil {
			gnomeDir := filepath.Join(dotfilesDir, "gnome")
			os.MkdirAll(gnomeDir, 0755)
			if err := os.WriteFile(filepath.Join(gnomeDir, "dconf.ini"), out, 0644); err != nil {
				fmt.Printf("Warning: failed to write dconf dump: %v\n", err)
			} else {
				fmt.Println("[OK] GNOME dconf dumped to gnome/dconf.ini")
			}
		}
	}

	// 2. Package list per distro
	pkgFile := filepath.Join(dotfilesDir, fmt.Sprintf("packages.%s.txt", sysInfo.DistroID))
	var pkgCmd *exec.Cmd
	switch sysInfo.DistroID {
	case "arch", "manjaro", "endeavouros":
		pkgCmd = exec.Command("bash", "-c", "pacman -Qe | awk '{print $1}'")
	case "fedora", "nobara":
		pkgCmd = exec.Command("bash", "-c", "dnf list installed 2>/dev/null | awk 'NR>1{print $1}'")
	default:
		pkgCmd = exec.Command("bash", "-c", "dpkg -l 2>/dev/null | awk '/^ii/{print $2}'")
	}
	if out, err := pkgCmd.Output(); err == nil {
		os.WriteFile(pkgFile, out, 0644)
		fmt.Printf("[OK] Package list written to %s\n", filepath.Base(pkgFile))
	}

	// 3. Commit + push
	ts := time.Now().Format("2006-01-02 15:04")
	add := exec.Command("git", "-C", dotfilesDir, "add", "-A")
	add.Stdout, add.Stderr = os.Stdout, os.Stderr
	_ = add.Run()
	commit := exec.Command("git", "-C", dotfilesDir, "commit", "-m", fmt.Sprintf("linutils backup %s (%s/%s)", ts, sysInfo.DistroID, sysInfo.DEID))
	commit.Stdout, commit.Stderr = os.Stdout, os.Stderr
	if err := commit.Run(); err != nil {
		fmt.Println("Nothing new to commit or commit failed — continuing.")
		return nil
	}
	push := exec.Command("git", "-C", dotfilesDir, "push")
	push.Stdin, push.Stdout, push.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := push.Run(); err != nil {
		return fmt.Errorf("git push failed (check gh auth / remotes): %v", err)
	}
	fmt.Println("[OK] Backup pushed to dotfiles repo.")
	return nil
}

// IsGitHubAuthenticated reports whether `gh auth status` succeeds.
func IsGitHubAuthenticated() bool {
	if _, err := exec.LookPath("gh"); err != nil {
		return false
	}
	return exec.Command("gh", "auth", "status").Run() == nil
}
