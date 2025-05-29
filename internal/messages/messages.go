package messages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	SwitchViewMsg struct {
		NewModel tea.Model
	}
	NewGame     struct{}
	JoinGame    struct{ GameCode string }
	RefreshGame struct{}
)
