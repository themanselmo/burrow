package ui

import (
	"strings"
	"time"

	"github.com/anselmo/burrow/internal/locale"
	"github.com/anselmo/burrow/internal/pet"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type abruptTickMsg time.Time

func abruptTickCmd() tea.Cmd {
	return tea.Tick(350*time.Millisecond, func(t time.Time) tea.Msg {
		return abruptTickMsg(t)
	})
}

type AbruptQuitModel struct {
	pet   *pet.Pet
	frame int
}

func NewAbruptQuitModel(p *pet.Pet) AbruptQuitModel {
	return AbruptQuitModel{pet: p}
}

func (m AbruptQuitModel) Init() tea.Cmd {
	return abruptTickCmd()
}

func (m AbruptQuitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(abruptTickMsg); ok {
		m.frame++
		if m.frame >= len(AnimAbruptQuit) {
			return m, tea.Quit
		}
		return m, abruptTickCmd()
	}
	return m, nil
}

func (m AbruptQuitModel) View() string {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	alertStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))

	frame := m.frame
	if frame >= len(AnimAbruptQuit) {
		frame = len(AnimAbruptQuit) - 1
	}
	sprite := AnimAbruptQuit[frame]
	spriteRender := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).
		Render(strings.Join(sprite, "\n"))

	content := lipgloss.JoinVertical(lipgloss.Center,
		spriteRender,
		"",
		alertStyle.Render(locale.T("quit.abrupt_line1")),
		dimStyle.Render(locale.Tf("quit.abrupt_line2", m.pet.Name)),
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("9")).
		Padding(1, 4)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, boxStyle.Render(content))
}
