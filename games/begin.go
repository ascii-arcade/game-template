package games

import "errors"

const (
	minimumPlayers = 2
	maximumPlayers = 5
)

func (s *Game) Begin() error {
	err := s.withLock(func() error {
		if error := s.IsPlayerCountOk(); error != nil {
			return error
		}

		s.inProgress = true
		return nil
	})

	return err
}

func (s *Game) IsPlayerCountOk() error {
	if len(s.players) > maximumPlayers {
		return errors.New("Too many players")
	}
	if len(s.players) < minimumPlayers {
		return errors.New("Not enough players")
	}
	return nil
}
