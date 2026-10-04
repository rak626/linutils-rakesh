package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Gruvbox Dark Medium — ASCII only, normal terminal font.
// No Nerd Font icons. All markers are plain ASCII.
var (
	GruvBg     = lipgloss.Color("#282828")
	GruvBg0H   = lipgloss.Color("#1d2021")
	GruvBg1    = lipgloss.Color("#3c3836")
	GruvBg2    = lipgloss.Color("#504945")
	GruvFg     = lipgloss.Color("#ebdbb2")
	GruvFg0    = lipgloss.Color("#fbf1c7")
	GruvGray   = lipgloss.Color("#928374")
	GruvOrange = lipgloss.Color("#fe8019")
	GruvGreen  = lipgloss.Color("#b8bb26")
	GruvYellow = lipgloss.Color("#fabd2f")
	GruvRed    = lipgloss.Color("#fb4934")
	GruvBlue   = lipgloss.Color("#83a598")
	GruvAqua   = lipgloss.Color("#8ec07c")
	GruvPurple = lipgloss.Color("#d3869b")

	ThemeHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(GruvOrange).
			Background(GruvBg0H).
			Padding(0, 2)

	ThemeCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(GruvBg2).
			Background(GruvBg)

	ThemeSidebar = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(GruvBg2).
			Background(GruvBg)

	ThemeMain = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(GruvOrange).
			Background(GruvBg)

	ThemeCursor = lipgloss.NewStyle().
			Background(GruvBg1).
			Foreground(GruvOrange).
			Bold(true)

	ThemeSelected = lipgloss.NewStyle().
			Foreground(GruvGreen).
			Bold(true)

	ThemeDim = lipgloss.NewStyle().
			Foreground(GruvGray)

	ThemeFooter = lipgloss.NewStyle().
			Background(GruvBg0H).
			Foreground(GruvFg).
			Padding(0, 2).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(GruvBg2)

	ThemeKey     = lipgloss.NewStyle().Foreground(GruvYellow).Bold(true)
	ThemeLabel   = lipgloss.NewStyle().Foreground(GruvOrange).Bold(true)
	ThemeSysKey  = lipgloss.NewStyle().Foreground(GruvAqua).Bold(true)
	ThemeSysVal  = lipgloss.NewStyle().Foreground(GruvFg)
	ThemeSuccess = lipgloss.NewStyle().Foreground(GruvGreen).Bold(true)
	ThemeDanger  = lipgloss.NewStyle().Foreground(GruvRed).Bold(true)
)

// ASCII markers — no Nerd Font dependency.
const (
	MarkSelected   = "[*]"
	MarkUnselected = "[ ]"
	MarkCursor     = "> "
	MarkSection    = "# "
	MarkOK         = "[OK]"
	MarkFail       = "[FAIL]"
	MarkRun        = "[..]"
)

var spinnerFrames = []string{"|", "/", "-", "\\"}

func SpinnerFrame(i int) string {
	return spinnerFrames[i%len(spinnerFrames)]
}

func ProgressBar(done, total, width int) string {
	if total <= 0 || width <= 0 {
		return ""
	}
	filled := done * width / total
	if filled > width {
		filled = width
	}
	bar := ""
	for i := 0; i < filled; i++ {
		bar += "="
	}
	for i := filled; i < width; i++ {
		bar += "-"
	}
	return "[" + bar + "]"
}
