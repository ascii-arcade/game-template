package games

import (
	"sort"
	"sync"

	"slices"

	"github.com/ascii-arcade/wish-template/generaterandom"
)

var games = make(map[string]*Game)

type Game struct {
	Code string

	inProgress bool
	mu         sync.Mutex
	players    []*Player
}

func New() *Game {
	game := &Game{
		Code:    generaterandom.Code(),
		players: make([]*Player, 0),
	}
	games[game.Code] = game

	return game
}

func Get(code string) (*Game, bool) {
	game, exists := games[code]
	return game, exists
}

func (s *Game) InProgress() bool {
	return s.inProgress
}

func (s *Game) Begin() {
	s.withLock(func() {
		s.inProgress = true
	})
}

func (s *Game) OrderedPlayers() []*Player {
	var players []*Player
	for _, p := range s.players {
		players = append(players, p)
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].TurnOrder < players[j].TurnOrder
	})

	return players
}

func (s *Game) refresh() {
	for _, p := range s.players {
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

func (s *Game) AddPlayer(isHost bool) *Player {
	var player *Player
	s.withLock(func() {
		maxTurnOrder := 0
		for _, p := range s.players {
			if p.TurnOrder > maxTurnOrder {
				maxTurnOrder = p.TurnOrder
			}
		}
		player = newPlayer(maxTurnOrder, isHost)
		s.players = append(s.players, player)
	})

	return player
}

func (s *Game) RemovePlayer(playerName string) {
	s.withLock(func() {
		if player, exists := s.getPlayer(playerName); exists {
			close(player.UpdateChan)
			for i, p := range s.players {
				if p.Name == playerName {
					s.players = slices.Delete(s.players, i, i+1)
					break
				}
			}

			if len(s.players) == 0 {
				delete(games, playerName)
			}
		}
	})
}

func (s *Game) getPlayer(name string) (*Player, bool) {
	for _, player := range s.players {
		if player.Name == name {
			return player, true
		}
	}

	return nil, false
}
