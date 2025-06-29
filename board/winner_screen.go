package board

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type winnerScreen struct {
	model *Model
	style lipgloss.Style
}

func (m *Model) newWinnerScreen() *winnerScreen {
	return &winnerScreen{
		model: m,
		style: m.style,
	}
}

func (s *winnerScreen) Update(msg tea.Msg) (any, tea.Cmd) {
	return s.model, nil
}

func (s *winnerScreen) View() string {
	winner := s.model.Game.GetWinner()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		s.model.style.Bold(true).Render(s.model.lang().Get("board", "game_over")),
		s.model.style.Bold(true).Render(fmt.Sprintf(s.model.lang().Get("board", "winner"), winner.Name)),
	)

	return s.style.Width(s.model.width).Height(s.model.height).Render(
		lipgloss.Place(
			s.model.width,
			s.model.height,
			lipgloss.Center,
			lipgloss.Center,
			s.style.
				Padding(2, 2).
				BorderStyle(lipgloss.NormalBorder()).
				Render(content),
		),
	)
}
