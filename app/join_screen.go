package app

import (
	"strings"

	"github.com/ascii-arcade/wish-template/colors"
	"github.com/ascii-arcade/wish-template/games"
	"github.com/ascii-arcade/wish-template/keys"
	"github.com/ascii-arcade/wish-template/messages"
	"github.com/ascii-arcade/wish-template/screen"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type joinScreen struct {
	model *Model
	style lipgloss.Style

	gameCodeInput textinput.Model
}

func (m *Model) newJoinScreen() *joinScreen {
	ti := textinput.New()
	ti.Width = 9
	ti.CharLimit = 7
	ti.Focus()

	s := &joinScreen{
		model:         m,
		style:         m.style,
		gameCodeInput: ti,
	}

	return s
}

func (s *joinScreen) WithModel(model any) screen.Screen {
	s.model = model.(*Model)
	return s
}

func (s *joinScreen) Update(msg tea.Msg) (any, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if keys.PreviousScreen.TriggeredBy(msg.String()) {
			return s.model, func() tea.Msg {
				return messages.SwitchScreenMsg{
					Screen: s.model.newOptionScreen(),
				}
			}
		}
		if keys.Submit.TriggeredBy(msg.String()) {
			if len(s.gameCodeInput.Value()) == 7 {
				code := strings.ToUpper(s.gameCodeInput.Value())
				if _, err := games.GetOpenGame(code); err != nil {
					s.model.setError(err.Error())
					return s.model, nil
				}
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
		}

		s.model.clearError()
		val := s.gameCodeInput.Value()
		if len(val) == 3 && msg.Type == tea.KeyRunes && msg.Runes[0] != '-' {
			val = val + "-"
			s.gameCodeInput.SetValue(val)
			s.gameCodeInput.CursorEnd()
		}
	}
	s.gameCodeInput, cmd = s.gameCodeInput.Update(msg)
	return s.model, cmd
}

func (s *joinScreen) View() string {
	style := s.style.Width(s.model.width).Height(s.model.height)
	paneStyle := s.style.Width(s.model.width).Height(s.model.height / 2)

	content := "Enter the game code to join:\n\n" + s.gameCodeInput.View()

	panes := lipgloss.JoinVertical(
		lipgloss.Center,
		paneStyle.MarginBottom(2).Align(lipgloss.Center, lipgloss.Bottom).Render(logo),
		paneStyle.Align(lipgloss.Center, lipgloss.Top).Render(content+"\n\n"+s.style.Foreground(colors.Error).Render(s.model.error)),
	)

	return style.Render(panes)
}
