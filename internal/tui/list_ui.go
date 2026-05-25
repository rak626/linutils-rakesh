package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	categoryStyle     = lipgloss.NewStyle().MarginLeft(2).MarginTop(1).Bold(true).Foreground(lipgloss.Color("#AF87FF"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#13BCED"))
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).MarginTop(1)
)

type ListItem struct {
	Key      string
	Name     string
	Category string
	Selected bool
}

type ListModel struct {
	Title    string
	Items    []ListItem
	Cursor   int
	Action   string // "r" for remove, "i" for install, "" for none
	Finished bool
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.Finished = true
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Items) - 1
			}
		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}
		case " ":
			m.Items[m.Cursor].Selected = !m.Items[m.Cursor].Selected
		case "a", "A", "ctrl+a":
			allSelected := true
			for _, item := range m.Items {
				if !item.Selected {
					allSelected = false
					break
				}
			}
			for i := range m.Items {
				m.Items[i].Selected = !allSelected
			}
		case "r", "R":
			m.Action = "r"
			m.Finished = true
			return m, tea.Quit
		case "i", "I", "enter":
			m.Action = "i"
			m.Finished = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ListModel) View() string {
	s := titleStyle.Render(m.Title) + "\n"

	var lastCategory string
	for i, item := range m.Items {
		if item.Category != "" && item.Category != lastCategory {
			s += categoryStyle.Render(item.Category) + "\n"
			lastCategory = item.Category
		}

		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := "•" // Using a dot instead of [ ] for cleaner look, or stick to [ ]
		if item.Selected {
			checked = "✓"
		}

		line := fmt.Sprintf("%s %s %s", cursor, checked, item.Name)
		if m.Cursor == i {
			s += selectedItemStyle.Render(line) + "\n"
		} else {
			s += itemStyle.Render(line) + "\n"
		}
	}

	s += helpStyle.Render("\n ↑/↓ Navigate • Space Select • Ctrl+A Toggle All • I Install • R Remove • Q Back")
	return s
}

func RunListUI(title string, items []ListItem) (string, []ListItem, error) {
	m := ListModel{
		Title: title,
		Items: items,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", nil, err
	}

	m = finalModel.(ListModel)
	return m.Action, m.Items, nil
}
