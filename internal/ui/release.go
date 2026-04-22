package ui

import (
	"strings"
	"time"

	"github.com/themanselmo/burrow/internal/locale"
	"github.com/themanselmo/burrow/internal/pet"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ReleaseCancelMsg struct{}
type ReleaseConfirmedMsg struct{}
type FarewellDoneMsg struct{}

// --- Confirmation screen ---

type ReleaseModel struct {
	pet *pet.Pet
}

func NewReleaseModel(p *pet.Pet) ReleaseModel {
	return ReleaseModel{pet: p}
}

func (m ReleaseModel) Init() tea.Cmd { return nil }

func (m ReleaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch strings.ToLower(key.String()) {
		case "y":
			return m, func() tea.Msg { return ReleaseConfirmedMsg{} }
		case "n", "esc":
			return m, func() tea.Msg { return ReleaseCancelMsg{} }
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ReleaseModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	bodyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Width(44)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	yesStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	noStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	sprite := SpriteFor("sad", false)
	spriteRender := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).
		Render(strings.Join(sprite, "\n"))

	content := lipgloss.JoinVertical(lipgloss.Center,
		spriteRender,
		"",
		headerStyle.Render(locale.Tf("release.confirm_header", m.pet.Name)),
		"",
		bodyStyle.Align(lipgloss.Center).Render(locale.Tf("release.confirm_body", m.pet.Name)),
		"",
		lipgloss.JoinHorizontal(lipgloss.Top,
			yesStyle.Render(locale.T("release.confirm_yes")),
			dimStyle.Render("   "),
			noStyle.Render(locale.T("release.confirm_no")),
		),
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("11")).
		Padding(1, 4)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, boxStyle.Render(content))
}

// --- Farewell animation screen ---

type farewellTickMsg time.Time

func farewellTickCmd() tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
		return farewellTickMsg(t)
	})
}

type FarewellModel struct {
	pet   *pet.Pet
	frame int
}

func NewFarewellModel(p *pet.Pet) FarewellModel {
	return FarewellModel{pet: p}
}

func (m FarewellModel) Init() tea.Cmd {
	return farewellTickCmd()
}

func (m FarewellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(farewellTickMsg); ok {
		m.frame++
		if m.frame >= len(AnimFarewell) {
			return m, func() tea.Msg { return FarewellDoneMsg{} }
		}
		return m, farewellTickCmd()
	}
	return m, nil
}

func (m FarewellModel) View() string {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))

	frame := m.frame
	if frame >= len(AnimFarewell) {
		frame = len(AnimFarewell) - 1
	}
	sprite := AnimFarewell[frame]
	spriteRender := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).
		Render(strings.Join(sprite, "\n"))

	content := lipgloss.JoinVertical(lipgloss.Center,
		spriteRender,
		"",
		nameStyle.Render(locale.Tf("release.farewell_line1", m.pet.Name)),
		dimStyle.Render(locale.T("release.farewell_line2")),
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 4)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, boxStyle.Render(content))
}
