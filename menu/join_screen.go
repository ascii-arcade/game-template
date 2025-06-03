package menu

import (
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
				if _, err := games.GetOpenGame(code); err != nil {
					s.model.setError(s.model.lang().Get("error", err.Error()))
					return s.model, nil
				}
				return s.model, func() tea.Msg { return messages.JoinGame{GameCode: code} }
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
