package pkgmanager

import (
	"fmt"
	"os/exec"
)

type PacmanManager struct{}

func (m *PacmanManager) getAURHelper() string {
	if IsCommandAvailable("yay") {
		return "yay"
	}
	if IsCommandAvailable("paru") {
		return "paru"
	}
	return ""
}

func (m *PacmanManager) Install(packages ...string) error {
	var toInstall []string
	for _, pkg := range packages {
		if !m.IsInstalled(pkg) {
			toInstall = append(toInstall, pkg)
		} else {
			fmt.Printf("Package %s is already installed. Skipping.\n", pkg)
		}
	}

	if len(toInstall) == 0 {
		return nil
	}

	helper := m.getAURHelper()

	// Split packages into official and AUR
	var official []string
	var aur []string

	for _, pkg := range toInstall {
		// Check if package is in official repos
		err := exec.Command("pacman", "-Si", pkg).Run()
		if err == nil {
			official = append(official, pkg)
		} else if helper != "" {
			aur = append(aur, pkg)
		} else {
			fmt.Printf("Warning: Package %s not found in official repositories and no AUR helper found.\n", pkg)
		}
	}

	if len(official) > 0 {
		args := append([]string{"-S", "--noconfirm"}, official...)
		fmt.Printf("Installing official packages: sudo pacman %v\n", args)
		if err := RunCommand("sudo", append([]string{"pacman"}, args...)...); err != nil {
			return err
		}
	}

	if len(aur) > 0 {
		args := append([]string{"-S", "--noconfirm"}, aur...)
		fmt.Printf("Installing AUR packages using %s: %s %v\n", helper, helper, args)
		// AUR helpers should not be run as root/sudo
		if err := RunCommand(helper, args...); err != nil {
			return err
		}
	}

	return nil
}

func (m *PacmanManager) Remove(packages ...string) error {
	if len(packages) == 0 {
		return nil
	}
	
	helper := m.getAURHelper()
	if helper != "" {
		args := append([]string{"-Rs", "--noconfirm"}, packages...)
		fmt.Printf("Removing packages using %s: %s %v\n", helper, helper, args)
		return RunCommand(helper, args...)
	}

	args := append([]string{"-Rs", "--noconfirm"}, packages...)
	fmt.Printf("Running: sudo pacman %v\n", args)
	return RunCommand("sudo", append([]string{"pacman"}, args...)...)
}

func (m *PacmanManager) Update() error {
	helper := m.getAURHelper()
	if helper != "" {
		fmt.Printf("Running: %s -Sy\n", helper)
		return RunCommand(helper, "-Sy")
	}
	fmt.Println("Running: sudo pacman -Sy")
	return RunCommand("sudo", "pacman", "-Sy")
}

func (m *PacmanManager) Upgrade() error {
	helper := m.getAURHelper()
	if helper != "" {
		fmt.Printf("Running: %s -Syu --noconfirm\n", helper)
		return RunCommand(helper, "-Syu", "--noconfirm")
	}
	fmt.Println("Running: sudo pacman -Syu --noconfirm")
	return RunCommand("sudo", "pacman", "-Syu", "--noconfirm")
}

func (m *PacmanManager) IsInstalled(pkg string) bool {
	// pacman -Qq pkg returns 0 if installed exactly, 1 if not.
	// -Qs is a search, might return partial matches.
	err := exec.Command("pacman", "-Qq", pkg).Run()
	return err == nil
}
