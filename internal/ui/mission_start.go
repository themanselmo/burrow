package ui

import (
	"strings"

	"github.com/anselmo/burrow/internal/locale"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MissionCancelMsg struct{}
type MissionStartMsg struct{ DurationMinutes int }

type missionOption struct {
	label   string
	minutes int
}

var missionOptions = []missionOption{
	{label: "25 min  — quick focus sprint", minutes: 25},
	{label: "45 min  — deep work session", minutes: 45},
	{label: "Custom  — set your own duration", minutes: 0},
}

type MissionStartModel struct {
	cursor      int
	customMode  bool
	customInput textinput.Model
	err         string
}

func NewMissionStartModel() MissionStartModel {
	ti := textinput.New()
	ti.Placeholder = "minutes (max 90 for full XP)"
	ti.CharLimit = 4
	ti.Width = 32
	return MissionStartModel{customInput: ti}
}

func (m MissionStartModel) Init() tea.Cmd { return nil }

func (m MissionStartModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.customMode {
		return m.updateCustom(msg)
	}

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(missionOptions)-1 {
				m.cursor++
			}
		case "enter":
			opt := missionOptions[m.cursor]
			if opt.minutes == 0 {
				m.customMode = true
				m.customInput.Focus()
				return m, textinput.Blink
			}
			return m, func() tea.Msg { return MissionStartMsg{DurationMinutes: opt.minutes} }
		case "esc", "q":
			return m, func() tea.Msg { return MissionCancelMsg{} }
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MissionStartModel) updateCustom(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyEnter:
			val := strings.TrimSpace(m.customInput.Value())
			var mins int
			_, err := parseMinutes(val, &mins)
			if err != "" {
				m.err = err
				return m, nil
			}
			return m, func() tea.Msg { return MissionStartMsg{DurationMinutes: mins} }
		case tea.KeyEsc:
			m.customMode = false
			m.customInput.Blur()
			m.err = ""
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.customInput, cmd = m.customInput.Update(msg)
	return m, cmd
}

func parseMinutes(s string, out *int) (struct{}, string) {
	if s == "" {
		return struct{}{}, locale.T("mission.custom_empty")
	}
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return struct{}{}, locale.T("mission.custom_invalid")
		}
		n = n*10 + int(c-'0')
	}
	if n < 1 {
		return struct{}{}, locale.T("mission.custom_too_short")
	}
	*out = n
	return struct{}{}, ""
}

func (m MissionStartModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).MarginBottom(1)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	var sb strings.Builder
	sb.WriteString(headerStyle.Render(locale.T("mission.start_header")) + "\n")
	sb.WriteString(dimStyle.Render(locale.T("mission.start_sub")) + "\n\n")

	if m.customMode {
		sb.WriteString(dimStyle.Render(locale.T("mission.custom_prompt")) + "\n")
		sb.WriteString(m.customInput.View() + "\n")
		if m.err != "" {
			sb.WriteString(errStyle.Render(m.err) + "\n")
		}
		sb.WriteString("\n" + dimStyle.Render(locale.T("mission.custom_back")))
	} else {
		for i, opt := range missionOptions {
			cursor := "  "
			style := normalStyle
			if i == m.cursor {
				cursor = "▶ "
				style = selectedStyle
			}
			sb.WriteString(cursor + style.Render(opt.label) + "\n")
		}
		sb.WriteString("\n" + dimStyle.Render(locale.T("mission.start_hint")))
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 4)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, boxStyle.Render(sb.String()))
}
