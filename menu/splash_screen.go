package menu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const splashLogo = `++------------------------------------------------------------------------------++
++------------------------------------------------------------------------------++
||                                                                              ||
||                                                                              ||
||      _    ____   ____ ___ ___        _    ____   ____    _    ____  _____    ||
||     / \  / ___| / ___|_ _|_ _|      / \  |  _ \ / ___|  / \  |  _ \| ____|   ||
||    / _ \ \___ \| |    | | | |_____ / _ \ | |_) | |     / _ \ | | | |  _|     ||
||   / ___ \ ___) | |___ | | | |_____/ ___ \|  _ <| |___ / ___ \| |_| | |___    ||
||  /_/   \_\____/ \____|___|___|   /_/   \_\_| \_\\____/_/   \_\____/|_____|   ||
||                                                                              ||
||                                                                              ||
||                                                                              ||
++------------------------------------------------------------------------------++
++------------------------------------------------------------------------------++`

type splashScreen struct {
	model *Model
}

func (s *splashScreen) setModel(model *Model) {
	s.model = model
}

func (s *splashScreen) update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	return s.model, nil
}

func (s *splashScreen) view() string {
	style := s.model.renderer.NewStyle().
		Width(s.model.Width).
		Height(s.model.Height)

	return style.Render(
		lipgloss.Place(
			s.model.Width,
			s.model.Height,
			lipgloss.Center,
			lipgloss.Center,
			splashLogo,
		),
	)
}
