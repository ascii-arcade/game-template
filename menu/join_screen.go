package menu

import (
	"errors"
	"strings"

	"github.com/ascii-arcade/game-template/colors"
	"github.com/ascii-arcade/game-template/games"
	"github.com/ascii-arcade/game-template/keys"
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

func (s *joinScreen) Update(msg tea.Msg) (any, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keys.PreviousScreen.TriggeredBy(msg.String()):
			s.model.screen = s.model.newTitleScreen()
			return s.model, nil
		case keys.Submit.TriggeredBy(msg.String()):
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

				return s.model, func() tea.Msg { return SwitchToBoardMsg{Game: game} }
			}
		}

		s.model.clearError()
		s.model.gameCodeInput, cmd = s.model.gameCodeInput.Update(msg)
		s.model.gameCodeInput.SetValue(strings.ToUpper(s.model.gameCodeInput.Value()))

		if len(s.model.gameCodeInput.Value()) == 3 && !strings.Contains(s.model.gameCodeInput.Value(), "-") {
			s.model.gameCodeInput.SetValue(s.model.gameCodeInput.Value() + "-")
			s.model.gameCodeInput.CursorEnd()
		}
	}
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
