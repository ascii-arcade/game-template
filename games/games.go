package games

import (
	"errors"
	"sort"
	"sync"

	"slices"

	"github.com/ascii-arcade/game-template/generaterandom"
)

var (
	ErrGameInProgress = errors.New("game_already_in_progress")
	ErrGameNotFound   = errors.New("game_not_found")
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

func GetOpenGame(code string) (*Game, error) {
	game, exists := games[code]
	if !exists {
		return nil, ErrGameNotFound
	}
	if game.inProgress {
		return nil, ErrGameInProgress
	}

	return game, nil
}

func (s *Game) InProgress() bool {
	return s.inProgress
}

func (s *Game) OrderedPlayers() []*Player {
	var players []*Player
	players = append(players, s.players...)
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

func (s *Game) withLock(fn func() error) error {
	s.mu.Lock()
	defer func() {
		s.refresh()
		s.mu.Unlock()
	}()
	return fn()
}

func (s *Game) AddPlayer(player *Player, isHost bool) error {
	return s.withLock(func() error {
		if s.inProgress {
			return ErrGameInProgress
		}

		maxTurnOrder := 0
		for _, p := range s.players {
			if p.TurnOrder > maxTurnOrder {
				maxTurnOrder = p.TurnOrder
			}
		}

		player.SetTurnOrder(maxTurnOrder + 1)
		if isHost {
			player.MakeHost()
		}

		go func() {
			<-player.ctx.Done()
			_ = s.RemovePlayer(player.Name)
		}()

		s.players = append(s.players, player)
		return nil
	})
}

func (s *Game) RemovePlayer(playerName string) error {
	return s.withLock(func() error {
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
		return nil
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
