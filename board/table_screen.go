package board

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tableScreen struct {
	model *Model
	style lipgloss.Style
}

func (m *Model) newTableScreen() *tableScreen {
	return &tableScreen{
		model: m,
		style: m.style,
	}
}

func (s *tableScreen) setModel(model *Model) {
	s.model = model
}

func (s *tableScreen) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			s.model.Game.Count(s.model.Player.Name)
		case "ctrl+c":
			s.model.Game.RemovePlayer(s.model.Player.Name)
			return s.model, tea.Quit
		}
	}

	return s.model, nil
}

func (s *tableScreen) view() string {
	counts := ""
	for _, p := range s.model.Game.OrderedPlayers() {
		counts += fmt.Sprintf("%s: %d\n", p.Name, p.Count)
	}

	return s.style.Render(fmt.Sprintf("You are %s", s.model.Player.Name)) +
		"\n\n" + counts +
		"\n\n" + s.style.Render("Press 'ctrl+c' to quit")
}
