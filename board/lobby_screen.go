package board

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type lobbyScreen struct {
}

func (l lobbyScreen) Update(m *Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s":
		m.screen = tableScreen{}
		m.gameState().Begin()
	}

	return m, nil
}

func (l lobbyScreen) View(m *Model) string {
	playerNames := make([]string, 0)
	for _, p := range m.gameState().OrderedPlayers() {
		playerNames = append(playerNames, p.Name)
	}

	return m.renderer.NewStyle().Render(fmt.Sprintf("You are %s", m.Player.Name)) +
		"\n\n'" + m.Game.Code + "'" +
		"\n\n" + strings.Join(playerNames, "\n") +
		"\n\n" + m.renderer.NewStyle().Render("Press 's' to start the game") +
		"\n\n" + m.renderer.NewStyle().Render("Press 'ctrl+c' to quit")
}
