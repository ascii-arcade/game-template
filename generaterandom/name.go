package generaterandom

import (
	"fmt"
	"math/rand/v2"
)

var adjectives = []string{
	"Brave", "Swift", "Clever", "Mighty", "Silent", "Wise",
	"Fierce", "Gentle", "Loyal", "Bold", "Nimble", "Sly",
	"Valiant", "Daring", "Calm", "Stealthy", "Fearless", "Radiant",
	"Sturdy", "Cheerful", "Gallant", "Serene", "Vigilant", "Eager",
	"Majestic", "Diligent", "Resourceful", "Charming", "Witty", "Patient",
	"Dynamic", "Inventive", "Resilient", "Ambitious", "Courageous", "Gracious",
	"Tenacious", "Optimistic", "Curious", "Adaptable", "Persistent", "Humble",
	"Trusty", "Playful", "Energetic", "Sincere", "Polite", "Generous",
}

var nouns = []string{
	"Lion", "Eagle", "Wizard", "Ninja", "Knight", "Dragon",
	"Panther", "Falcon", "Samurai", "Ranger", "Wolf", "Tiger",
	"Phoenix", "Bear", "Shark", "Viking", "Pirate", "Giant",
	"Crusader", "Sage", "Hawk", "Otter", "Fox", "Griffin",
	"Jaguar", "Cobra", "Monk", "Sentinel", "Warden", "Druid",
	"Rogue", "Paladin", "Bard", "Seeker", "Hermit", "Oracle",
	"Minotaur", "Hydra", "Pegasus", "Golem", "Sprite", "Imp",
	"Raven", "Stag", "Mammoth", "Coyote", "Badger", "Mantis",
}

func Name() string {
	adj := adjectives[rand.IntN(len(adjectives))]
	noun := nouns[rand.IntN(len(nouns))]

	return fmt.Sprintf("%s %s", adj, noun)
}
