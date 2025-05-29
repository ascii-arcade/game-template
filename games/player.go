package games

type Player struct {
	Name      string
	Count     int
	TurnOrder int

	UpdateChan chan struct{}
}
