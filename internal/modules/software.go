package modules

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/config"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
	"github.com/rakesh/linutils-rakesh/internal/tui"
)

func InstallSoftware(manager pkgmanager.PackageManager, sysInfo system.Info, items []tui.ListItem) ([]tui.ListItem, error) {
	if len(items) == 0 {
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

		// 6. Flatpak Installs
		for key, inst := range config.FlatpakInstalls {
			items = append(items, tui.ListItem{Key: key, Name: inst.Name, Category: "Flatpak Installs"})
		}
	}

	for {
		action, results, err := tui.RunListUI("Software Installer", items)
		if err != nil {
			return results, err
		}

		if action == "" {
			break
		}

		// Sync selection state
		items = results

		selectedCount := 0
		for _, item := range results {
			if item.Selected {
				selectedCount++
			}
		}

		if selectedCount == 0 {
			fmt.Println("No items selected.")
			continue
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
					installFromConfig(manager, sysInfo, config.ManualInstalls[item.Key])
				case "Web Apps (Chromium based)":
					if !manager.IsInstalled("chromium") {
						fmt.Println("Installing Chromium for WebApps...")
						manager.Install("chromium")
					}

					if item.Key == "custom" {
						var name, url, iconURL string
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
								huh.NewInput().
									Title("Icon URL (Optional)").
									Placeholder("https://dashboardicons.com/...").
									Value(&iconURL),
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
							createWebApp(name, url, iconURL)
						}
					} else {
						createWebApp(item.Name, item.Key, "")
					}
				case "AI Tools":
					installFromConfig(manager, sysInfo, config.AIInstalls[item.Key])
				case "Helper Tools":
					installFromConfig(manager, sysInfo, config.HelperInstalls[item.Key])
				case "Flatpak Installs":
					if !isFlatpakReady() {
						fmt.Println("Flatpak not ready. Setting up Flatpak first...")
						if err := SetupFlatpak(manager, sysInfo); err != nil {
							fmt.Printf("Error setting up Flatpak: %v\n", err)
							continue
						}
					}
					installFromConfig(manager, sysInfo, config.FlatpakInstalls[item.Key])
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
					if inst, ok := config.ManualInstalls[item.Key]; ok && len(inst.Remove) > 0 {
						fmt.Printf("Removing %s...\n", item.Name)
						runCommands(inst.Remove)
					} else {
						fmt.Printf("Manual removal not yet implemented for: %s\n", item.Name)
					}
				case "Web Apps (Chromium based)":
					if item.Key != "custom" {
						removeWebApp(item.Name)
					}
				case "AI Tools":
					fmt.Printf("Manual removal not yet implemented for: %s\n", item.Name)
				case "Helper Tools":
					fmt.Printf("Manual removal not yet implemented for: %s\n", item.Name)
				case "Flatpak Installs":
					if inst, ok := config.FlatpakInstalls[item.Key]; ok && len(inst.Remove) > 0 {
						fmt.Printf("Removing %s...\n", item.Name)
						runCommands(inst.Remove)
					} else {
						fmt.Printf("Flatpak removal not yet implemented for: %s\n", item.Name)
					}
				}
			}
		}

		fmt.Println("\nOperation complete! Press Enter to return to software list...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	return items, nil
}

func isFlatpakReady() bool {
	_, err := exec.LookPath("flatpak")
	if err != nil {
		return false
	}
	// Also check if flathub is added
	cmd := exec.Command("flatpak", "remotes")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "flathub")
}

func removeWebApp(name string) {
	fmt.Printf("Removing WebApp: %s\n", name)
	home := os.Getenv("HOME")
	kebabName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	filePath := filepath.Join(home, ".local/share/applications", kebabName+".desktop")
	os.Remove(filePath)
	runSimpleCmd("update-desktop-database ~/.local/share/applications")
}

func osGroup(osName string) string {
	switch osName {
	case "debian", "ubuntu", "pop", "linuxmint":
		return "apt"
	case "arch", "manjaro":
		return "arch"
	case "fedora":
		return "fedora"
	}
	return "default"
}

func installFromConfig(manager pkgmanager.PackageManager, sysInfo system.Info, inst config.InstallConfig) {
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

	// Pick commands: CommandByOS takes priority, fall back to Command
	cmds := inst.Command
	group := osGroup(sysInfo.OS)
	if osCmds, ok := inst.CommandByOS[group]; ok {
		cmds = osCmds
	} else if osCmds, ok := inst.CommandByOS[sysInfo.OS]; ok {
		cmds = osCmds
	}

	fmt.Printf("Installing %s...\n", inst.Name)
	runCommands(cmds)
}



func createWebApp(name, url, iconURL string) {
	fmt.Printf("Creating WebApp: %s\n", name)

	home := os.Getenv("HOME")
	iconDir := filepath.Join(home, ".local/share/applications/icons")
	os.MkdirAll(iconDir, 0755)
	
	// Base icon name without extension
	iconBaseName := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	iconPath := "" // Will be set once we know the extension

	// Normalize URL for icon fetching
	domain := url
	if strings.Contains(url, "://") {
		parts := strings.Split(url, "/")
		if len(parts) > 2 {
			domain = parts[2]
		}
	}

	// Generate a kebab-case name for dashboard-icons guessing
	kebabName := strings.ReplaceAll(strings.ToLower(name), " ", "-")

	fetched := false

	// Check for builtin icons first
	if svgContent, ok := builtinIcons[kebabName]; ok {
		fmt.Printf("Using builtin high-quality SVG for: %s\n", name)
		iconPath = filepath.Join(iconDir, iconBaseName+".svg")
		if err := os.WriteFile(iconPath, []byte(svgContent), 0644); err == nil {
			fetched = true
		}
	}

	if !fetched {
		// Icon fetching strategy
		type iconSource struct {
			url string
			ext string
		}

		sources := []iconSource{}
		if iconURL != "" {
			ext := ".png"
			if strings.HasSuffix(strings.ToLower(iconURL), ".svg") {
				ext = ".svg"
			}
			sources = append(sources, iconSource{url: iconURL, ext: ext})
		}

		// Priority 1: Automated Dashboard Icons (SVG then PNG)
		sources = append(sources,
			iconSource{url: fmt.Sprintf("https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/svg/%s.svg", kebabName), ext: ".svg"},
			iconSource{url: fmt.Sprintf("https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/%s.png", kebabName), ext: ".png"},
		)

		// Priority 2: High-res Google Favicon & DuckDuckGo Fallback
		sources = append(sources,
			iconSource{url: fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=512", domain), ext: ".png"},
			iconSource{url: fmt.Sprintf("https://icons.duckduckgo.com/ip3/%s.ico", domain), ext: ".ico"},
		)

		for _, src := range sources {
			tempPath := filepath.Join(iconDir, iconBaseName+src.ext)
			fmt.Printf("Trying icon source: %s\n", src.url)
			if err := downloadIcon(src.url, tempPath); err == nil {
				if info, err := os.Stat(tempPath); err == nil && info.Size() > 500 {
					fetched = true
					iconPath = tempPath
					break
				} else {
					os.Remove(tempPath) // Clean up empty/small files
				}
			}
		}

		if !fetched {
			fmt.Println("Warning: Could not fetch high-quality icon, falling back to chromium icon")
			iconPath = "chromium"
		}
	}

	chromiumBin := findChromiumBinary()
	appID := strings.ReplaceAll(name, " ", "")
	desktopFile := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=%s
Comment=%s
Exec=%s --class=%s --app=%s
Icon=%s
Terminal=false
StartupNotify=true
StartupWMClass=%s
Categories=Network;WebBrowser;
`, name, name, chromiumBin, appID, url, iconPath, appID)

	filePath := filepath.Join(home, ".local/share/applications", kebabName+".desktop")

	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
	// Ensure both directories have correct permissions
	appsDir := filepath.Dir(filePath)
	os.MkdirAll(appsDir, 0755)
	os.Chmod(appsDir, 0755)
	os.Chmod(iconDir, 0755)

	err = os.WriteFile(filePath, []byte(desktopFile), 0644)
	if err != nil {
		fmt.Printf("Error writing desktop file: %v\n", err)
	}

	fmt.Println("Refreshing desktop environment...")
	runSimpleCmd("update-desktop-database ~/.local/share/applications")
	runSimpleCmd("gtk-update-icon-cache ~/.local/share/icons >/dev/null 2>&1 || true")
}

func findChromiumBinary() string {
	binaries := []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable", "brave-browser"}
	for _, bin := range binaries {
		if _, err := exec.LookPath(bin); err == nil {
			return bin
		}
	}
	return "chromium" // Fallback
}


func downloadIcon(url, path string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
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
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func runCommands(commands []string) {
	for _, cmdStr := range commands {
		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Step failed, aborting: %v\n", err)
			return
		}
	}
}
