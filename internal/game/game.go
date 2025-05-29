package game

import (
	"sort"
	"sync"
)

var Games = make(map[string]*Game)

type Game struct {
	mu      sync.Mutex
	Clients map[chan int]struct{}
	Players map[string]*Player
}

type Player struct {
	Name      string
	Count     int
	TurnOrder int
}

func NewGame() *Game {
	return &Game{
		Clients: make(map[chan int]struct{}),
		Players: make(map[string]*Player),
	}
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
	for ch := range s.Clients {
		select {
		case ch <- 0:
		default:
		}
	}
}

func (s *Game) AddClient(ch chan int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Clients[ch] = struct{}{}
}

func (s *Game) RemoveClient(ch chan int, player string) {
	delete(s.Clients, ch)
	delete(s.Players, player)

	if len(s.Players) == 0 {
		delete(Games, player)
		return
	}
}
