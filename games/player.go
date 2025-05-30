package games

import "github.com/ascii-arcade/wish-template/generaterandom"

type Player struct {
	Name      string
	Count     int
	TurnOrder int

	isHost bool

	UpdateChan chan struct{}
}

func newPlayer(maxTurnOrder int, host bool) *Player {
	return &Player{
		Name:       generaterandom.Name(),
		Count:      0,
		TurnOrder:  maxTurnOrder + 1,
		UpdateChan: make(chan struct{}),
		isHost:     host,
	}
}

func (p *Player) IsHost() bool {
	return p.isHost
}

func (p *Player) incrementCount() {
	p.Count++
}
