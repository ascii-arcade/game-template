package board

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type lobbyScreen struct {
	model *Model
}

func (s *lobbyScreen) setModel(model *Model) {
	s.model = model
}

func (s *lobbyScreen) update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s":
		if s.model.Player.IsHost() {
			s.model.screen = &tableScreen{}
			s.model.Game.Begin()
		}
	}

	return s.model, nil
}

func (s *lobbyScreen) view() string {
	style := s.model.renderer.NewStyle().Width(s.model.Width / 3)

	footer := "\nWaiting for host to start the game..."
	if s.model.Player.IsHost() {
		footer = "Press 's' to start the game."
	}
	footer += "\nPress 'ctrl+c' to quit."

	header := s.model.Game.Code
	playerList := s.model.renderer.NewStyle().Render(s.playerList())

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		style.Align(lipgloss.Center).MarginBottom(2).Render(header),
		style.Render(playerList),
		style.Render(footer),
	)

	return s.model.renderer.NewStyle().Width(s.model.Width).Height(s.model.Height).Render(
		lipgloss.Place(
			s.model.Width,
			s.model.Height,
			lipgloss.Center,
			lipgloss.Center,
			s.model.renderer.NewStyle().
				Padding(2, 2).
				BorderStyle(lipgloss.NormalBorder()).
				Render(content),
		),
	)
}

func (s *lobbyScreen) playerList() string {
	playerList := ""
	for _, p := range s.model.Game.OrderedPlayers() {
		playerList += "* " + p.Name
		if p.Name == s.model.Player.Name {
			playerList += " (you)"
		}
		if p.IsHost() {
			playerList += " (host)"
		}
		playerList += "\n"
	}
	return playerList
}
