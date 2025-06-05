package messages

import (
	"github.com/ascii-arcade/game-template/games"
	"github.com/ascii-arcade/game-template/screen"
)

type (
	// SwitchViewMsg    struct{ Model tea.Model }
	SwitchToMenuMsg  struct{}
	SwitchToBoardMsg struct{ Game *games.Game }
	SwitchScreenMsg  struct{ Screen screen.Screen }
	// NewGame          struct{}
	// JoinGame         struct{ GameCode string }
	RefreshBoard struct{}
)
