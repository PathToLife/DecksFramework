package db

import "decksframework/decks"

// DeckStore stores decks, placeholder struct for future functionality such as persistent storage
type DeckStore struct {
	decksMap map[string]*decks.Deck
}

func (ds DeckStore) AddDeck(d *decks.Deck) {
	ds.decksMap[d.DeckId] = d
}

func (ds DeckStore) RemoveDeck(d *decks.Deck) {
	delete(ds.decksMap, d.DeckId)
}

func (ds DeckStore) GetDeck(uuid string) *decks.Deck {
	return ds.decksMap[uuid]
}

// NewDeckStore initialize deckStore
func NewDeckStore() DeckStore {
	return DeckStore{
		decksMap: map[string]*decks.Deck{},
	}
}
