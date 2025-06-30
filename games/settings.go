package games

type Settings struct {
	EndPoints int
}

func NewSettings() Settings {
	return Settings{
		EndPoints: 10,
	}
}
