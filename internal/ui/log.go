package ui

import (
	"fmt"
	"strings"

	"github.com/themanselmo/burrow/internal/locale"
	"github.com/themanselmo/burrow/internal/pet"
	"github.com/themanselmo/burrow/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogBackMsg struct{}

type LogModel struct {
	entries    []pet.LogEntry
	currentPet *pet.Pet
}

func NewLogModel(currentPet *pet.Pet) LogModel {
	log, _ := storage.LoadLog()
	entries := []pet.LogEntry{}
	if log != nil {
		entries = log.Entries
	}
	return LogModel{entries: entries, currentPet: currentPet}
}

func (m LogModel) Init() tea.Cmd { return nil }

func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch strings.ToLower(key.String()) {
		case "b", "esc":
			return m, func() tea.Msg { return LogBackMsg{} }
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m LogModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).MarginBottom(1)
	colStyle := lipgloss.NewStyle().Width(16).Foreground(lipgloss.Color("8"))
	nameStyle := lipgloss.NewStyle().Width(16).Foreground(lipgloss.Color("15"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	currentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	var sb strings.Builder
	sb.WriteString(headerStyle.Render(locale.T("log.header")) + "\n\n")

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		colStyle.Render(locale.T("log.col_name")),
		colStyle.Render(locale.T("log.col_species")),
		colStyle.Render(locale.T("log.col_level")),
		colStyle.Render(locale.T("log.col_owned_since")),
		colStyle.Render(locale.T("log.col_released")),
	)
	sb.WriteString(dimStyle.Render(header) + "\n")
	sb.WriteString(dimStyle.Render(strings.Repeat("─", 80)) + "\n")

	if m.currentPet != nil {
		released := currentStyle.Render(locale.T("log.current_label"))
		row := lipgloss.JoinHorizontal(lipgloss.Top,
			nameStyle.Render(m.currentPet.Name),
			colStyle.Render(string(m.currentPet.Species)),
			colStyle.Render(fmt.Sprintf("%d", m.currentPet.Level)),
			colStyle.Render(m.currentPet.OwnedSince.Format("2006-01-02")),
			released,
		)
		sb.WriteString(row + "\n")
	}

	for _, e := range m.entries {
		releasedAt := "—"
		if e.ReleasedAt != nil {
			releasedAt = e.ReleasedAt.Format("2006-01-02")
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top,
			nameStyle.Render(e.Name),
			colStyle.Render(string(e.Species)),
			colStyle.Render(fmt.Sprintf("%d", e.Level)),
			colStyle.Render(e.OwnedSince.Format("2006-01-02")),
			colStyle.Render(releasedAt),
		)
		sb.WriteString(row + "\n")
	}

	if m.currentPet == nil && len(m.entries) == 0 {
		sb.WriteString(dimStyle.Render(locale.T("log.empty")) + "\n")
	}

	sb.WriteString("\n" + dimStyle.Render(locale.T("log.back")))
	return sb.String()
}
