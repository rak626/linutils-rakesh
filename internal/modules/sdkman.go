package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func SetupSDKMan() {
	fmt.Println("\n--- Setting up SDKMAN ---")

	home := os.Getenv("HOME")
	sdkmanDir := filepath.Join(home, ".sdkman")

	if _, err := os.Stat(sdkmanDir); os.IsNotExist(err) {
		fmt.Println("Installing SDKMAN...")
		cmd := exec.Command("bash", "-c", `curl -s "https://get.sdkman.io" | bash`)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error installing SDKMAN: %v\n", err)
			return
		}
	} else {
		fmt.Println("SDKMAN already installed.")
	}

	fmt.Println("Installing Java and Maven via SDKMAN...")

	sdkInit := filepath.Join(sdkmanDir, "bin", "sdkman-init.sh")
	script := fmt.Sprintf(`source "%s"

sdk install java 21-temurin

sdk install java 25-temurin

latest_lts=$(sdk list java 2>&1 | grep -i tem | grep -v '\*' | awk '{print $NF}' | while read ver; do
  major="${ver%%.*}"
  case "$major" in
    8|11|17|21) echo "$ver" ;;
  esac
done | sort -V | tail -1)
[ -n "$latest_lts" ] && sdk install java "$latest_lts"

latest_installed=$(sdk list java 2>&1 | grep installed | grep -i tem | awk '{print $NF}' | sort -V | tail -1)
[ -n "$latest_installed" ] && sdk default java "$latest_installed"

sdk install maven
`, sdkInit)

	cmd := exec.Command("bash", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error during SDKMAN setup: %v\n", err)
		return
	}

	fmt.Println("SDKMAN setup complete!")
	time.Sleep(1500 * time.Millisecond)
}
