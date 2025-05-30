package board

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type lobbyScreen struct {
	model *Model
}

func (s *lobbyScreen) setModel(model *Model) {
	s.model = model
}

func (s *lobbyScreen) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s":
		if s.model.Player.IsHost() {
			s.model.screen = &tableScreen{}
			s.model.gameState().Begin()
		}
	}

	return s.model, nil
}

func (s *lobbyScreen) View() string {
	playerList := ""
	for _, p := range s.model.gameState().OrderedPlayers() {
		playerList += p.Name
		if p.Name == s.model.Player.Name {
			playerList += " (you)"
		}
		if p.IsHost() {
			playerList += " (host)"
		}
		playerList += "\n"
	}

	waitingMessage := "Waiting for host to start the game..."
	if s.model.Player.IsHost() {
		waitingMessage = "You are the host. Press 's' to start the game."
	}

	return s.model.renderer.NewStyle().Render(fmt.Sprintf("You are %s", s.model.Player.Name)) +
		"\n\n'" + s.model.Game.Code + "'" +
		"\n\n" + s.model.renderer.NewStyle().Render(playerList) +
		"\n\n" + s.model.renderer.NewStyle().Render(waitingMessage) +
		"\n\n" + s.model.renderer.NewStyle().Render("Press 'ctrl+c' to quit")
}
