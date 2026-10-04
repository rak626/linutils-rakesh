package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rak626/linutils-rakesh/internal/system"
)

var (
	// Gruvbox Dark Medium aliases — single source is theme.go.
	fgColor      = GruvFg
	accentColor  = GruvOrange
	successColor = GruvGreen
	warningColor = GruvYellow
	grayColor    = GruvBg2
	dimColor     = GruvGray
	white        = GruvFg0

	// Styles
	headerStyle = ThemeHeader.MarginBottom(1)

	sidebarStyle = ThemeSidebar.Width(38)

	mainContentStyle = ThemeMain

	tabStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(GruvFg)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(GruvOrange).
			Padding(0, 1)

	selectedItemStyle = ThemeSelected

	cursorItemStyle = ThemeCursor

	systemInfoStyle = lipgloss.NewStyle().
			MarginTop(1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(GruvBg2).
			PaddingTop(1)

	sysKeyStyle = ThemeSysKey
	sysValStyle = ThemeSysVal

	footerStyle = ThemeFooter.MarginTop(1)

	helpLabelStyle = ThemeLabel
	helpKeyStyle   = ThemeKey
	helpTextStyle  = ThemeDim

	searchInputStyle = lipgloss.NewStyle().
				Foreground(GruvOrange).
				Border(lipgloss.NormalBorder()).
				BorderForeground(GruvBg2).
				Padding(0, 1)
)

type ListItem struct {
	Key         string
	Name        string
	Category    string
	Description string
	Selected    bool
	Hidden      bool
}

type ListModel struct {
	Title       string
	Description string
	SysInfo     system.Info
	Items       []ListItem
	Filtered    []int // indices of original Items
	Cursor      int
	Action      string // "r" for remove, "i" for install, "" for none
	Finished    bool
	ShowHelp    bool

	Tabs         []string
	ActiveTab    int
	SearchInput  textinput.Model
	ScrollOffset int
	Width        int
	Height       int
}

func (m ListModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Action = "quit"
			m.Finished = true
			return m, tea.Quit
		case "?", "h":
			if !m.SearchInput.Focused() {
				m.ShowHelp = !m.ShowHelp
				return m, nil
			}
		case "esc":
			if m.ShowHelp {
				m.ShowHelp = false
				return m, nil
			}
			if m.SearchInput.Focused() {
				m.SearchInput.Blur()
				m.SearchInput.SetValue("")
				m.filterItems()
			} else {
				m.Action = "back"
				m.Finished = true
				return m, tea.Quit
			}
		case "up", "k":
			if len(m.Filtered) == 0 {
				break
			}
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Filtered) - 1
			}
			m.fixScroll()
		case "down", "j":
			if len(m.Filtered) == 0 {
				break
			}
			if m.Cursor < len(m.Filtered)-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}
			m.fixScroll()
		case "tab":
			if len(m.Tabs) == 0 {
				break
			}
			m.ActiveTab = (m.ActiveTab + 1) % len(m.Tabs)
			m.Cursor = 0
			m.ScrollOffset = 0
			m.filterItems()
		case "shift+tab":
			if len(m.Tabs) == 0 {
				break
			}
			m.ActiveTab = (m.ActiveTab - 1 + len(m.Tabs)) % len(m.Tabs)
			m.Cursor = 0
			m.ScrollOffset = 0
			m.filterItems()
		case " ":
			if len(m.Filtered) > 0 {
				idx := m.Filtered[m.Cursor]
				m.Items[idx].Selected = !m.Items[idx].Selected
			}
		case "/":
			if !m.SearchInput.Focused() {
				m.SearchInput.Focus()
				return m, nil
			}
		case "ctrl+v":
			// Toggle selection for all items in the filtered list
			allSelected := true
			for _, idx := range m.Filtered {
				if !m.Items[idx].Selected {
					allSelected = false
					break
				}
			}
			for _, idx := range m.Filtered {
				m.Items[idx].Selected = !allSelected
			}
		case "enter":
			if m.SearchInput.Focused() {
				m.SearchInput.Blur()
			} else {
				// If no items are selected via Space, we only run the current item
				// but we don't mark it as "Selected" in the persistent list.
				anySelected := false
				for _, item := range m.Items {
					if item.Selected {
						anySelected = true
						break
					}
				}

				if !anySelected && len(m.Filtered) > 0 {
					// We'll signal this by setting a special state or returning it
					// For now, let's just make sure the calling code knows.
					m.Items[m.Filtered[m.Cursor]].Selected = true
					m.Action = "i_single"
				} else {
					m.Action = "i"
				}

				m.Finished = true
				return m, tea.Quit
			}
		case "r", "R":
			if !m.SearchInput.Focused() {
				m.Action = "r"
				m.Finished = true
				return m, tea.Quit
			}
		}
	}

	if m.SearchInput.Focused() {
		m.SearchInput, cmd = m.SearchInput.Update(msg)
		m.filterItems()
	}

	return m, cmd
}

func (m *ListModel) fixScroll() {
	if m.Height == 0 {
		return
	}
	visibleHeight := m.Height - 16 // Buffer for header, footer, search, padding
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	if m.Cursor < m.ScrollOffset {
		m.ScrollOffset = m.Cursor
	} else if m.Cursor >= m.ScrollOffset+visibleHeight {
		m.ScrollOffset = m.Cursor - visibleHeight + 1
	}
}

func (m *ListModel) filterItems() {
	m.Filtered = []int{}
	searchTerm := strings.ToLower(m.SearchInput.Value())
	currentCategory := "All"
	if len(m.Tabs) > 0 {
		currentCategory = m.Tabs[m.ActiveTab]
	}

	for i, item := range m.Items {
		if item.Hidden {
			continue
		}
		matchesCategory := currentCategory == "All" || item.Category == currentCategory
		matchesSearch := searchTerm == "" || strings.Contains(strings.ToLower(item.Name), searchTerm) || strings.Contains(strings.ToLower(item.Category), searchTerm) || strings.Contains(strings.ToLower(item.Description), searchTerm)

		if matchesCategory && matchesSearch {
			m.Filtered = append(m.Filtered, i)
		}
	}

	if m.Cursor >= len(m.Filtered) {
		m.Cursor = 0
	}
	m.fixScroll()
}

func (m ListModel) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Initializing..."
	}

	// 1. Header — ASCII only
	title := strings.ToUpper(m.Title)
	header := headerStyle.Render("[#] LINUTILS RAKESH  >  " + title)

	// 2. Sidebar (System Info) — responsive widths
	bodyHeight := m.Height - 10
	if bodyHeight < 10 {
		bodyHeight = 10
	}

	osDisplay := m.SysInfo.OS
	if m.SysInfo.OSVersion != "" && !strings.Contains(m.SysInfo.OS, m.SysInfo.OSVersion) {
		osDisplay += " " + m.SysInfo.OSVersion
	}

	deDisplay := m.SysInfo.DE
	if m.SysInfo.DEVersion != "" && !strings.Contains(m.SysInfo.DE, m.SysInfo.DEVersion) {
		deDisplay += " " + m.SysInfo.DEVersion
	}

	sysInfoContent := fmt.Sprintf("%s\n%s\n\n%s\n\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s",
		helpLabelStyle.Render("# CURRENT TASK"),
		lipgloss.NewStyle().Foreground(accentColor).Bold(true).Render(m.Title),
		helpLabelStyle.Render("# SYSTEM"),
		sysKeyStyle.Render("OS:  "), sysValStyle.Render(osDisplay),
		sysKeyStyle.Render("DE:  "), sysValStyle.Render(deDisplay),
		sysKeyStyle.Render("CPU: "), sysValStyle.Render(m.SysInfo.CPU),
		sysKeyStyle.Render("RAM: "), sysValStyle.Render(m.SysInfo.RAM),
		sysKeyStyle.Render("DISK:"), sysValStyle.Render(m.SysInfo.Disk),
		sysKeyStyle.Render("GPU: "), sysValStyle.Render(m.SysInfo.GPU),
	)

	showSidebar := m.Width >= 100
	var sidebar string
	if showSidebar {
		sidebar = sidebarStyle.
			Height(bodyHeight).
			Render(sysInfoContent)
	}

	// 3. Main Content (Presets) — ASCII only
	boxTitle := "SELECTION"
	if m.Title != "" {
		boxTitle = m.Title
	}
	listContent := helpLabelStyle.Render("# "+strings.ToUpper(boxTitle)) + "\n\n"

	visibleHeight := bodyHeight - 4
	end := m.ScrollOffset + visibleHeight
	if end > len(m.Filtered) {
		end = len(m.Filtered)
	}

	for i := m.ScrollOffset; i < end; i++ {
		idx := m.Filtered[i]
		item := m.Items[idx]

		// Selection indicator — ASCII
		checked := MarkUnselected
		if item.Selected {
			checked = MarkSelected
		}

		// Cursor indicator and item text
		line := fmt.Sprintf(" %s %s", checked, item.Name)

		if m.Cursor == i {
			w := m.Width - 12
			if !showSidebar {
				w = m.Width - 8
			} else {
				w = m.Width - 52
			}
			if w < 20 {
				w = 20
			}
			listContent += cursorItemStyle.Width(w).Render(MarkCursor+line) + "\n"
		} else {
			if item.Selected {
				listContent += selectedItemStyle.Render(" "+line) + "\n"
			} else {
				listContent += "  " + line + "\n"
			}
		}
	}

	mainWidth := m.Width - 6
	if showSidebar {
		mainWidth = m.Width - 46
	}
	if mainWidth < 20 {
		mainWidth = 20
	}
	main := mainContentStyle.
		Height(bodyHeight).
		Width(mainWidth).
		Render(listContent)

	// Combine Sidebar and Main — stack on narrow terminals
	var body string
	if showSidebar {
		body = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
	} else {
		body = main
	}

	// 4. Footer — ASCII, rune-safe truncation
	footerContent := ""
	if len(m.Filtered) > 0 && m.Cursor < len(m.Filtered) {
		item := m.Items[m.Filtered[m.Cursor]]
		desc := item.Description
		maxDesc := m.Width - 15
		if maxDesc < 10 {
			maxDesc = 10
		}
		if len([]rune(desc)) > maxDesc {
			desc = string([]rune(desc)[:maxDesc-3]) + "..."
		}
		footerContent += fmt.Sprintf("%s %s\n", helpLabelStyle.Render("DESC:"), desc)
	}

	commands := fmt.Sprintf("%s Quit  %s Move  %s Select  %s All  %s Run  %s Help  %s Search",
		helpKeyStyle.Render("[q]"),
		helpKeyStyle.Render("[j/k]"),
		helpKeyStyle.Render("[Space]"),
		helpKeyStyle.Render("[Ctrl+v]"),
		helpKeyStyle.Render("[Enter]"),
		helpKeyStyle.Render("[?]"),
		helpKeyStyle.Render("[/]"),
	)
	if m.Width < 100 {
		commands += "\n" + ThemeDim.Render("Small screen: sidebar hidden. All keys still work.")
	}

	footer := footerStyle.Width(m.Width - 4).Render(footerContent + commands)

	view := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
	if m.ShowHelp {
		help := ThemeCard.
			Width(60).
			Render(ThemeLabel.Render("HELP") + "\n\n" +
				"j/k or Up/Down : move\n" +
				"Space          : select\n" +
				"Enter          : run (single if none selected)\n" +
				"Ctrl+v         : toggle all visible\n" +
				"/ then type    : search, Esc clears\n" +
				"Tab            : next category\n" +
				"r              : remove mode (debloat)\n" +
				"q / Esc        : back / quit\n\n" +
				"Press ? or Esc to close")
		_ = help
		// Overlay help by appending below — keeps it simple and ASCII-safe.
		view = lipgloss.JoinVertical(lipgloss.Left, view, "", ThemeDim.Render("HELP OPEN: j/k move, Space select, Enter run, / search, q back. Press ? to close."))
	}

	return view
}

func RunListUI(title string, items []ListItem) (string, []ListItem, error) {
	return RunListUIWithDesc(title, "", items)
}

func RunListUIWithDesc(title, desc string, items []ListItem) (string, []ListItem, error) {
	sysInfo := system.GetSystemInfo()

	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 50
	ti.Width = 30

	// Extract unique categories for tabs
	categories := []string{"All"}
	catMap := make(map[string]bool)
	for _, item := range items {
		if item.Category != "" && !catMap[item.Category] {
			categories = append(categories, item.Category)
			catMap[item.Category] = true
		}
	}

	m := ListModel{
		Title:       title,
		Description: desc,
		SysInfo:     sysInfo,
		Items:       items,
		Tabs:        categories,
		ActiveTab:   0,
		SearchInput: ti,
	}
	m.filterItems()

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", nil, err
	}

	m = finalModel.(ListModel)
	return m.Action, m.Items, nil
}
