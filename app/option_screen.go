package app

import (
	"github.com/ascii-arcade/wish-template/keys"
	"github.com/ascii-arcade/wish-template/messages"
	"github.com/ascii-arcade/wish-template/screen"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type optionScreen struct {
	model *Model
	style lipgloss.Style
}

func (m *Model) newOptionScreen() *optionScreen {
	return &optionScreen{
		model: m,
		style: m.style,
	}
}

func (s *optionScreen) WithModel(model any) screen.Screen {
	s.model = model.(*Model)
	return s
}

func (s *optionScreen) Update(msg tea.Msg) (any, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if keys.MenuStartNewGame.TriggeredBy(msg.String()) {
			code := s.model.newGame()
			if err := s.model.joinGame(code, true); err != nil {
				s.model.setError(err.Error())
				return s.model, nil
			}
			return s.model, tea.Batch(
				func() tea.Msg {
					return messages.SwitchScreenMsg{
						Screen: s.model.newLobbyScreen(),
					}
				},
				func() tea.Msg {
					return messages.RefreshBoard{}
				},
			)
		}
		if keys.MenuJoinGame.TriggeredBy(msg.String()) {
			return s.model, func() tea.Msg {
				return messages.SwitchScreenMsg{
					Screen: s.model.newJoinScreen(),
				}
			}
		}
	}

	return s.model, nil
}

func (s *optionScreen) View() string {
	style := s.style.Width(s.model.width).Height(s.model.height)
	paneStyle := s.style.Width(s.model.width).Height(s.model.height / 2)

	content := "Welcome to the Game!\n\n"
	content += "Press " + keys.MenuStartNewGame.String(s.style) + " to create a new game.\n"
	content += "Press " + keys.MenuJoinGame.String(s.style) + " to join an existing game.\n"

	panes := lipgloss.JoinVertical(
		lipgloss.Center,
		paneStyle.MarginBottom(2).Align(lipgloss.Center, lipgloss.Bottom).Render(logo),
		paneStyle.Align(lipgloss.Center, lipgloss.Top).Render(content),
	)

	return style.Render(panes)
}
