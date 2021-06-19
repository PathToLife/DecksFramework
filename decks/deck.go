package decks

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type Card struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type Deck struct {
	lock       sync.Mutex // thread safety in a concurrent API server
	DeckId     string
	Shuffled   bool
	Remaining  int    // len(Cards) number of cards able to be drawn
	Cards      []Card // cards able to be drawn
	DrawnCards []Card // cards that are drawn
}

func (d *Deck) DrawCard(drawCount int) ([]Card, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if drawCount == 0 {
		return []Card{}, nil
	}

	if d.Remaining == 0 || len(d.Cards) == 0 {
		return nil, errors.New("no more cards to draw")
	}

	if drawCount > d.Remaining || drawCount > len(d.Cards) {
		return nil, errors.New("draw count too high")
	}

	cards := make([]Card, drawCount, drawCount)

	for i := 0; i < drawCount; i++ {
		drawn := d.Cards[0]
		cards[i] = drawn
		d.Cards = d.Cards[1:]
	}
	d.Remaining = len(d.Cards)

	return cards, nil
}

// ShuffleDeck shuffle deck.Cards
// Can shuffle a deck multiple times
func (d *Deck) ShuffleDeck() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	// cannot shuffle deck with less than 2 cards, no error
	if len(d.Cards) < 2 {
		return nil
	}

	// Fisher-Yates shuffle algorithm
	for i := len(d.Cards) - 1; i > 0; i-- {
		j, err := SecureRandomInt(int64(i + 1))
		if err != nil {
			return err
		}
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}

	d.Shuffled = true
	return nil
}

// GetClosedDeck Get a instance of deck that is closed for json purposes
// Close deck is a deck with deck.Cards hidden
func (d *Deck) GetClosedDeck() ClosedDeck {
	return ClosedDeck{
		DeckId:     d.DeckId,
		Shuffled:   d.Shuffled,
		Remaining:  d.Remaining,
	}
}

// GetOpenDeck Get a instance of deck that is open for json purposes
// Open deck is a deck with deck.Cards shown
func (d *Deck) GetOpenDeck() OpenDeck {
	return OpenDeck{
		DeckId:     d.DeckId,
		Shuffled:   d.Shuffled,
		Remaining:  d.Remaining,
		Cards:      d.Cards,
	}
}

// SetCards set the cards in the deck, should only be called at start
func (d *Deck) setCards(c []Card) {
	d.Cards = c
	d.Remaining = len(d.Cards)
}

func NewDeck() Deck {
	return Deck{
		DeckId:     uuid.NewString(),
		Shuffled:   false,
		Remaining:  0,
		Cards:      []Card{},
		DrawnCards: []Card{},
	}
}

// NewStandardDeck generates a standard french 52 card deck without jokers
func NewStandardDeck(shuffle bool) (*Deck, error) {
	deck := NewDeck()
	deck.setCards(DefaultDeckCards())

	if shuffle {
		err := deck.ShuffleDeck()
		if err != nil {
			return nil, err
		}
	}

	return &deck, nil
}

// NewStandardPartialDeck generates a standard deck with a subset of the default 52 cards
// will duplicate cards if same cardCode is specified twice
// returns error if cardCode is not found
func NewStandardPartialDeck(cardCodes []string, shuffle bool) (*Deck, error) {
	deck := NewDeck()

	availableCards := DefaultDeckCards()
	subDeck, err := SubDeck(availableCards, cardCodes)
	if err != nil {
		return nil, err
	}
	deck.setCards(subDeck)

	if shuffle {
		err := deck.ShuffleDeck()
		if err != nil {
			return nil, err
		}
	}

	return &deck, nil
}
