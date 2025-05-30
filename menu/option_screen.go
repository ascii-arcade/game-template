package menu

import (
	"github.com/ascii-arcade/wish-template/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type optionScreen struct {
	model *Model
}

func (s *optionScreen) setModel(model *Model) {
	s.model = model
}

func (s *optionScreen) update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n":
		return s.model, func() tea.Msg { return messages.NewGame{} }
	case "j":
		s.model.screen = &joinScreen{}
		s.model.gameCodeInput.Focus()
		s.model.gameCodeInput.SetValue("")
	}

	return s.model, nil
}

func (s *optionScreen) view() string {
	style := s.model.renderer.NewStyle().Width(s.model.Width).Height(s.model.Height)
	paneStyle := s.model.renderer.NewStyle().Width(s.model.Width).Height(s.model.Height / 2)

	content := "Welcome to the Game!\n\n"
	content += "Press 'n' to create a new game.\n"
	content += "Press 'j' to join an existing game.\n"

	panes := lipgloss.JoinVertical(
		lipgloss.Center,
		paneStyle.MarginBottom(2).Align(lipgloss.Center, lipgloss.Bottom).Render(logo),
		paneStyle.Align(lipgloss.Center, lipgloss.Top).Render(content),
	)

	return style.Render(panes)
}
