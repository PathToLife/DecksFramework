package test

import "decks/decks"

func CardsEqual(c1 []decks.Card, c2 []decks.Card) bool {
	if len(c1) != len(c2) {
		return false
	}

	for i := range c1 {
		cc1 := c1[i]
		cc2 := c2[i]
		if cc1.Code != cc2.Code {
			return false
		}
	}

	return true
}