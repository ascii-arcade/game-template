package game

import (
	"sort"
	"sync"

	generaterandom "github.com/ascii-arcade/wish-template/internal/generate_random"
	"github.com/ascii-arcade/wish-template/internal/messages"
)

var games = make(map[string]*Game)

type Game struct {
	mu      sync.Mutex
	Players map[string]*Player
	Code    string
}

type Player struct {
	Name      string
	Count     int
	TurnOrder int
	RefreshCh chan messages.RefreshGame
}

func NewGame() *Game {
	g := &Game{
		Players: make(map[string]*Player),
		Code:    generaterandom.Code(),
	}
	games[g.Code] = g
	return g
}

func GetGame(code string) (*Game, bool) {
	g, exists := games[code]
	return g, exists
}

func (s *Game) LockState() {
	s.mu.Lock()
}

func (s *Game) UnlockState() {
	s.mu.Unlock()
}

func (s *Game) OrderedPlayers() []*Player {
	var players []*Player
	for _, p := range s.Players {
		players = append(players, p)
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].TurnOrder < players[j].TurnOrder
	})

	return players
}

func (s *Game) Refresh() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.Players {
		if p.RefreshCh != nil {
			select {
			case p.RefreshCh <- messages.RefreshGame{}:
			default:
			}
		}
	}
}

func (s *Game) RemovePlayer(player *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, exists := s.Players[player.Name]; exists {
		close(p.RefreshCh) // Close the channel to signal no more updates
		delete(s.Players, player.Name)
	}
}

func (s *Game) AddPlayer() *Player {
	s.mu.Lock()
	defer s.mu.Unlock()

	player := &Player{
		Name:      generaterandom.Name(),
		Count:     0,
		TurnOrder: len(s.Players) + 1,
		RefreshCh: make(chan messages.RefreshGame, 1),
	}
	s.Players[player.Name] = player
	return player
}
