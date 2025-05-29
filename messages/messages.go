package messages

type (
	SwitchToMenu struct{}
	SwitchToGame struct{}
	NewGame      struct{}
	JoinGame     struct{ GameCode string }
	RefreshGame  int
)
