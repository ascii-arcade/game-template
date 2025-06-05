package games

import (
	"slices"
	"sort"
	"sync"

	"github.com/charmbracelet/ssh"
)

type Game struct {
	Code string

	inProgress bool
	mu         sync.Mutex
	players    []*Player
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
		if _, ok := s.getPlayer(player.Sess); ok {
			return nil
		}

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

		s.players = append(s.players, player)
		return nil
	})
}

func (s *Game) RemovePlayer(player *Player) {
	_ = s.withLock(func() error {
		if player, exists := s.getPlayer(player.Sess); exists {
			close(player.UpdateChan)
			for i, p := range s.players {
				if p.Sess.User() == player.Sess.User() {
					s.players = slices.Delete(s.players, i, i+1)
					break
				}
			}

			if len(s.players) == 0 {
				delete(games, player.Sess.User())
			}
		}
		return nil
	})
}

func (s *Game) getPlayer(sess ssh.Session) (*Player, bool) {
	var player *Player
	_ = s.withLock(func() error {
		for _, p := range s.players {
			if p.Sess.User() == sess.User() {
				player = p
				break
			}
		}
		return nil
	})
	return player, player != nil
}
