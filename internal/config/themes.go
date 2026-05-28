package config

type ThemeConfig struct {
	Name      string
	Alacritty string
	Zed       string
	Neovim    string
	Vim       string
	GTK       string
	Ulauncher string
	Starship  string
	VSCodium  string
	Ghostty   string
	Btop      string
	Kitty     string
	Icons     string
	Cursor    string
	Hyprland  string // Content for active_theme.conf
	I3        string // Content for active_theme.i3
	Waybar    string // CSS variables
	Mako      string // Config content
	Hyprlock  string // Config content
	SwayOSD   string // CSS content
	GNOME     struct {
		ShellTheme string
		Wallpaper  string
	}
}

var GlobalThemes = []ThemeConfig{
	{
		Name:      "Rose Pine Moon",
		Alacritty: "rose-pine-moon.toml",
		Zed:       "Rosé Pine Moon",
		Neovim:    "rose-pine",
		Vim:       "rose-pine",
		GTK:       "rose-pine-moon-gtk",
		Ulauncher: "rose-pine-moon",
		Starship:  "rose-pine-moon",
		VSCodium:  "Rosé Pine Moon",
		Ghostty:   "rose-pine-moon",
		Btop:      "rose-pine-moon",
		Kitty:     "Rosé Pine Moon",
		Icons:     "rose-pine",
		Cursor:    "rose-pine-cursor",
		Hyprland: `
$base = 0xff232136
$surface = 0xff2a273f
$overlay = 0xff393552
$muted = 0xff6e6a86
$subtle = 0xff908caa
$text = 0xffe0def4
$love = 0xffeb6f92
$gold = 0xfff6c177
$rose = 0xffea9a97
$pine = 0xff3e8fb0
$foam = 0xff9ccfd8
$iris = 0xffc4a7e7

general {
    col.active_border = $iris $pine 45deg
    col.inactive_border = $muted
}
`,
		Waybar: `
@define-color base #232136;
@define-color surface #2a273f;
@define-color overlay #393552;
@define-color text #e0def4;
@define-color love #eb6f92;
@define-color gold #f6c177;
@define-color rose #ea9a97;
@define-color pine #3e8fb0;
@define-color foam #9ccfd8;
@define-color iris #c4a7e7;
`,
		Mako: `
background-color=#232136
text-color=#e0def4
border-color=#c4a7e7
`,
		Hyprlock: `
$text = rgb(e0def4)
$accent = rgb(c4a7e7)
`,
		SwayOSD: `
@define-color accent #c4a7e7;
@define-color bg #232136;
`,
	},
	{
		Name:      "Catppuccin Macchiato",
		Alacritty: "catppuccin-macchiato.toml",
		Zed:       "Catppuccin Macchiato",
		Neovim:    "catppuccin",
		Vim:       "catppuccin",
		GTK:       "catppuccin-macchiato-gtk",
		Ulauncher: "catppuccin-macchiato",
		Starship:  "catppuccin-macchiato",
		VSCodium:  "Catppuccin Macchiato",
		Ghostty:   "catppuccin-macchiato",
		Btop:      "catppuccin-macchiato",
		Kitty:     "Catppuccin Macchiato",
		Icons:     "Papirus-Dark",
		Cursor:    "catppuccin-macchiato-cursors",
		Hyprland: `
$base = 0xff24273a
$text = 0xffcad3f5
$mauve = 0xffc6a0f6
$blue = 0xff8aadf4
$surface0 = 0xff363a4f

general {
    col.active_border = $mauve $blue 45deg
    col.inactive_border = $surface0
}
`,
		Waybar: `
@define-color base #24273a;
@define-color text #cad3f5;
@define-color mauve #c6a0f6;
@define-color blue #8aadf4;
`,
		Mako: `
background-color=#24273a
text-color=#cad3f5
border-color=#c6a0f6
`,
		Hyprlock: `
$text = rgb(cad3f5)
$accent = rgb(c6a0f6)
`,
	},
	{
		Name:      "Everforest",
		Alacritty: "everforest.toml",
		Zed:       "Everforest Dark",
		Neovim:    "everforest",
		Vim:       "everforest",
		GTK:       "Everforest-Dark-B",
		Ulauncher: "everforest-dark",
		Starship:  "everforest",
		VSCodium:  "Everforest Dark",
		Ghostty:   "everforest",
		Btop:      "everforest",
		Kitty:     "Everforest Dark",
		Icons:     "Papirus-Dark",
		Cursor:    "Bibata-Modern-Classic",
		Hyprland: `
$bg = 0xff2d353b
$fg = 0xffd3c6aa
$green = 0xffa7c080
$blue = 0xff7fbbb3

general {
    col.active_border = $green $blue 45deg
    col.inactive_border = $bg
}
`,
		Waybar: `
@define-color base #2d353b;
@define-color text #d3c6aa;
@define-color green #a7c080;
@define-color blue #7fbbb3;
`,
	},
	{
		Name:      "One Dark",
		Alacritty: "one-dark.toml",
		Zed:       "One Dark",
		Neovim:    "onedark",
		Vim:       "onedark",
		GTK:       "one-dark-gtk",
		Ulauncher: "one-dark",
		Starship:  "one-dark",
		VSCodium:  "One Dark",
		Ghostty:   "one-dark",
		Btop:      "one-dark",
		Kitty:     "One Dark",
		Icons:     "Papirus-Dark",
		Cursor:    "Bibata-Modern-Classic",
		Hyprland: `
$bg = 0xff282c34
$fg = 0xffabb2bf
$blue = 0xff61afef
$magenta = 0xffc678dd
$grey = 0xff5c6370

general {
    col.active_border = $blue $magenta 45deg
    col.inactive_border = $grey
}
`,
		Waybar: `
@define-color base #282c34;
@define-color text #abb2bf;
@define-color blue #61afef;
@define-color magenta #c678dd;
`,
	},
	{
		Name:      "Gruvbox Dark",
		Alacritty: "gruvbox-dark.toml",
		Zed:       "Gruvbox Dark Medium",
		Neovim:    "gruvbox",
		Vim:       "gruvbox",
		GTK:       "gruvbox-dark-gtk",
		Ulauncher: "gruvbox-dark",
		Starship:  "gruvbox-dark",
		VSCodium:  "Gruvbox Dark Medium",
		Ghostty:   "gruvbox-dark",
		Btop:      "gruvbox-dark",
		Kitty:     "Gruvbox Dark",
		Icons:     "Gruvbox-Plus-Dark",
		Cursor:    "Gruvbox-Cursor",
		Hyprland: `
$bg = 0xff282828
$fg = 0xffebdbb2
$aqua = 0xff8ec07c
$blue = 0xff83a598

general {
    col.active_border = $aqua $blue 45deg
    col.inactive_border = $bg
}
`,
		Waybar: `
@define-color base #282828;
@define-color text #ebdbb2;
@define-color aqua #8ec07c;
@define-color blue #83a598;
`,
	},
	{
		Name:      "Miasma",
		Alacritty: "miasma.toml",
		Zed:       "Miasma",
		Neovim:    "miasma",
		Vim:       "miasma",
		GTK:       "miasma-gtk",
		Ulauncher: "miasma",
		Starship:  "miasma",
		VSCodium:  "Miasma",
		Ghostty:   "miasma",
		Btop:      "miasma",
		Kitty:     "Miasma",
		Icons:     "Papirus-Dark",
		Cursor:    "Bibata-Modern-Classic",
		Hyprland: `
$bg = 0xff222222
$fg = 0xffc2c2b0
$green = 0xff5f875f
$blue = 0xff78824b

general {
    col.active_border = $green $blue 45deg
    col.inactive_border = $bg
}
`,
		Waybar: `
@define-color base #222222;
@define-color text #c2c2b0;
@define-color green #5f875f;
@define-color blue #78824b;
`,
	},
}
