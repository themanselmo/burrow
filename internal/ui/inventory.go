package ui

import (
	"strings"

	"github.com/themanselmo/burrow/internal/locale"
	"github.com/themanselmo/burrow/internal/mission"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InventoryBackMsg struct{}

type InventoryModel struct {
	items []mission.Item
}

func NewInventoryModel(items []mission.Item) InventoryModel {
	return InventoryModel{items: items}
}

func (m InventoryModel) Init() tea.Cmd { return nil }

func (m InventoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch strings.ToLower(key.String()) {
		case "b", "esc":
			return m, func() tea.Msg { return InventoryBackMsg{} }
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m InventoryModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).MarginBottom(1)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	sourceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Width(20)

	var sb strings.Builder
	sb.WriteString(headerStyle.Render(locale.T("inventory.header")) + "\n\n")

	if len(m.items) == 0 {
		sb.WriteString(dimStyle.Render(locale.T("inventory.empty")) + "\n")
	} else {
		sb.WriteString(dimStyle.Render(locale.T("inventory.col_item")+"          "+locale.T("inventory.col_source")) + "\n")
		sb.WriteString(dimStyle.Render(strings.Repeat("─", 40)) + "\n")

		counts := map[string]int{}
		order := []string{}
		for _, item := range m.items {
			if counts[item.Name] == 0 {
				order = append(order, item.Name)
			}
			counts[item.Name]++
		}

		themeNames := map[string]string{
			"forest_path":   "Forest Path",
			"mountain_pass": "Mountain Pass",
			"ancient_ruins": "Ancient Ruins",
			"frozen_peak":   "Frozen Peak",
			"ocean_shore":   "Ocean Shore",
		}

		itemTheme := map[string]string{}
		for _, item := range m.items {
			itemTheme[item.Name] = item.Theme
		}

		for _, name := range order {
			count := counts[name]
			themeName := themeNames[itemTheme[name]]
			label := name
			if count > 1 {
				label = name + dimStyle.Render(" ×"+itoa(count))
			}
			row := lipgloss.JoinHorizontal(lipgloss.Top,
				itemStyle.Width(22).Render("✦ "+label),
				sourceStyle.Render(themeName),
			)
			sb.WriteString(row + "\n")
		}
	}

	sb.WriteString("\n" + dimStyle.Render(locale.T("inventory.back")))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 3).
		Width(50)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, boxStyle.Render(sb.String()))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
