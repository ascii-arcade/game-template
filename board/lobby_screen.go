package board

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type lobbyScreen struct {
}

func (l lobbyScreen) Update(m *Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s":
		if m.Player.IsHost() {
			m.screen = tableScreen{}
			m.gameState().Begin()
		}
	}

	return m, nil
}

func (l lobbyScreen) View(m *Model) string {
	playerList := ""
	for _, p := range m.gameState().OrderedPlayers() {
		playerList += p.Name
		if p.Name == m.Player.Name {
			playerList += " (you)"
		}
		if p.IsHost() {
			playerList += " (host)"
		}
		playerList += "\n"
	}

	waitingMessage := "Waiting for host to start the game..."
	if m.Player.IsHost() {
		waitingMessage = "You are the host. Press 's' to start the game."
	}

	return m.renderer.NewStyle().Render(fmt.Sprintf("You are %s", m.Player.Name)) +
		"\n\n'" + m.Game.Code + "'" +
		"\n\n" + m.renderer.NewStyle().Render(playerList) +
		"\n\n" + m.renderer.NewStyle().Render(waitingMessage) +
		"\n\n" + m.renderer.NewStyle().Render("Press 'ctrl+c' to quit")
}
