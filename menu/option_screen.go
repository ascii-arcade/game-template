package menu

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
		if keys.MenuStartNewGame.Contains(msg.String()) {
			return s.model, func() tea.Msg { return messages.NewGame{} }
		}
		if keys.MenuJoinGame.Contains(msg.String()) {
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
	style := s.style.Width(s.model.Width).Height(s.model.Height)
	paneStyle := s.style.Width(s.model.Width).Height(s.model.Height / 2)

	content := "Welcome to the Game!\n\n"
	content += "Press '" + keys.MenuStartNewGame.String() + "' to create a new game.\n"
	content += "Press '" + keys.MenuJoinGame.String() + "' to join an existing game.\n"

	panes := lipgloss.JoinVertical(
		lipgloss.Center,
		paneStyle.MarginBottom(2).Align(lipgloss.Center, lipgloss.Bottom).Render(logo),
		paneStyle.Align(lipgloss.Center, lipgloss.Top).Render(content),
	)

	return style.Render(panes)
}
