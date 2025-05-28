package generate_random

import (
	"crypto/rand"
	"encoding/hex"
)

func Code(existing []string) string {
	exists := make(map[string]struct{}, len(existing))
	for _, v := range existing {
		exists[v] = struct{}{}
	}

	for {
		b := make([]byte, 6)
		_, err := rand.Read(b)
		if err != nil {
			continue
		}
		code := hex.EncodeToString(b)[:6]
		if _, found := exists[code]; !found {
			return code
		}
	}
}
