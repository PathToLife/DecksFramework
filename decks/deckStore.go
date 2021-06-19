package decks

// DeckStore stores decks, placeholder struct for future functionality such as persistent storage
type DeckStore struct {
	decksMap map[string]*Deck
}

func (ds DeckStore) AddDeck(d *Deck) {
	ds.decksMap[d.DeckId] = d
}

func (ds DeckStore) RemoveDeck(d *Deck) {
	delete(ds.decksMap, d.DeckId)
}

func (ds DeckStore) GetDeck(uuid string) *Deck {
	return ds.decksMap[uuid]
}

// NewDeckStore initialize deckStore
func NewDeckStore() DeckStore {
	return DeckStore{
		decksMap: map[string]*Deck{},
	}
}
