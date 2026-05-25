package modules

import (
	"fmt"
	"os/exec"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

var BloatServices = []string{
	"tracker-miner-fs-3.service",
	"tracker-extract-3.service",
	"packagekit.service",
	"flatpak-system-helper.service",
	"abrtd.service",
	"ModemManager.service",
	"avahi-daemon.service",
	"cups.service",
}

func GetBloatServices(sysInfo system.Info) []string {
	services := make([]string, len(BloatServices))
	copy(services, BloatServices)
	if sysInfo.OS != "fedora" && sysInfo.OS != "arch" && sysInfo.OS != "manjaro" {
		services = append(services, "gnome-software.service")
	}
	return services
}

func MaskServices(services []string) error {
	for _, s := range services {
		if err := setServiceMask(s, true); err != nil {
			fmt.Printf("Error masking %s: %v\n", s, err)
		}
	}
	return nil
}

func UnmaskServices(services []string) error {
	for _, s := range services {
		if err := setServiceMask(s, false); err != nil {
			fmt.Printf("Error unmasking %s: %v\n", s, err)
		}
	}
	return nil
}

func setServiceMask(service string, mask bool) error {
	action := "mask"
	if !mask {
		action = "unmask"
	}

	userServices := map[string]bool{
		"tracker-miner-fs-3.service": true,
		"tracker-extract-3.service": true,
		"gnome-software.service":     true,
	}

	var cmd *exec.Cmd
	if userServices[service] {
		cmd = exec.Command("sudo", "systemctl", "--global", action, service)
	} else {
		cmd = exec.Command("sudo", "systemctl", action, service)
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Printf("%sed %s\n", action, service)
	return nil
}
