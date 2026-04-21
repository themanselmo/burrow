package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/anselmo/burrow/internal/locale"
	"github.com/anselmo/burrow/internal/pet"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type attrTickMsg time.Time
type animTickMsg time.Time

type OpenLogMsg struct{}
type OpenReleaseMsg struct{}
type OpenMissionStartMsg struct{}
type OpenInventoryMsg struct{}
type MissionCompleteCheckMsg struct{}
type QuitMsg struct{ Abrupt bool }

func attrTickCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return attrTickMsg(t)
	})
}

func animTickCmd() tea.Cmd {
	return tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg {
		return animTickMsg(t)
	})
}

type HomeModel struct {
	pet        *pet.Pet
	width      int
	height     int
	anim       Animation
	animFrame  int
	idleFrames int // ticks remaining before next animation plays
}

func NewHomeModel(p *pet.Pet) HomeModel {
	return HomeModel{
		pet:        p,
		anim:       AnimIdle,
		idleFrames: randomIdleWait(),
	}
}

func (m HomeModel) Init() tea.Cmd {
	return tea.Batch(attrTickCmd(), animTickCmd())
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case attrTickMsg:
		m.pet.Tick()
		return m, attrTickCmd()

	case animTickMsg:
		m = m.stepAnim()
		return m, animTickCmd()

	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "f":
			m.pet.Feed()
		case "p":
			m.pet.Pet_()
		case "a":
			m.pet.Play()
		case "l":
			return m, func() tea.Msg { return OpenLogMsg{} }
		case "i":
			return m, func() tea.Msg { return OpenInventoryMsg{} }
		case "m":
			return m, func() tea.Msg { return OpenMissionStartMsg{} }
		case "r":
			return m, func() tea.Msg { return OpenReleaseMsg{} }
		case "q", "ctrl+c":
			return m, func() tea.Msg { return QuitMsg{Abrupt: msg.String() == "ctrl+c"} }
		}
	}
	return m, nil
}

func (m HomeModel) stepAnim() HomeModel {
	// If mid-animation, advance the frame.
	if m.animFrame < len(m.anim)-1 {
		m.animFrame++
		return m
	}

	// Animation finished — return to idle countdown.
	m.anim = AnimIdle
	m.animFrame = 0

	if m.idleFrames > 0 {
		m.idleFrames--
		return m
	}

	// Idle wait over — pick the next animation.
	if m.pet.IsSleepy() {
		m.anim = AnimSleepyNod
	} else {
		m.anim = RandomIdleAnim()
	}
	m.animFrame = 0
	m.idleFrames = randomIdleWait()
	return m
}

func (m HomeModel) currentSprite() []string {
	if m.animFrame < len(m.anim) {
		return m.anim[m.animFrame]
	}
	return SpriteFor(m.pet.Mood(), m.pet.IsSleepy())
}

func (m HomeModel) View() string {
	mood := m.pet.Mood()
	sleepy := m.pet.IsSleepy()
	sprite := m.currentSprite()

	spriteColor := lipgloss.Color("14")
	if sleepy {
		spriteColor = lipgloss.Color("8")
	}
	spriteStyle := lipgloss.NewStyle().Foreground(spriteColor)
	renderedSprite := spriteStyle.Render(strings.Join(sprite, "\n"))

	timeLabel := locale.T("time.awake")
	if sleepy {
		timeLabel = locale.T("time.sleepy")
	}

	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Width(12)
	accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	nameHeader := lipgloss.JoinHorizontal(lipgloss.Bottom,
		nameStyle.Render(m.pet.Name),
		"  ",
		dimStyle.Render(fmt.Sprintf("%s %d  ·  %d/%d %s  ·  %s",
			locale.T("ui.level"),
			m.pet.Level,
			m.pet.XP,
			xpForDisplay(m.pet.Level),
			locale.T("ui.xp"),
			timeLabel,
		)),
	)

	bar := func(label string, val float64) string {
		filled := int(val / 10)
		empty := 10 - filled
		b := accentStyle.Render(strings.Repeat("█", filled)) + dimStyle.Render(strings.Repeat("░", empty))
		return lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render(label), b, dimStyle.Render(fmt.Sprintf(" %d%%", int(val))))
	}

	stats := strings.Join([]string{
		bar(locale.T("ui.hunger"), m.pet.Hunger),
		bar(locale.T("ui.happiness"), m.pet.Happiness),
		bar(locale.T("ui.energy"), m.pet.Energy),
		bar(locale.T("ui.affection"), m.pet.Affection),
	}, "\n")

	actions := dimStyle.Render(strings.Join([]string{
		locale.T("ui.action_feed"),
		locale.T("ui.action_pet"),
		locale.T("ui.action_play"),
		locale.T("ui.action_log"),
		locale.T("ui.action_inventory"),
		locale.T("ui.action_mission"),
		locale.T("ui.action_release"),
		locale.T("ui.action_quit"),
	}, "  "))

	habitatStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 3).
		Width(20).
		Align(lipgloss.Center)

	habitat := habitatStyle.Render(renderedSprite + "\n\n" + dimStyle.Render(locale.T("mood."+mood)))

	statsPanel := lipgloss.NewStyle().
		Padding(0, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			nameHeader,
			"",
			stats,
			"",
			actions,
		))

	main := lipgloss.JoinHorizontal(lipgloss.Top, habitat, statsPanel)
	return lipgloss.NewStyle().Padding(1, 2).Render(main)
}

func xpForDisplay(level int) int {
	v := 100
	for i := 1; i < level; i++ {
		v = int(float64(v) * 1.2)
	}
	return v
}

// randomIdleWait returns a tick count (at 400ms each) to wait between animations,
// giving roughly 3–7 seconds of stillness.
func randomIdleWait() int {
	return 8 + rand.Intn(10)
}
