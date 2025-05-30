package board

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type lobbyScreen struct {
}

func (l lobbyScreen) Update(m *Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	game := m.gameState()

	switch msg.String() {
	case "a":
		game.Count(m.Player.Name)
	case "q", "ctrl+c":
		game.RemovePlayer(m.Player.Name)
		return m, tea.Quit
	}

	return m, nil
}

func (l lobbyScreen) View(m *Model) string {
	counts := ""
	for _, p := range m.gameState().OrderedPlayers() {
		counts += fmt.Sprintf("%s: %d\n", p.Name, p.Count)
	}

	return m.renderer.NewStyle().Render(fmt.Sprintf("You are %s", m.Player.Name)) +
		"\n\n" + counts +
		"\n\n'" + m.Game.Code + "'" +
		"\n\n" + m.renderer.NewStyle().Render("Press 'q' to quit")
}
