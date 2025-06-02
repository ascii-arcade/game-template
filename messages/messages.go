package messages

import (
	"github.com/ascii-arcade/game-template/screen"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	SwitchViewMsg   struct{ Model tea.Model }
	SwitchScreenMsg struct{ Screen screen.Screen }
	NewGame         struct{}
	JoinGame        struct{ GameCode string }
	RefreshBoard    struct{}
)
