package decks

import (
	"errors"
)

// DeckCardConfig configuration structure for a deck
// CardSuits and CardValues: the first letter of CardSuits and CardValues combined should be unique!
// Order of values affects initial order of deck
type DeckCardConfig struct {
	CardSuits  []string
	CardValues []string
	ExtraCards []Card
}

// default 52 card deck, no jokers
var defaultDeckConfig = DeckCardConfig{
	CardSuits:  []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"},
	CardValues: []string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING"},
	ExtraCards: []Card{},
}

// GenerateCards Generates all the cards for given DeckCardConfig
// Output order is dependent on order of strings in config
func GenerateCards(config DeckCardConfig) []Card {
	numCards := len(config.CardSuits) * len(config.CardValues) + len(config.ExtraCards)
	cards := make([]Card, numCards, numCards)
	i := 0
	for _, cardSuit := range config.CardSuits {
		for _, cardValue := range config.CardValues {

			cardCode := string(cardValue[0]) + string(cardSuit[0])

			cards[i] = Card{
				Value: cardValue,
				Suit:  cardSuit,
				Code:  cardCode,
			}
			i += 1
		}
	}

	// add extra cards, e.g joker
	for _, extraCard := range config.ExtraCards {
		cards[i] = extraCard
		i += 1
	}

	return cards
}

// DefaultDeckCards returns a default 52 card deck, no jokers
func DefaultDeckCards() []Card {
	return GenerateCards(defaultDeckConfig)
}

// FindCard returns index of card in a list of cards of the first card that matches cardCode
// returns error if not found and idx -1
func FindCard(cards []Card, cardCode string) (int, error) {
	for i, card := range cards {
		if card.Code == cardCode {
			return i, nil
		}
	}

	return -1, errors.New("card cardCode not found "  + cardCode)
}

// SubDeck returns a subset of the given cards in the order specified by cardCodes param
// allows duplicate cardCodes
// returns error if card code not found
func SubDeck(cards []Card, cardCodes []string) ([]Card, error) {
	cardSubset := make([]Card, len(cardCodes), len(cardCodes))

	for i, cardCode := range cardCodes {
		foundIdx, err := FindCard(cards, cardCode)
		if err != nil {
			return cardSubset, err
		}
		c := cards[foundIdx]
		cardSubset[i] = Card{
			Value: c.Value,
			Suit:  c.Suit,
			Code:  c.Code,
		}
	}

	return cardSubset, nil
}
