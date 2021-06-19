package decks

import (
	"crypto/rand"
	"math/big"
)

// SecureRandomInt return cryptographically secure int between {0, maxInt}
// Used for shuffling Decks
// Bit overkill, consider removing if seeding math/rand in main()
func SecureRandomInt(maxInt int64) (int64, error) {
	jBigInt, err := rand.Int(rand.Reader, big.NewInt(maxInt))
	if err != nil {
		return 0, err
	}
	return jBigInt.Int64(), nil
}

// OpenDeck for json of a open deck
type OpenDeck struct {
	DeckId     string `json:"deck_id"`
	Shuffled   bool   `json:"shuffled"`
	Remaining  int    `json:"remaining"` // len(Cards) number of cards able to be drawn
	Cards      []Card `json:"cards"`     // cards able to be drawn
}

// ClosedDeck for json of a closed deck
type ClosedDeck struct {
	DeckId     string `json:"deck_id"`
	Shuffled   bool   `json:"shuffled"`
	Remaining  int    `json:"remaining"` // len(Cards) number of cards able to be drawn
}
