package config

type InstallConfig struct {
	Name        string
	Command     []string            // Shell commands run sequentially (fallback)
	CommandByOS map[string][]string // Per-OS commands (key: "apt", "arch", "fedora"), overrides Command
	Remove      []string            // Shell commands run sequentially (optional)
	Check       string              // Path or binary name to check for existence
	Deps        []string            // System packages needed
}

var ManualInstalls = map[string]InstallConfig{
	"github-desktop": {
		Name: "GitHub Desktop",
		Command: []string{
			"sudo rpm --import https://mirror.mwt.me/shiftkey-desktop/gpgkey",
			`sudo sh -c 'echo -e "[mwt-packages]\nname=GitHub Desktop\nbaseurl=https://mirror.mwt.me/shiftkey-desktop/rpm\nenabled=1\ngpgcheck=1\nrepo_gpgcheck=1\ngpgkey=https://mirror.mwt.me/shiftkey-desktop/gpgkey" > /etc/yum.repos.d/mwt-packages.repo'`,
			"sudo dnf install -y github-desktop",
		},
		Remove: []string{
			"sudo dnf remove -y github-desktop",
			"sudo rm -f /etc/yum.repos.d/mwt-packages.repo",
		},
		Check: "github-desktop",
	},
	"brave": {
		Name: "Brave Browser",
		CommandByOS: map[string][]string{
			"apt": {
				"sudo curl -fsSLo /usr/share/keyrings/brave-browser-archive-keyring.gpg https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg",
				"sudo curl -fsSLo /etc/apt/sources.list.d/brave-browser-release.sources https://brave-browser-apt-release.s3.brave.com/brave-browser.sources",
				"sudo apt update",
				"sudo apt install -y brave-browser",
			},
			"fedora": {
				"sudo dnf install -y dnf-plugins-core",
				"sudo dnf config-manager addrepo --from-repofile=https://brave-browser-rpm-release.s3.brave.com/brave-browser.repo",
				"sudo dnf install -y brave-browser",
			},
			"arch": {
				"curl -fsS https://dl.brave.com/install.sh | sh",
			},
		},
		Remove: []string{
			"sudo apt remove -y brave-browser 2>/dev/null",
			"sudo dnf remove -y brave-browser 2>/dev/null",
			"sudo pacman -Rs --noconfirm brave-browser 2>/dev/null",
		},
		Check: "brave-browser",
		Deps:  []string{"curl"},
	},
	"zed": {
		Name:    "Zed Editor",
		Command: []string{"curl -f https://zed.dev/install.sh | sh"},
		Check:   "zed",
	},
	"sdkman": {
		Name:    "SDKMAN!",
		Command: []string{`curl -s "https://get.sdkman.io" | bash`},
		Check:   "~/.sdkman",
	},
	"nvm": {
		Name:    "NVM (Node Version Manager)",
		Command: []string{"curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash"},
		Check:   "~/.nvm",
	},
	"starship": {
		Name:    "Starship Prompt",
		Command: []string{"curl -sS https://starship.rs/install.sh | sh -s -- -y"},
		Check:   "starship",
	},
	"bun": {
		Name:    "Bun Runtime",
		Command: []string{"curl -fsSL https://bun.sh/install | bash"},
		Check:   "bun",
	},
	"cliamp": {
		Name:    "cliamp (Music Player)",
		Command: []string{"curl -fsSL https://raw.githubusercontent.com/bjarneo/cliamp/HEAD/install.sh | sh"},
		Check:   "cliamp",
		Deps:    []string{"ffmpeg", "yt-dlp"},
	},
	"dbeaver": {
		Name: "DBeaver Community (Database Manager)",
		CommandByOS: map[string][]string{
			"fedora": {
				`wget -O /tmp/dbeaver.rpm "https://dbeaver.io/files/dbeaver-ce-latest-linux-x86_64.rpm"`,
				"sudo rpm -ivh /tmp/dbeaver.rpm",
				"rm /tmp/dbeaver.rpm",
			},
			"apt": {
				`wget -O /tmp/dbeaver.deb "https://dbeaver.io/files/dbeaver-ce-latest-linux.x86_64.deb"`,
				"sudo dpkg -i /tmp/dbeaver.deb",
				"rm /tmp/dbeaver.deb",
			},
			"arch": {
				"sudo pacman -Sy dbeaver",
			},
		},
		Remove: []string{
			"sudo dnf remove -y dbeaver-ce 2>/dev/null",
			"sudo dpkg -r dbeaver-ce 2>/dev/null",
			"sudo pacman -Rs dbeaver 2>/dev/null",
			"sudo rm -f /usr/share/applications/dbeaver.desktop",
		},
		Check: "dbeaver",
	},
	"jetbrains-toolbox": {
		Name: "JetBrains Toolbox",
		Command: []string{
			"mkdir -p ~/.local/share/JetBrains/Toolbox/bin ~/.local/bin ~/.config/JetBrains/Toolbox",
			`curl -s "https://data.services.jetbrains.com/products/releases?code=TBA&latest=true&type=release" | jq -r '.TBA[0].downloads.linux.link' | wget -O /tmp/toolbox.tar.gz -i -`,
			"tar -xzf /tmp/toolbox.tar.gz -C ~/.local/share/JetBrains/Toolbox/bin --strip-components=1",
			"rm /tmp/toolbox.tar.gz",
			`echo '{"autostart":false,"keep_running":false,"shell_scripts":{"enabled":true,"location":"$HOME/.local/bin"}}' > ~/.config/JetBrains/Toolbox/state.json`,
			`sed -i "s|\$HOME|$HOME|g" ~/.config/JetBrains/Toolbox/state.json`,
			`if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc; [[ -f ~/.zshrc ]] && echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc; fi`,
			`~/.local/share/JetBrains/Toolbox/bin/bin/jetbrains-toolbox --minimize & sleep 10; pkill -f jetbrains-toolbox`,
		},
		Remove: []string{
			"pkill -f jetbrains-toolbox 2>/dev/null",
			"rm -rf ~/.local/share/JetBrains",
			"rm -f ~/.local/share/applications/jetbrains-toolbox.desktop",
			"rm -rf ~/.cache/JetBrains",
			"rm -rf ~/.config/JetBrains",
			"if [ -d ~/.local/bin ]; then rm -f ~/.local/bin/idea ~/.local/bin/webstorm ~/.local/bin/pycharm ~/.local/bin/datagrip ~/.local/bin/clion ~/.local/bin/goland ~/.local/bin/rider ~/.local/bin/phpstorm; fi",
			`sed -i '/export PATH="\$HOME\/\.local\/bin:\$PATH"/d' ~/.bashrc`,
			`sed -i '/export PATH="\$HOME\/\.local\/bin:\$PATH"/d' ~/.zshrc`,
			`export PATH=$(echo $PATH | tr ':' '\n' | grep -v "$HOME/.local/bin" | tr '\n' ':' | sed 's/:$//')`,
		},
		Check: "~/.local/share/JetBrains/Toolbox/bin/bin/jetbrains-toolbox",
		Deps:  []string{"jq"},
	},
}

var AIInstalls = map[string]InstallConfig{
	"gemini-cli": {
		Name:    "Gemini CLI",
		Command: []string{"npm install -g @google/gemini-cli"},
		Check:   "gemini",
	},
	"claude-code": {
		Name:    "Claude Code",
		Command: []string{"npm install -g @anthropic-ai/claude-code"},
		Check:   "claude",
	},
}

var HelperInstalls = map[string]InstallConfig{
	"rtk": {
		Name:    "RTK (Rust Token Killer)",
		Command: []string{"curl -fsSL https://raw.githubusercontent.com/rtk-ai/rtk/refs/heads/master/install.sh | sh"},
		Check:   "rtk",
	},
	"graphify": {
		Name:    "Graphify",
		Command: []string{},
		Check:   "graphify",
	},
}

var FlatpakInstalls = map[string]InstallConfig{
	"obsidian": {
		Name:    "Obsidian",
		Command: []string{"sudo flatpak install -y flathub md.obsidian.Obsidian"},
		Remove:  []string{"sudo flatpak uninstall -y md.obsidian.Obsidian"},
		Check:   "md.obsidian.Obsidian",
		Deps:    []string{"flatpak"},
	},
}

var WebApps = map[string]string{
	"Discord":         "https://discord.com/app",
	"WhatsApp":        "https://web.whatsapp.com",
	"YouTube":         "https://www.youtube.com",
	"X.com":           "https://x.com",
	"Google Drive":    "https://drive.google.com",
	"Google Messages": "https://messages.google.com/web",
	"Google Meet":     "https://meet.google.com",
}
