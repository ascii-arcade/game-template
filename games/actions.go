package games

import "errors"

func (s *Game) Count(pName string) error {
	return s.withLock(func() error {
		player, exists := s.getPlayer(pName)
		if !exists {
			return errors.New("Player not found")
		}
		player.incrementCount()
		return nil
	})
}
