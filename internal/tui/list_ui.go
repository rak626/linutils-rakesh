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
	// Colors - Nord-inspired / Modern
	fgColor      = lipgloss.Color("#D8DEE9")
	accentColor  = lipgloss.Color("#88C0D0") // Frost blue
	successColor = lipgloss.Color("#A3BE8C") // Aurora green
	warningColor = lipgloss.Color("#EBCB8B") // Aurora yellow
	grayColor    = lipgloss.Color("#4C566A")
	dimColor     = lipgloss.Color("#626262")
	white        = lipgloss.Color("#FFFFFF")

	// Styles
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			MarginBottom(1).
			Padding(0, 2)

	sidebarStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(grayColor).
			Width(40) // Increased width

	mainContentStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(accentColor)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(fgColor)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	cursorItemStyle = lipgloss.NewStyle().
			Background(grayColor).
			Foreground(white).
			Bold(true)

	systemInfoStyle = lipgloss.NewStyle().
			MarginTop(1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(grayColor).
			PaddingTop(1)

	sysKeyStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	sysValStyle = lipgloss.NewStyle().Foreground(fgColor)

	footerStyle = lipgloss.NewStyle().
			MarginTop(1).
			Padding(0, 2).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(grayColor)

	helpLabelStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	helpKeyStyle   = lipgloss.NewStyle().Foreground(warningColor)
	helpTextStyle  = lipgloss.NewStyle().Foreground(dimColor)

	searchInputStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Border(lipgloss.NormalBorder()).
				BorderForeground(grayColor).
				Padding(0, 1)
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
	ScrollOffset int
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
				m.Action = "back"
				m.Finished = true
				return m, tea.Quit
			}
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Filtered) - 1
			}
			m.fixScroll()
		case "down", "j":
			if m.Cursor < len(m.Filtered)-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}
			m.fixScroll()
		case "tab":
			m.ActiveTab = (m.ActiveTab + 1) % len(m.Tabs)
			m.Cursor = 0
			m.ScrollOffset = 0
			m.filterItems()
		case "shift+tab":
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
	m.fixScroll()
}

func (m ListModel) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Initializing..."
	}

	// 1. Header
	title := strings.ToUpper(m.Title)
	header := headerStyle.Render("󱄅 LINUTILS RAKESH  󰁔  " + title)

	// 2. Sidebar (System Info) - More Width
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
		helpLabelStyle.Render("󰄾 CURRENT TASK"),
		lipgloss.NewStyle().Foreground(accentColor).Bold(true).Render(m.Title),
		helpLabelStyle.Render("󰄾 SYSTEM SPECIFICATIONS"),
		sysKeyStyle.Render("OS:  "), sysValStyle.Render(osDisplay),
		sysKeyStyle.Render("DE:  "), sysValStyle.Render(deDisplay),
		sysKeyStyle.Render("CPU: "), sysValStyle.Render(m.SysInfo.CPU),
		sysKeyStyle.Render("RAM: "), sysValStyle.Render(m.SysInfo.RAM),
		sysKeyStyle.Render("DISK:"), sysValStyle.Render(m.SysInfo.Disk),
		sysKeyStyle.Render("GPU: "), sysValStyle.Render(m.SysInfo.GPU),
	)
	
	sidebar := sidebarStyle.
		Height(bodyHeight).
		Render(sysInfoContent)

	// 3. Main Content (Presets)
	boxTitle := "SELECTION"
	if m.Title != "" {
		boxTitle = m.Title
	}
	listContent := helpLabelStyle.Render("󰄾 " + strings.ToUpper(boxTitle)) + "\n\n"
	
	visibleHeight := bodyHeight - 4
	end := m.ScrollOffset + visibleHeight
	if end > len(m.Filtered) {
		end = len(m.Filtered)
	}

	for i := m.ScrollOffset; i < end; i++ {
		idx := m.Filtered[i]
		item := m.Items[idx]
		
		// Selection indicator
		checked := "○"
		if item.Selected {
			checked = lipgloss.NewStyle().Foreground(successColor).Render("●")
		}

		// Cursor indicator and item text
		line := fmt.Sprintf(" %s %s", checked, item.Name)
		
		if m.Cursor == i {
			listContent += cursorItemStyle.Width(m.Width - 55).Render(""+line) + "\n"
		} else {
			if item.Selected {
				listContent += selectedItemStyle.Render(" "+line) + "\n"
			} else {
				listContent += "  " + line + "\n"
			}
		}
	}

	main := mainContentStyle.
		Height(bodyHeight).
		Width(m.Width - 48).
		Render(listContent)

	// Combine Sidebar and Main
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)

	// 4. Footer
	footerContent := ""
	if len(m.Filtered) > 0 && m.Cursor < len(m.Filtered) {
		item := m.Items[m.Filtered[m.Cursor]]
		desc := item.Description
		// Truncate description if too long
		if len(desc) > m.Width-15 {
			desc = desc[:m.Width-18] + "..."
		}
		footerContent += fmt.Sprintf("%s %s\n", helpLabelStyle.Render("󰛨 DESC:"), desc)
	}
	
	commands := fmt.Sprintf("%s Quit  %s Navigate  %s Select  %s Toggle All  %s Install",
		helpKeyStyle.Render("[q]"),
		helpKeyStyle.Render("[j/k]"),
		helpKeyStyle.Render("[Space]"),
		helpKeyStyle.Render("[Ctrl+v]"),
		helpKeyStyle.Render("[Enter]"),
	)
	
	footer := footerStyle.Width(m.Width - 4).Render(footerContent + commands)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
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
