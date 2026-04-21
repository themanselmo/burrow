package ui

import (
	"strings"

	"github.com/anselmo/burrow/internal/locale"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NameSubmittedMsg struct{ Name string }

type NamingModel struct {
	input textinput.Model
	err   string
}

func NewNamingModel() NamingModel {
	ti := textinput.New()
	ti.Placeholder = locale.T("greeting.name_placeholder")
	ti.Focus()
	ti.CharLimit = 24
	ti.Width = 30

	return NamingModel{input: ti}
}

func (m NamingModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m NamingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			name := strings.TrimSpace(m.input.Value())
			if name == "" {
				m.err = locale.T("greeting.name_empty_error")
				return m, nil
			}
			return m, func() tea.Msg { return NameSubmittedMsg{Name: name} }
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m NamingModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 3).
		Width(40)

	sealLines := SpriteFor("content", false)
	seal := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Render(strings.Join(sealLines, "\n"))

	errLine := ""
	if m.err != "" {
		errLine = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.err)
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		seal,
		"",
		titleStyle.Render(locale.T("greeting.welcome")),
		locale.T("greeting.name_prompt"),
		"",
		m.input.View(),
		errLine,
	)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, boxStyle.Render(content))
}
