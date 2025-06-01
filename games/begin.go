package games

const (
	minimumPlayers = 2
	maximumPlayers = 5
)

func (s *Game) Begin() {
	s.withLock(func() {
		if _, ok := s.IsPlayerCountOk(); !ok {
			return
		}

		s.inProgress = true
	})
}

func (s *Game) IsPlayerCountOk() (string, bool) {
	if len(s.players) > maximumPlayers {
		return "Too many players", false
	}
	if len(s.players) < minimumPlayers {
		return "Not enough players", false
	}
	return "", true
}
