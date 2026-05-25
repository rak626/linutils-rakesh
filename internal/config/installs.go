package config

type InstallConfig struct {
	Name    string
	Command string
	Check   string // Path or binary name to check for existence
}

var ManualInstalls = map[string]InstallConfig{
	"zed": {
		Name:    "Zed Editor",
		Command: "curl -f https://zed.dev/install.sh | sh",
		Check:   "zed",
	},
	"sdkman": {
		Name:    "SDKMAN!",
		Command: "curl -s \"https://get.sdkman.io\" | bash",
		Check:   "~/.sdkman",
	},
	"nvm": {
		Name:    "NVM (Node Version Manager)",
		Command: "curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash",
		Check:   "~/.nvm",
	},
	"starship": {
		Name:    "Starship Prompt",
		Command: "curl -sS https://starship.rs/install.sh | sh -s -- -y",
		Check:   "starship",
	},
	"bun": {
		Name:    "Bun Runtime",
		Command: "curl -fsSL https://bun.sh/install | bash",
		Check:   "bun",
	},
}

var AIInstalls = map[string]InstallConfig{
	"gemini-cli": {
		Name:    "Gemini CLI",
		Command: "npm install -g @google/gemini-cli",
		Check:   "gemini",
	},
	"claude-code": {
		Name:    "Claude Code",
		Command: "npm install -g @anthropic-ai/claude-code",
		Check:   "claude",
	},
}

var HelperInstalls = map[string]InstallConfig{
	"rtk": {
		Name:    "RTK (Rust Token Killer)",
		Command: "curl -fsSL https://raw.githubusercontent.com/rtk-ai/rtk/refs/heads/master/install.sh | sh",
		Check:   "rtk",
	},
	"graphify": {
		Name:    "Graphify",
		Command: "echo 'Graphify installation placeholder'",
		Check:   "graphify",
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

