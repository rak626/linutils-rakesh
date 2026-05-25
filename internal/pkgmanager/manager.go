package pkgmanager

import (
	"fmt"
	"os"
	"os/exec"
)

// PackageManager defines the interface for different distro package managers.
type PackageManager interface {
	Install(packages ...string) error
	Remove(packages ...string) error
	Update() error
	Upgrade() error
	IsInstalled(pkg string) bool
}

// GetManager returns the appropriate PackageManager for the detected OS.
func GetManager(osName string) (PackageManager, error) {
	switch osName {
	case "debian", "ubuntu", "pop", "linuxmint":
		return &AptManager{}, nil
	case "arch", "manjaro":
		return &PacmanManager{}, nil
	case "fedora":
		return &DnfManager{}, nil
	default:
		return nil, fmt.Errorf("unsupported distribution: %s", osName)
	}
}

// RunCommand is a helper to execute shell commands and pipe output to stdout/stderr.
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsCommandAvailable checks if a command exists in the system PATH.
func IsCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// checkPackageInstalled is a generic helper that uses 'which' or 'command -v'.
// Note: This is a fallback; specific managers should implement better checks.
func checkPackageInstalled(pkg string) bool {
	_, err := exec.LookPath(pkg)
	return err == nil
}
