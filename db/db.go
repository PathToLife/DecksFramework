package db

/*
Normally would have db connection handing function here, but for now, we just have a DeckStore inmemory storage
*/

var deckStore DeckStore

func InitDb() {
	deckStore = NewDeckStore()
}

func GetDeckStore() DeckStore {
	return deckStore
}
