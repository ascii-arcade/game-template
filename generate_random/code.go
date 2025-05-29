package generaterandom

import (
	"crypto/rand"
	"math/big"
)

func Code(existing []string) string {
	exists := make(map[string]struct{}, len(existing))
	for _, v := range existing {
		exists[v] = struct{}{}
	}

	for {
		b := make([]byte, 7)
		letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		for i := range b {
			if i == 3 {
				b[i] = '-'
				continue
			}
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
			if err != nil {
				continue
			}
			b[i] = letters[num.Int64()]
		}
		code := string(b)
		return code
	}
}
