package ui

import (
	"fmt"
	"strings"

	"github.com/themanselmo/burrow/internal/locale"
	"github.com/themanselmo/burrow/internal/pet"
	"github.com/themanselmo/burrow/internal/shop"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ShopBackMsg struct{}
type ShopPurchasedMsg struct{ EnergyRestore float64 }

type ShopModel struct {
	pet    *pet.Pet
	cursor int
	notice string
}

func NewShopModel(p *pet.Pet) ShopModel {
	return ShopModel{pet: p}
}

func (m ShopModel) Init() tea.Cmd { return nil }

func (m ShopModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			m.notice = ""
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			m.notice = ""
			if m.cursor < len(shop.FoodItems)-1 {
				m.cursor++
			}
		case "enter":
			item := shop.FoodItems[m.cursor]
			if !m.pet.SpendCoins(item.Cost) {
				m.notice = locale.T("shop.cant_afford")
				return m, nil
			}
			restore := item.EnergyRestore
			m.notice = locale.Tf("shop.bought", item.Name)
			return m, func() tea.Msg { return ShopPurchasedMsg{EnergyRestore: restore} }
		case "b", "esc":
			return m, func() tea.Msg { return ShopBackMsg{} }
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ShopModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	coinStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	noticeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	var sb strings.Builder
	sb.WriteString(headerStyle.Render(locale.T("shop.header")) + "\n")
	sb.WriteString(dimStyle.Render(locale.T("shop.sub")) + "\n\n")
	sb.WriteString(coinStyle.Render(locale.Tf("shop.coins_label", m.pet.Coins)) + "\n\n")

	for i, item := range shop.FoodItems {
		cursor := "  "
		style := normalStyle
		if i == m.cursor {
			cursor = "▶ "
			style = selectedStyle
		}
		costStr := coinStyle.Render(fmt.Sprintf("%dc", item.Cost))
		line := fmt.Sprintf("%-16s %-28s %s",
			item.Name,
			dimStyle.Render(fmt.Sprintf("+%d energy  %s", int(item.EnergyRestore), item.Description)),
			costStr,
		)
		sb.WriteString(cursor + style.Render(item.Name) +
			dimStyle.Render(fmt.Sprintf("   +%d energy  —  %s   ", int(item.EnergyRestore), item.Description)) +
			coinStyle.Render(fmt.Sprintf("%dc", item.Cost)) + "\n")
		_ = line
	}

	sb.WriteString("\n")
	if m.notice != "" {
		if m.notice == locale.T("shop.cant_afford") {
			sb.WriteString(errStyle.Render(m.notice) + "\n")
		} else {
			sb.WriteString(noticeStyle.Render(m.notice) + "\n")
		}
	}

	sb.WriteString("\n" + dimStyle.Render(locale.T("shop.hint")))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 3).
		Width(56)

	return lipgloss.Place(80, 24, lipgloss.Center, lipgloss.Center, boxStyle.Render(sb.String()))
}
