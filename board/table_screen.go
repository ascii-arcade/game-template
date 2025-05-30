package board

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type tableScreen struct {
}

func (l tableScreen) Update(m *Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "a":
		m.gameState().Count(m.Player.Name)
	case "q":
		m.screen = lobbyScreen{}
	}

	return m, nil
}

func (l tableScreen) View(m *Model) string {
	counts := ""
	for _, p := range m.gameState().OrderedPlayers() {
		counts += fmt.Sprintf("%s: %d\n", p.Name, p.Count)
	}

	return m.renderer.NewStyle().Render(fmt.Sprintf("You are %s", m.Player.Name)) +
		"\n\n" + counts +
		"\n\n'" + m.Game.Code + "'" +
		"\n\n" + m.renderer.NewStyle().Render("Press 'q' to quit")
}
