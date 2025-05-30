package board

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type tableScreen struct {
	model *Model
}

func (s *tableScreen) setModel(model *Model) {
	s.model = model
}

func (s *tableScreen) update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "a":
		s.model.Game.Count(s.model.Player.Name)
	}

	return s.model, nil
}

func (s *tableScreen) view() string {
	counts := ""
	for _, p := range s.model.Game.OrderedPlayers() {
		counts += fmt.Sprintf("%s: %d\n", p.Name, p.Count)
	}

	return s.model.renderer.NewStyle().Render(fmt.Sprintf("You are %s", s.model.Player.Name)) +
		"\n\n" + counts +
		"\n\n" + s.model.renderer.NewStyle().Render("Press 'ctrl+c' to quit")
}
