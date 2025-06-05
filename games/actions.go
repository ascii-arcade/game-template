package games

func (s *Game) Count(player *Player) {
	player.incrementCount()
}
