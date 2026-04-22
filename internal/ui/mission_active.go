package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/themanselmo/burrow/internal/locale"
	"github.com/themanselmo/burrow/internal/mission"
	"github.com/themanselmo/burrow/internal/pet"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type missionProgressTickMsg time.Time
type missionAnimTickMsg time.Time

func missionProgressTickCmd() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return missionProgressTickMsg(t)
	})
}

func missionAnimTickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return missionAnimTickMsg(t)
	})
}

type MissionActiveModel struct {
	mission   *mission.Mission
	pet       *pet.Pet
	walkFrame int
}

func NewMissionActiveModel(m *mission.Mission, p *pet.Pet) MissionActiveModel {
	return MissionActiveModel{mission: m, pet: p}
}

func (m MissionActiveModel) Init() tea.Cmd {
	return tea.Batch(missionProgressTickCmd(), missionAnimTickCmd())
}

func (m MissionActiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case missionProgressTickMsg:
		if m.mission.IsComplete() {
			return m, func() tea.Msg { return MissionCompleteCheckMsg{} }
		}
		return m, missionProgressTickCmd()
	case missionAnimTickMsg:
		m.walkFrame = (m.walkFrame + 1) % len(AnimWalk)
		return m, missionAnimTickCmd()
	}
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "ctrl+c" || strings.ToLower(key.String()) == "q" {
			return m, func() tea.Msg { return QuitMsg{Abrupt: true} }
		}
	}
	return m, nil
}

func (m MissionActiveModel) View() string {
	isSleep := m.mission.Type == mission.TypeSleep

	var waypoints []mission.Waypoint
	if isSleep {
		waypoints = mission.DreamWaypoints
	} else {
		theme, ok := mission.Themes[m.mission.Theme]
		if !ok {
			theme = mission.Themes["forest_path"]
		}
		waypoints = theme.Waypoints
	}

	fraction := m.mission.ElapsedFraction()
	numWaypoints := len(waypoints)

	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)

	// Journey path bar
	totalWidth := 52
	segmentWidth := totalWidth / (numWaypoints - 1)
	filledWidth := int(fraction * float64(totalWidth))

	pathBar := ""
	for i := 0; i < totalWidth; i++ {
		if i <= filledWidth {
			pathBar += accentStyle.Render("─")
		} else {
			pathBar += dimStyle.Render("─")
		}
	}

	// Waypoint markers and labels
	markers := ""
	labels := ""
	for i, wp := range waypoints {
		pos := i * segmentWidth
		if i == numWaypoints-1 {
			pos = totalWidth
		}
		isDone := float64(pos)/float64(totalWidth) <= fraction
		if isDone {
			markers += accentStyle.Render("●")
		} else {
			markers += dimStyle.Render("○")
		}
		if i < numWaypoints-1 {
			pad := segmentWidth - 1
			markers += strings.Repeat(" ", pad)
		}
		label := wp.Label
		labelPad := segmentWidth - len(label)
		if i == numWaypoints-1 {
			labels += labelStyle.Render(label)
		} else {
			labels += labelStyle.Render(label) + strings.Repeat(" ", labelPad)
		}
	}

	// Pet position along the path
	petPos := int(fraction * float64(totalWidth))
	petLine := strings.Repeat(" ", petPos) + accentStyle.Render("▲")

	// Current waypoint flavor
	currentWaypointIdx := int(fraction * float64(numWaypoints-1))
	if currentWaypointIdx >= numWaypoints {
		currentWaypointIdx = numWaypoints - 1
	}
	flavorText := locale.Tf("mission.flavor", m.pet.Name, waypoints[currentWaypointIdx].FlavorIn)
	if isSleep {
		flavorText += locale.T("mission.in_their_dreams")
	}

	// Sprite — walking for explore, sleeping for sleep
	var spriteFrames Animation
	if m.mission.Type == mission.TypeSleep {
		spriteFrames = AnimAsleep
	} else {
		spriteFrames = AnimWalk
	}
	frame := m.walkFrame % len(spriteFrames)
	spriteRender := accentStyle.Render(strings.Join(spriteFrames[frame], "\n"))

	minsRemaining := m.mission.MinutesRemaining()
	timeStr := fmt.Sprintf("%d %s", minsRemaining, locale.T("mission.minutes_remaining"))
	if minsRemaining <= 1 {
		timeStr = locale.T("mission.almost_back")
	}

	header := locale.T("mission.sleep_active_header")
	if !isSleep {
		theme, ok := mission.Themes[m.mission.Theme]
		if ok {
			header = theme.Name
		}
	}

	var sb strings.Builder
	sb.WriteString(headerStyle.Render(header) + "\n\n")
	sb.WriteString(spriteRender + "\n\n")
	sb.WriteString(dimStyle.Render(flavorText) + "\n\n")
	sb.WriteString(markers + "\n")
	sb.WriteString(pathBar + "\n")
	sb.WriteString(petLine + "\n")
	sb.WriteString(labels + "\n\n")
	sb.WriteString(activeStyle.Render(timeStr) + "\n\n")
	hint := locale.T("mission.active_hint")
	if m.mission.Type == mission.TypeSleep {
		hint = locale.T("mission.sleeping_hint")
	}
	sb.WriteString(dimStyle.Render(hint))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 3)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, boxStyle.Render(sb.String()))
}
