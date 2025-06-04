package games

import (
	"github.com/ascii-arcade/game-template/generaterandom"
	"github.com/ascii-arcade/game-template/language"
)

type Player struct {
	Name      string
	Count     int
	TurnOrder int

	isHost bool

	UpdateChan         chan struct{}
	LanguagePreference *language.LanguagePreference
}

func NewPlayer(langPref *language.LanguagePreference) *Player {
	return &Player{
		Name:               generaterandom.Name(langPref.Lang),
		Count:              0,
		UpdateChan:         make(chan struct{}),
		LanguagePreference: langPref,
	}
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
