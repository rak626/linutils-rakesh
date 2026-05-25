package modules

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
	"github.com/rakesh/linutils-rakesh/internal/tui"
)

func DebloatGnome(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	for {
		var choice string
		
		options := []huh.Option[string]{
			huh.NewOption("Debloat Apps", "apps"),
			huh.NewOption("Debloat Services", "services"),
		}

		if sysInfo.OS == "fedora" && sysInfo.DE == "gnome" {
			options = append(options, huh.NewOption("Debloat Extensions", "extensions"))
		}

		options = append(options, huh.NewOption("Back", "back"))

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Gnome Debloater").
					Options(options...).
					Value(&choice),
			),
		)

		err := form.Run()
		if err != nil {
			return err
		}

		if choice == "back" {
			return nil
		}

		switch choice {
		case "apps":
			appsList := GetBloatApps(sysInfo)
			items := make([]tui.ListItem, len(appsList))
			for i, name := range appsList {
				items[i] = tui.ListItem{Key: name, Name: name}
			}
			action, selectedItems, err := tui.RunListUI("Debloat Apps", items)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			selected := getSelectedKeys(selectedItems)
			if action == "r" && len(selected) > 0 {
				RemoveApps(manager, selected)
			} else if action == "i" && len(selected) > 0 {
				InstallApps(manager, selected)
			}

		case "services":
			servicesList := GetBloatServices(sysInfo)
			items := make([]tui.ListItem, len(servicesList))
			for i, name := range servicesList {
				items[i] = tui.ListItem{Key: name, Name: name}
			}
			action, selectedItems, err := tui.RunListUI("Debloat Services", items)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			selected := getSelectedKeys(selectedItems)
			if action == "r" && len(selected) > 0 {
				MaskServices(selected)
			} else if action == "i" && len(selected) > 0 {
				UnmaskServices(selected)
			}

		case "extensions":
			extensionsList := GetBloatExtensions()
			items := make([]tui.ListItem, len(extensionsList))
			for i, name := range extensionsList {
				items[i] = tui.ListItem{Key: name, Name: name}
			}
			action, selectedItems, err := tui.RunListUI("Debloat Extensions", items)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			selected := getSelectedKeys(selectedItems)
			if action == "r" && len(selected) > 0 {
				RemoveExtensions(manager, selected)
			} else if action == "i" && len(selected) > 0 {
				InstallExtensions(manager, selected)
			}
		}
		
		fmt.Println("\nPress Enter to continue...")
		fmt.Scanln()
	}
}

func getSelectedKeys(items []tui.ListItem) []string {
	var selected []string
	for _, item := range items {
		if item.Selected {
			selected = append(selected, item.Key)
		}
	}
	return selected
}
