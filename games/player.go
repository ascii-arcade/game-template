package games

import (
	"context"

	"github.com/ascii-arcade/game-template/language"
	"github.com/charmbracelet/ssh"
)

type Player struct {
	Name      string
	Count     int
	TurnOrder int

	isHost      bool
	IsConnected bool

	UpdateChan         chan struct{}
	LanguagePreference *language.LanguagePreference

	Sess ssh.Session
	ctx  context.Context
}

func (p *Player) OnDisconnect(fn func()) {
	go func() {
		<-p.ctx.Done()
		fn()
	}()
}

func (p *Player) SetName(name string) *Player {
	p.Name = name
	return p
}

func (p *Player) SetTurnOrder(order int) *Player {
	p.TurnOrder = order
	return p
}

func (p *Player) MakeHost() *Player {
	p.isHost = true
	return p
}

func (p *Player) IsHost() bool {
	return p.isHost
}

func (p *Player) incrementCount() {
	p.Count++
}
