package keys

import "slices"

type Keys []string

func (k Keys) Contains(msg string) bool {
	return slices.Contains(k, msg)
}

func (k Keys) String() string {
	if len(k) == 0 {
		return ""
	}
	return k[0]
}

var (
	MenuJoinGame     = Keys{"j"}
	MenuStartNewGame = Keys{"n"}

	PreviousScreen = Keys{"esc"}
	Submit         = Keys{"enter"}

	ExitApplication    = Keys{"ctrl+c"}
	GameIncrementPoint = Keys{"a"}
	LobbyStartGame     = Keys{"s"}
)
