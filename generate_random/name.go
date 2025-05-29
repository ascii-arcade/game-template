package generaterandom

import (
	"fmt"
	"math/rand/v2"
)

var adjectives = []string{"Brave", "Swift", "Clever", "Mighty", "Silent", "Wise"}
var nouns = []string{"Lion", "Eagle", "Wizard", "Ninja", "Knight", "Dragon"}

func Name() string {
	adj := adjectives[rand.IntN(len(adjectives))]
	noun := nouns[rand.IntN(len(nouns))]
	num := rand.IntN(1000)

	return fmt.Sprintf("%s%s%d", adj, noun, num)
}
