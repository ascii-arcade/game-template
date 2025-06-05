package menu

import (
	"errors"
	"strings"

	"github.com/ascii-arcade/game-template/colors"
	"github.com/ascii-arcade/game-template/games"
	"github.com/ascii-arcade/game-template/keys"
	"github.com/ascii-arcade/game-template/messages"
	"github.com/ascii-arcade/game-template/screen"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type joinScreen struct {
	model *Model
	style lipgloss.Style
}

func (m *Model) newJoinScreen() *joinScreen {
	m.gameCodeInput.Focus()
	m.gameCodeInput.SetValue("")
	return &joinScreen{
		model: m,
		style: m.style,
	}
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
			if len(s.model.gameCodeInput.Value()) == 7 {
				code := strings.ToUpper(s.model.gameCodeInput.Value())

				game, err := games.GetOpenGame(code)
				if err != nil {
					if !errors.Is(err, games.ErrGameInProgress) && !game.HasPlayer(s.model.player) {
						s.model.setError(err.Error())
						return s.model, nil
					}
				}

				if err := s.model.joinGame(code, false); err != nil {
					s.model.setError(err.Error())
					return s.model, nil
				}

				return s.model, func() tea.Msg { return messages.SwitchToBoardMsg{Game: game} }
			}
		}

		s.model.clearError()
		val := s.model.gameCodeInput.Value()
		if len(val) == 3 && msg.Type == tea.KeyRunes && msg.Runes[0] != '-' {
			val = val + "-"
			s.model.gameCodeInput.SetValue(val)
			s.model.gameCodeInput.CursorEnd()
		}
	}

	s.model.gameCodeInput, cmd = s.model.gameCodeInput.Update(msg)
	return s.model, cmd
}

func (s *joinScreen) View() string {
	errorMessage := ""
	if s.model.errorCode != "" {
		errorMessage = s.model.lang().Get("error", s.model.errorCode)
	}

	var content strings.Builder
	content.WriteString(s.model.lang().Get("menu", "enter_code") + "\n\n")
	content.WriteString(s.model.gameCodeInput.View() + "\n\n")
	content.WriteString(s.style.Foreground(colors.Error).Render(errorMessage))

	return content.String()
}
