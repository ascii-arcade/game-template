package screen

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
	SetModel(any)
	Update(tea.Msg) (any, tea.Cmd)
	View() string
}
