package games

import (
	"sort"
	"sync"

	generaterandom "github.com/ascii-arcade/wish-template/generate_random"
)

var Games = make(map[string]*Game)

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
	Games[game.Code] = game

	return game
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
		select {
		case p.UpdateChan <- struct{}{}:
		default:
		}
	}
}

func (s *Game) Count(pName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if player, exists := s.Players[pName]; exists {
		player.Count++
	}
}

func Get(code string) (*Game, bool) {
	game, exists := Games[code]
	return game, exists
}

func (s *Game) AddPlayer(updateChan chan struct{}) *Player {
	s.mu.Lock()
	defer s.mu.Unlock()

	player := &Player{
		Name:       generaterandom.Name(),
		Count:      0,
		TurnOrder:  len(s.Players),
		UpdateChan: updateChan,
	}

	s.Players[player.Name] = player
	return player
}

func (s *Game) RemovePlayer(playerName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if player, exists := s.Players[playerName]; exists {
		close(player.UpdateChan)
		delete(s.Players, playerName)

		if len(s.Players) == 0 {
			delete(Games, playerName)
		}
	}
}
