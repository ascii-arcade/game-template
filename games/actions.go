package games

import "errors"

func (s *Game) Count(pName string) error {
	err := s.withLock(func() error {
		player, exists := s.getPlayer(pName)
		if !exists {
			return errors.New("player not found")
		}
		player.incrementCount()
		return nil
	})

	return err
}
