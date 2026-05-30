package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

var (
	// Colors
	accentColor = lipgloss.Color("#7D56F4")
	white       = lipgloss.Color("#FFFFFF")
	gray        = lipgloss.Color("#626262")
	blue        = lipgloss.Color("#13BCED")
	yellow      = lipgloss.Color("#F1FA8C")

	// Styles
	sidebarStyle = lipgloss.NewStyle().
			Width(30).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor)

	mainContentStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(accentColor)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Padding(0, 1)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(white)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Padding(0, 1).
			Background(lipgloss.Color("#303030"))

	systemInfoStyle = lipgloss.NewStyle().
			MarginTop(1).
			Foreground(gray)

	sysKeyStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	sysValStyle = lipgloss.NewStyle().Foreground(white)

	helpLabelStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	helpKeyStyle   = lipgloss.NewStyle().Foreground(yellow)
)

type ListItem struct {
	Key         string
	Name        string
	Category    string
	Description string
	Selected    bool
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

	Tabs       []string
	ActiveTab  int
	SearchInput textinput.Model
	Width      int
	Height     int
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
		sidebarStyle.Height(m.Height - 10)
		mainContentStyle.Height(m.Height - 10)
		mainContentStyle.Width(m.Width - 40)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Finished = true
			return m, tea.Quit
		case "esc":
			if m.SearchInput.Focused() {
				m.SearchInput.Blur()
				m.SearchInput.SetValue("")
				m.filterItems()
			} else {
				m.Finished = true
				return m, tea.Quit
			}
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Filtered) - 1
			}
		case "down", "j":
			if m.Cursor < len(m.Filtered)-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}
		case "tab":
			m.ActiveTab = (m.ActiveTab + 1) % len(m.Tabs)
			m.Cursor = 0
			m.filterItems()
		case "shift+tab":
			m.ActiveTab = (m.ActiveTab - 1 + len(m.Tabs)) % len(m.Tabs)
			m.Cursor = 0
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
		case "enter":
			if m.SearchInput.Focused() {
				m.SearchInput.Blur()
			} else {
				// Install action
				anySelected := false
				for _, item := range m.Items {
					if item.Selected {
						anySelected = true
						break
					}
				}
				if !anySelected && len(m.Filtered) > 0 {
					m.Items[m.Filtered[m.Cursor]].Selected = true
				}
				m.Action = "i"
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

func (m *ListModel) filterItems() {
	m.Filtered = []int{}
	searchTerm := strings.ToLower(m.SearchInput.Value())
	currentCategory := m.Tabs[m.ActiveTab]

	for i, item := range m.Items {
		matchesCategory := currentCategory == "All" || item.Category == currentCategory
		matchesSearch := searchTerm == "" || strings.Contains(strings.ToLower(item.Name), searchTerm) || strings.Contains(strings.ToLower(item.Category), searchTerm)

		if matchesCategory && matchesSearch {
			m.Filtered = append(m.Filtered, i)
		}
	}

	if m.Cursor >= len(m.Filtered) {
		m.Cursor = 0
	}
}

func (m ListModel) View() string {
	// 1. Header
	header := headerStyle.Render("LINUTILS RAKESH - " + m.Title) + "\n"

	// 2. Sidebar
	sidebarContent := ""
	for i, tab := range m.Tabs {
		if i == m.ActiveTab {
			sidebarContent += activeTabStyle.Render(">> "+tab) + "\n"
		} else {
			sidebarContent += tabStyle.Render("   "+tab) + "\n"
		}
	}

	sidebarContent += "\n" + sysKeyStyle.Render("--- SYSTEM ---") + "\n"
	sidebarContent += fmt.Sprintf("%s %s\n", sysKeyStyle.Render("OS:"), sysValStyle.Render(m.SysInfo.OS))
	sidebarContent += fmt.Sprintf("%s %s\n", sysKeyStyle.Render("DE:"), sysValStyle.Render(m.SysInfo.DE))
	sidebarContent += fmt.Sprintf("%s %s\n", sysKeyStyle.Render("CPU:"), sysValStyle.Render(m.SysInfo.CPU))
	sidebarContent += fmt.Sprintf("%s %s\n", sysKeyStyle.Render("RAM:"), sysValStyle.Render(m.SysInfo.RAM))
	sidebarContent += fmt.Sprintf("%s %s\n", sysKeyStyle.Render("DISK:"), sysValStyle.Render(m.SysInfo.Disk))
	sidebarContent += fmt.Sprintf("%s %s\n", sysKeyStyle.Render("GPU:"), sysValStyle.Render(m.SysInfo.GPU))

	sidebar := sidebarStyle.Render(sidebarContent)

	// 3. Main Content
	mainContent := ""
	mainContent += fmt.Sprintf("%s %s\n\n", sysKeyStyle.Render("SEARCH"), m.SearchInput.View())

	if len(m.Filtered) == 0 {
		mainContent += "No items found.\n"
	} else {
		for i, idx := range m.Filtered {
			item := m.Items[idx]
			cursor := " "
			if m.Cursor == i {
				cursor = ">"
			}

			checked := "[ ]"
			if item.Selected {
				checked = "[*]"
			}

			line := fmt.Sprintf("%s %s %s", cursor, checked, item.Name)
			if m.Cursor == i {
				mainContent += activeTabStyle.Render(line) + "\n"
			} else {
				mainContent += tabStyle.Render(line) + "\n"
			}
		}
	}

	main := mainContentStyle.Render(mainContent)

	// Combine Sidebar and Main
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)

	// 4. Footer
	footer := "\n"
	if len(m.Filtered) > 0 && m.Cursor < len(m.Filtered) {
		item := m.Items[m.Filtered[m.Cursor]]
		desc := item.Description
		if desc == "" {
			desc = item.Category + " - " + item.Name
		}
		footer += fmt.Sprintf("%s %s\n", helpLabelStyle.Render("DESC:"), desc)
	}
	
	footer += "\n" + helpLabelStyle.Render("COMMANDS:") + "\n"
	footer += fmt.Sprintf("%s Quit  %s Navigate  %s Select  %s Install  %s Remove  %s Tab Navigation\n",
		helpKeyStyle.Render("[q]"),
		helpKeyStyle.Render("[j/k]"),
		helpKeyStyle.Render("[Space]"),
		helpKeyStyle.Render("[Enter/i]"),
		helpKeyStyle.Render("[r]"),
		helpKeyStyle.Render("[Tab/Shift+Tab]"),
	)

	return header + body + footer
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
