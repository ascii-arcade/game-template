package menu

import (
	"strings"

	"github.com/ascii-arcade/wish-template/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type joinScreen struct {
	model *Model
	style lipgloss.Style
}

func (m *Model) newJoinScreen() *joinScreen {
	return &joinScreen{
		model: m,
		style: m.style,
	}
}

func (s *joinScreen) setModel(model *Model) {
	s.model = model
}

func (s *joinScreen) update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		s.model.screen = s.model.newOptionScreen()
	case "enter":
		if len(s.model.gameCodeInput.Value()) == 7 {
			code := strings.ToUpper(s.model.gameCodeInput.Value())
			return s.model, func() tea.Msg { return messages.JoinGame{GameCode: code} }
		}
	default:
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

func (s *joinScreen) view() string {
	style := s.style.Width(s.model.Width).Height(s.model.Height)
	paneStyle := s.model.style.Width(s.model.Width).Height(s.model.Height / 2)

	content := "Enter the game code to join:\n\n" + s.model.gameCodeInput.View()

	panes := lipgloss.JoinVertical(
		lipgloss.Center,
		paneStyle.MarginBottom(2).Align(lipgloss.Center, lipgloss.Bottom).Render(logo),
		paneStyle.Align(lipgloss.Center, lipgloss.Top).Render(content),
	)

	return style.Render(panes)
}
