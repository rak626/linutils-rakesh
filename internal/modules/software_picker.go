package modules

import (
	"fmt"
	"sort"

	"github.com/rak626/linutils-rakesh/internal/config"
	"github.com/rak626/linutils-rakesh/internal/pkgmanager"
	"github.com/rak626/linutils-rakesh/internal/system"
	"github.com/rak626/linutils-rakesh/internal/tui"
)

// SoftwareCatalogItem joins all install maps into one pickable list.
func BuildSoftwareCatalog() []tui.ListItem {
	var items []tui.ListItem
	add := func(cat string, m map[string]config.InstallConfig) {
		for key, inst := range m {
			items = append(items, tui.ListItem{
				Key:         cat + ":" + key,
				Name:        inst.Name,
				Category:    cat,
				Description: "Install " + inst.Name,
			})
		}
	}
	add("Browser", map[string]config.InstallConfig{
		"chromium": {Name: "Chromium Browser"},
		"zen":      config.ManualInstalls["zen"],
		"brave":    config.ManualInstalls["brave"],
	})
	add("Editor", map[string]config.InstallConfig{
		"zed": config.ManualInstalls["zed"],
	})
	add("Dev", map[string]config.InstallConfig{
		"sdkman":            config.ManualInstalls["sdkman"],
		"nvm":               config.ManualInstalls["nvm"],
		"starship":          config.ManualInstalls["starship"],
		"bun":               config.ManualInstalls["bun"],
		"jetbrains-toolbox": config.ManualInstalls["jetbrains-toolbox"],
		"dbeaver":           config.ManualInstalls["dbeaver"],
	})
	add("AI", config.AIInstalls)
	add("Helper", config.HelperInstalls)
	add("Flatpak", config.FlatpakInstalls)
	sort.Slice(items, func(a, b int) bool {
		if items[a].Category == items[b].Category {
			return items[a].Name < items[b].Name
		}
		return items[a].Category < items[b].Category
	})
	return items
}

// InstallSoftwarePicker asks every time (no fixed defaults) and installs only selected.
func InstallSoftwarePicker(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	items := BuildSoftwareCatalog()
	action, results, err := tui.RunListUIWithDesc("Software Picker", "Space selects, Enter installs. Nothing auto-installed.", items)
	if err != nil || action == "" || action == "back" || action == "quit" {
		return err
	}
	lookup := map[string]config.InstallConfig{}
	for k, v := range config.ManualInstalls {
		lookup[k] = v
	}
	for k, v := range config.AIInstalls {
		lookup[k] = v
	}
	for k, v := range config.HelperInstalls {
		lookup[k] = v
	}
	for k, v := range config.FlatpakInstalls {
		lookup[k] = v
	}
	// Chromium is a plain package, not in maps.
	selected := 0
	for _, it := range results {
		if !it.Selected {
			continue
		}
		selected++
		// Key format "Cat:key"
		key := it.Key
		for i := len(key) - 1; i >= 0; i-- {
			if key[i] == ':' {
				key = key[i+1:]
				break
			}
		}
		if key == "chromium" {
			fmt.Println("Installing Chromium Browser...")
			_ = manager.Install("chromium")
			continue
		}
		if inst, ok := lookup[key]; ok {
			installFromConfig(manager, sysInfo, inst)
		}
	}
	if selected == 0 {
		fmt.Println("No software selected.")
	}
	return nil
}
