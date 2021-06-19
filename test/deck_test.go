package test

import (
	"decks/decks"
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestCardGeneration(t *testing.T) {
	cards := decks.DefaultDeckCards()

	t.Run("has 52 cards", func(t *testing.T) {
		if len(cards) != 52 {
			t.Errorf("len(cards)=%d; expected 52 cards", len(cards))
		}
	})

	t.Run("generated card codes are unique for default deck", func(t *testing.T) {
		codes := make(map[string]bool)
		for _, card := range cards {
			_, value := codes[card.Code]
			if value {
				t.Errorf("card.Code not unique, duplicate card.Code: \"%s\"", card.Code)
			}
			codes[card.Code] = value
		}
	})

	t.Run("all cards have non empty string in all fields", func(t *testing.T) {
		for i, card := range cards {
			if card.Value == "" {
				t.Errorf("card.Value[%d]=\"%s\"; expected not empty", i, card.Value)
			}
			if card.Suit == "" {
				t.Errorf("card.Suit[%d]=\"%s\"; expected not empty", i, card.Suit)
			}
			if card.Code == "" {
				t.Errorf("card.Code[%d]=\"%s\"; expected not empty", i, card.Code)
			}
		}
	})
}

func TestDefaultDeck(t *testing.T) {
	deckUnShuffled, _ := decks.NewStandardDeck(false)

	t.Run("deck should have valid uuid", func(t *testing.T) {
		_, err := uuid.Parse(deckUnShuffled.DeckId)
		if err != nil {
			t.Errorf("deckUnShuffled.DeckId is not a uuid; err=%s", err.Error())
		}
	})

	t.Run("has 52 cards in deckUnShuffled", func(t *testing.T) {
		if len(deckUnShuffled.Cards) != 52 {
			t.Errorf("len(deckUnShuffled.Cards)=%d; expected 52", len(deckUnShuffled.Cards))
		}
	})

	t.Run("deckUnShuffled is unShuffled", func(t *testing.T) {
		if deckUnShuffled.Shuffled {
			t.Errorf("deckUnShuffled.Shuffled=%t; expected false", deckUnShuffled.Shuffled)
		}
	})

	deckShuffled, _ := decks.NewStandardDeck(true)

	t.Run("deck uuids should not be the same", func(t *testing.T) {
		if deckUnShuffled.DeckId == deckShuffled.DeckId {
			t.Errorf("deckUnShuffled.DeckId=%s deckShuffled.DeckId=%s; expected not equal", deckUnShuffled.DeckId, deckShuffled.DeckId)
		}
	})

	t.Run("deck is shuffled", func(t *testing.T) {
		if !deckShuffled.Shuffled {
			t.Errorf("deckShuffled.Shuffled=%t; expected true", deckShuffled.Shuffled)
		}

		cardsEqual := CardsEqual(deckUnShuffled.Cards, deckShuffled.Cards)
		if cardsEqual {
			t.Errorf("cardsEqual=%t; expected false", cardsEqual)
		}

		//s, _ := json.MarshalIndent(deckShuffled.Cards, "", "  ")
		//fmt.Println(string(s))
	})

	t.Run("shuffled deck right length", func(t *testing.T) {
		if len(deckUnShuffled.Cards) != 52 {
			t.Errorf("len(deckShuffled.Cards)=%d; expected 52", len(deckUnShuffled.Cards))
		}
	})

	t.Run("shuffled deck no duplicates", func(t *testing.T) {
		codes := make(map[string]bool)
		for _, card := range deckShuffled.Cards {
			_, value := codes[card.Code]
			if value {
				t.Errorf("card code not unique, duplicate code: \"%s\"", card.Code)
			}
			codes[card.Code] = value
		}
	})
}

func DrawCards(t *testing.T, deck *decks.Deck, drawCount int) {
	remainingBefore := deck.Remaining
	remainingAfter := remainingBefore - drawCount

	card, err := deck.DrawCard(drawCount)

	if err != nil {
		t.Errorf("draw %d card, returned error %s", drawCount, err.Error())
	}
	if len(card) != drawCount {
		t.Errorf("len(card)=%d; expected %d card", len(card), drawCount)
	}
	if deck.Remaining != remainingAfter {
		t.Errorf("deck.Remaining=%d; expected %d", deck.Remaining, remainingAfter)
	}
}

func TestCardDrawing(t *testing.T) {
	deck, _ := decks.NewStandardDeck(true)

	t.Run("draw 1 card", func(t *testing.T) {
		DrawCards(t, deck, 1)
	})

	t.Run("draw 5 cards", func(t *testing.T) {
		DrawCards(t, deck, 5)
	})

	t.Run("draw 0 cards", func(t *testing.T) {
		remainingBefore := deck.Remaining
		card, _ := deck.DrawCard(0)
		if remainingBefore != deck.Remaining {
			t.Errorf("deck.remaining=%d; expected %d", deck.Remaining, remainingBefore)
		}
		if len(card) != 0 {
			t.Errorf("len(card)=%d; expected 0 cards drawn", len(card))
		}
	})

	t.Run("draw all cards", func(t *testing.T) {
		deck1, _ := decks.NewStandardDeck(false)
		DrawCards(t, deck1, 52)
	})

	t.Run("draw when all cards drawn", func(t *testing.T) {
		deck1, _ := decks.NewStandardDeck(false)
		DrawCards(t, deck1, 52)
		if deck1.Remaining != 0 {
			t.Errorf("deck1.Remaining=%d; expected 0", deck1.Remaining)
		}
		cards, err := deck1.DrawCard(1)
		if err == nil {
			t.Errorf("err=nil; expected some error when overdrawing cards")
		}
		if len(cards) != 0 {
			t.Errorf("len(card)=%d; expected 0 cards drawn", len(cards))
		}
	})

	t.Run("draw too many cards", func(t *testing.T) {
		remainingBefore := deck.Remaining
		deck1, _ := decks.NewStandardDeck(false)
		cards, err := deck1.DrawCard(deck1.Remaining + 1)
		if err == nil {
			t.Errorf("err=nil; expected some error")
		}
		if len(cards) != 0 {
			t.Errorf("len(card)=%d; expected 0 cards drawn", len(cards))
		}
		if remainingBefore != deck.Remaining {
			t.Errorf("deck.remaining=%d; expected %d", deck.Remaining, remainingBefore)
		}
	})
}

func TestCustomDeck(t *testing.T) {

	t.Run("shuffle empty deck", func(t *testing.T) {
		emptyDeck, _ := decks.NewStandardPartialDeck([]string{}, false)
		if len(emptyDeck.Cards) != 0 {
			t.Errorf("len(card)=%d; expected 0 cards", len(emptyDeck.Cards))
		}
		_ = emptyDeck.ShuffleDeck()
		if len(emptyDeck.Cards) != 0 {
			t.Errorf("len(card)=%d; expected 0 cards after shuffling", len(emptyDeck.Cards))
		}
	})

	t.Run("shuffle 2 card deck", func(t *testing.T) {
		deck, err := decks.NewStandardPartialDeck([]string{"AS", "1S"}, true)
		if err != nil {
			t.Errorf("%s", err.Error())
		} else if deck != nil && len(deck.Cards) != 2 {
			t.Errorf("len(card)=%d; expected 2 cards", len(deck.Cards))
		}
	})

	t.Run("should return error A0 card not found", func(t *testing.T) {
		_, err := decks.NewStandardPartialDeck([]string{"A0", "1S"}, true)
		if err == nil {
			t.Errorf("err=nil; expected some error")
		} else if !strings.Contains(err.Error(), "A0") {
			t.Errorf("err.Error()=%s; expected to contain A0", err.Error())
		}
	})

	t.Run("custom deck with duplicate cards", func(t *testing.T) {
		deck, _ := decks.NewStandardPartialDeck([]string{"KH", "KH"}, false)

		if len(deck.Cards) != 2 {
			t.Errorf("len(deck.Cards)=%d; expected 2", len(deck.Cards))
		}

		for _, card := range deck.Cards {
			if card.Code != "KH" {
				t.Errorf("card.Code=%s; expected KH", card.Code)
				break
			}
		}
	})
}

func TestGenerateCards(t *testing.T) {

	t.Run("extra card test", func(t *testing.T) {
		extraCards := []decks.Card{
			{
				Value: "GodFather",
				Suit:  "Mafia",
				Code:  "GM",
			}, {
				Value: "Doctor",
				Suit:  "Town",
				Code:  "DT",
			},
		}

		cards := decks.GenerateCards(decks.DeckCardConfig{
			ExtraCards: extraCards,
		})

		if len(cards) != 2 {
			t.Errorf("len(cards)=%d; expected 2", len(cards))
		}

		if cards[0].Code != extraCards[0].Code {
			t.Errorf("cards[0].Code=%s; expected %s", cards[0].Code, extraCards[0].Code)
		}

		if cards[1].Code != extraCards[1].Code {
			t.Errorf("cards[1].Code=%s; expected %s", cards[1].Code, extraCards[1].Code)
		}
	})
}