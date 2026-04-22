package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/themanselmo/burrow/internal/locale"
	"github.com/themanselmo/burrow/internal/pet"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type attrTickMsg time.Time
type animTickMsg time.Time

type OpenLogMsg struct{}
type OpenReleaseMsg struct{}
type OpenMissionStartMsg struct{ Sleep bool }
type OpenInventoryMsg struct{}
type OpenShopMsg struct{}
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
	idleFrames int
	notice     string
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
		m.notice = ""
		return m, attrTickCmd()

	case animTickMsg:
		m = m.stepAnim()
		return m, animTickCmd()

	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "p":
			m.pet.PetAction()
		case "a":
			m.pet.Play()
		case "l":
			return m, func() tea.Msg { return OpenLogMsg{} }
		case "i":
			return m, func() tea.Msg { return OpenInventoryMsg{} }
		case "s":
			return m, func() tea.Msg { return OpenShopMsg{} }
		case "m":
			if !m.pet.CanGoOnMission() {
				m.notice = locale.T("ui.no_energy")
				return m, nil
			}
			return m, func() tea.Msg { return OpenMissionStartMsg{Sleep: false} }
		case "z":
			return m, func() tea.Msg { return OpenMissionStartMsg{Sleep: true} }
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

	coinStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	stats := strings.Join([]string{
		bar(locale.T("ui.energy"), m.pet.Energy),
		bar(locale.T("ui.mood"), m.pet.MoodLevel),
		"",
		coinStyle.Render(fmt.Sprintf("  ◈ %d %s", m.pet.Coins, locale.T("ui.coins"))),
	}, "\n")

	noticeStr := ""
	if m.notice != "" {
		noticeStr = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.notice)
	}

	row := func(keys ...string) string {
		return dimStyle.Render(strings.Join(keys, "  "))
	}
	actions := strings.Join([]string{
		row(locale.T("ui.action_pet"), locale.T("ui.action_play"), locale.T("ui.action_shop")),
		row(locale.T("ui.action_mission"), locale.T("ui.action_sleep")),
		row(locale.T("ui.action_log"), locale.T("ui.action_inventory"), locale.T("ui.action_release"), locale.T("ui.action_quit")),
	}, "\n")

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
			noticeStr,
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
