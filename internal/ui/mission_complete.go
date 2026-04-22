package ui

import (
	"strings"

	"github.com/themanselmo/burrow/internal/locale"
	"github.com/themanselmo/burrow/internal/mission"
	"github.com/themanselmo/burrow/internal/pet"
	// mission type used for sleep check
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MissionDismissMsg struct{ Result mission.Result }

type MissionCompleteModel struct {
	result mission.Result
	pet    *pet.Pet
}

func NewMissionCompleteModel(r mission.Result, p *pet.Pet) MissionCompleteModel {
	return MissionCompleteModel{result: r, pet: p}
}

func (m MissionCompleteModel) Init() tea.Cmd { return nil }

func (m MissionCompleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter", " ":
			r := m.result
			return m, func() tea.Msg { return MissionDismissMsg{Result: r} }
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MissionCompleteModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	xpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true)

	sprite := SpriteFor("happy", false)
	spriteRender := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).
		Render(strings.Join(sprite, "\n"))

	coinStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)

	isSleep := m.result.Type == mission.TypeSleep

	header := locale.Tf("mission.complete_header", m.pet.Name)
	sub := locale.T("mission.complete_sub")
	if isSleep {
		header = locale.Tf("mission.sleep_complete_header", m.pet.Name)
		sub = locale.T("mission.sleep_complete_sub")
	}

	var sb strings.Builder
	sb.WriteString(spriteRender + "\n\n")
	sb.WriteString(headerStyle.Render(header) + "\n")
	sb.WriteString(dimStyle.Render(sub) + "\n\n")

	if isSleep {
		sb.WriteString(dimStyle.Render("  Energy and mood fully restored.") + "\n")
	} else {
		sb.WriteString(nameStyle.Render(locale.T("mission.rewards_header")) + "\n")
		sb.WriteString(dimStyle.Render("  ") + xpStyle.Render(locale.Tf("mission.xp_gained", m.result.XP)) + "\n")
		sb.WriteString(dimStyle.Render("  ") + coinStyle.Render(locale.Tf("mission.coins_gained", m.result.Coins)) + "\n")
		if m.result.Item != nil {
			sb.WriteString(dimStyle.Render("  ") + itemStyle.Render(locale.Tf("mission.item_found", m.result.Item.Name)) + "\n")
			sb.WriteString(dimStyle.Render("    "+locale.T("mission.item_rare")) + "\n")
		} else {
			sb.WriteString(dimStyle.Render("  "+locale.T("mission.no_item")) + "\n")
		}
	}

	sb.WriteString("\n" + dimStyle.Render(locale.T("mission.complete_hint")))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("10")).
		Padding(1, 4)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, boxStyle.Render(sb.String()))
}
