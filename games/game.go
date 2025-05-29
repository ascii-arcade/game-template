package games

import (
	"sort"
	"sync"

	generaterandom "github.com/ascii-arcade/wish-template/generate_random"
)

var games = make(map[string]*Game)

type Game struct {
	Code string

	mu      sync.Mutex
	Players map[string]*Player
}

func New() *Game {
	game := &Game{
		Code:    generaterandom.Code(),
		Players: make(map[string]*Player),
	}
	games[game.Code] = game

	return game
}

func Get(code string) (*Game, bool) {
	game, exists := games[code]
	return game, exists
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

func (s *Game) refresh() {
	for _, p := range s.Players {
		select {
		case p.UpdateChan <- struct{}{}:
		default:
		}
	}
}

func (s *Game) withLock(fn func()) {
	s.mu.Lock()
	defer func() {
		s.refresh()
		s.mu.Unlock()
	}()
	fn()
}

func (s *Game) AddPlayer(updateChan chan struct{}) *Player {
	var player *Player
	s.withLock(func() {
		maxTurnOrder := 0
		for _, p := range s.Players {
			if p.TurnOrder > maxTurnOrder {
				maxTurnOrder = p.TurnOrder
			}
		}
		player = &Player{
			Name:       generaterandom.Name(),
			Count:      0,
			TurnOrder:  maxTurnOrder + 1,
			UpdateChan: updateChan,
		}

		s.Players[player.Name] = player
	})

	return player
}

func (s *Game) RemovePlayer(playerName string) {
	s.withLock(func() {
		if player, exists := s.Players[playerName]; exists {
			close(player.UpdateChan)
			delete(s.Players, playerName)

			if len(s.Players) == 0 {
				delete(games, playerName)
			}
		}
	})
}
