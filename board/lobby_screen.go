package board

import (
	"github.com/ascii-arcade/wish-template/colors"
	"github.com/ascii-arcade/wish-template/keys"
	"github.com/ascii-arcade/wish-template/screen"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type lobbyScreen struct {
	model *Model
	style lipgloss.Style
}

func (m *Model) newLobbyScreen() *lobbyScreen {
	return &lobbyScreen{
		model: m,
		style: m.style,
	}
}

func (s *lobbyScreen) WithModel(model any) screen.Screen {
	s.model = model.(*Model)
	return s
}

func (s *lobbyScreen) Update(msg tea.Msg) (any, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if keys.LobbyStartGame.TriggeredBy(msg.String()) {
			if s.model.Player.IsHost() {
				s.model.Game.Begin()
			}
		}
	}

	return s.model, nil
}

func (s *lobbyScreen) View() string {
	style := s.style.Width(s.model.Width / 3)

	footer := "\nWaiting for host to start the game..."
	if s.model.Player.IsHost() {
		err := s.model.Game.IsPlayerCountOk()
		footer = "\nPress " + keys.MenuStartNewGame.String(s.style) + " to start the game."
		if err != nil {
			footer = s.style.Foreground(colors.Error).Render(err.Error())
		}
	}
	footer += "\nPress " + keys.ExitApplication.String(s.style) + " to quit."

	header := s.model.Game.Code
	playerList := s.style.Render(s.playerList())

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		style.Align(lipgloss.Center).MarginBottom(2).Render(header),
		style.Render(playerList),
		style.Render(footer),
	)

	return s.style.Width(s.model.Width).Height(s.model.Height).Render(
		lipgloss.Place(
			s.model.Width,
			s.model.Height,
			lipgloss.Center,
			lipgloss.Center,
			s.style.
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
