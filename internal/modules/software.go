package modules

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/config"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/tui"
)

type SoftwareConfig struct {
	General []string
	WebApps []string
	Manual  []string
	AI      []string
	Helpers []string
}

func InstallSoftware(manager pkgmanager.PackageManager) error {
	var items []tui.ListItem

	// 1. General Software
	items = append(items, tui.ListItem{Key: "chromium", Name: "Chromium Browser", Category: "General Software"})

	// 2. Manual Installs
	for key, inst := range config.ManualInstalls {
		items = append(items, tui.ListItem{Key: key, Name: inst.Name, Category: "Manual Installs (curl/fsSL)"})
	}

	// 3. Web Apps
	items = append(items, tui.ListItem{Key: "custom", Name: "Add Custom Web App", Category: "Web Apps (Chromium based)"})
	for name, url := range config.WebApps {
		items = append(items, tui.ListItem{Key: url, Name: name, Category: "Web Apps (Chromium based)"})
	}

	// 4. AI Tools
	for key, inst := range config.AIInstalls {
		items = append(items, tui.ListItem{Key: key, Name: inst.Name, Category: "AI Tools"})
	}

	// 5. Helper Tools
	for key, inst := range config.HelperInstalls {
		items = append(items, tui.ListItem{Key: key, Name: inst.Name, Category: "Helper Tools"})
	}

	action, results, err := tui.RunListUI("Software Installer", items)
	if err != nil {
		return err
	}

	if action == "" {
		return nil
	}

	selectedCount := 0
	for _, item := range results {
		if item.Selected {
			selectedCount++
		}
	}

	if selectedCount == 0 {
		fmt.Println("No items selected.")
		return nil
	}

	if action == "i" {
		fmt.Println("\n--- Installing Selected Software ---")
		for _, item := range results {
			if !item.Selected {
				continue
			}

			switch item.Category {
			case "General Software":
				manager.Install(item.Key)
			case "Manual Installs (curl/fsSL)":
				installFromConfig(manager, config.ManualInstalls[item.Key])
			case "Web Apps (Chromium based)":
				if !manager.IsInstalled("chromium") {
					fmt.Println("Installing Chromium for WebApps...")
					manager.Install("chromium")
				}

				if item.Key == "custom" {
					var name, url string
					form := huh.NewForm(
						huh.NewGroup(
							huh.NewInput().
								Title("App Name").
								Placeholder("My Web App").
								Value(&name),
							huh.NewInput().
								Title("App URL").
								Placeholder("https://example.com").
								Value(&url),
						),
					)
					err := form.Run()
					if err != nil {
						fmt.Printf("Error getting custom web app details: %v\n", err)
						continue
					}
					if name != "" && url != "" {
						if !strings.HasPrefix(url, "http") {
							url = "https://" + url
						}
						createWebApp(name, url)
					}
				} else {
					createWebApp(item.Name, item.Key)
				}
			case "AI Tools":
				installFromConfig(manager, config.AIInstalls[item.Key])
			case "Helper Tools":
				installFromConfig(manager, config.HelperInstalls[item.Key])
			}
		}
	} else if action == "r" {
		fmt.Println("\n--- Removing Selected Software ---")
		for _, item := range results {
			if !item.Selected {
				continue
			}

			switch item.Category {
			case "General Software":
				manager.Remove(item.Key)
			case "Manual Installs (curl/fsSL)":
				fmt.Printf("Manual removal not yet implemented for: %s\n", item.Name)
			case "Web Apps (Chromium based)":
				if item.Key != "custom" {
					removeWebApp(item.Name)
				}
			case "AI Tools":
				fmt.Printf("Manual removal not yet implemented for: %s\n", item.Name)
			case "Helper Tools":
				fmt.Printf("Manual removal not yet implemented for: %s\n", item.Name)
			}
		}
	}

	return nil
}

func removeWebApp(name string) {
	fmt.Printf("Removing WebApp: %s\n", name)
	home := os.Getenv("HOME")
	filePath := filepath.Join(home, ".local/share/applications", strings.ToLower(name)+".desktop")
	os.Remove(filePath)
	runSimpleCmd("update-desktop-database ~/.local/share/applications")
}

func installFromConfig(manager pkgmanager.PackageManager, inst config.InstallConfig) {
	if inst.Check != "" {
		if strings.HasPrefix(inst.Check, "~/") {
			home := os.Getenv("HOME")
			path := filepath.Join(home, inst.Check[2:])
			if _, err := os.Stat(path); err == nil {
				fmt.Printf("%s is already installed.\n", inst.Name)
				return
			}
		} else {
			if _, err := exec.LookPath(inst.Check); err == nil {
				fmt.Printf("%s is already installed.\n", inst.Name)
				return
			}
		}
	}

	// Install dependencies if any
	if len(inst.Deps) > 0 {
		fmt.Printf("Checking dependencies for %s: %v\n", inst.Name, inst.Deps)
		var missingDeps []string
		for _, dep := range inst.Deps {
			if !manager.IsInstalled(dep) {
				missingDeps = append(missingDeps, dep)
			}
		}

		if len(missingDeps) > 0 {
			fmt.Printf("Installing missing dependencies: %v\n", missingDeps)
			if err := manager.Install(missingDeps...); err != nil {
				fmt.Printf("Warning: Failed to install dependencies: %v\n", err)
			}
		}
	}

	fmt.Printf("Installing %s...\n", inst.Name)
	runSimpleCmd(inst.Command)
}



func createWebApp(name, url string) {
	fmt.Printf("Creating WebApp: %s\n", name)

	home := os.Getenv("HOME")
	iconDir := filepath.Join(home, ".local/share/applications/icons")
	os.MkdirAll(iconDir, 0755)
	iconPath := filepath.Join(iconDir, strings.ToLower(name)+".png")

	// Download favicon
	faviconURL := fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=128", url)
	if err := downloadIcon(faviconURL, iconPath); err != nil {
		fmt.Printf("Warning: Could not download icon: %v\n", err)
		iconPath = "chromium" // Fallback to chromium icon
	}

	desktopFile := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=%s
Comment=%s
Exec=chromium --app=%s
Icon=%s
Terminal=false
StartupNotify=true
Categories=Network;WebBrowser;
`, name, name, url, iconPath)

	filePath := filepath.Join(home, ".local/share/applications", strings.ToLower(name)+".desktop")

	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	err = os.WriteFile(filePath, []byte(desktopFile), 0755) // Mark as executable
	if err != nil {
		fmt.Printf("Error writing desktop file: %v\n", err)
	}

	runSimpleCmd("update-desktop-database ~/.local/share/applications")
}

func downloadIcon(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func runSimpleCmd(shellCmd string) {
	cmd := exec.Command("bash", "-c", shellCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
