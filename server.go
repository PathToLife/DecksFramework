package main

import (
	"decksframework/decks"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

/*
A Gin Api Server :)

Notes: Please branch out server routes into separate handler files when adding new functionality
Single file for now
*/

func httpServer(port string, debug bool) {
	if debug == false {
		gin.SetMode("release")
	}

	r := gin.Default()

	// in memory deck storage
	deckStore := decks.NewDeckStore()

	r.POST("/deck/create", func(c *gin.Context) {
		shouldShuffleDeck := false

		// ?shouldShuffleDeck=true for deck shuffling
		shuffleStr := c.DefaultQuery("shouldShuffleDeck", "false")
		if strings.ToLower(shuffleStr) == "true" {
			shouldShuffleDeck = true
		}

		var deck *decks.Deck

		// ?cards=AS,2S,KH
		if cardCodesStr, isCustom := c.GetQuery("cards"); isCustom {
			if cardCodesStr == "" {
				c.AbortWithStatusJSON(400, ErrorMsg("cards requires comma separated list of card codes"))
				return
			}
			cardCodes := strings.Split(cardCodesStr, ",")
			customDeck, err := decks.NewStandardPartialDeck(cardCodes, shouldShuffleDeck)
			if err != nil {
				c.AbortWithStatusJSON(500, ErrorMsg(err.Error()))
				return
			}
			deck = customDeck
		} else { // standard deck fi cards is not specified
			standardDeck, err := decks.NewStandardDeck(shouldShuffleDeck)
			if err != nil {
				_ = c.AbortWithError(500, err)
				return
			}
			deck = standardDeck
		}

		deckStore.AddDeck(deck)
		c.JSON(201, deck.GetClosedDeck())
	})

	r.POST("/deck/draw", func(c *gin.Context) {
		uuid, ok := c.GetQuery("uuid")
		if !ok || uuid == "" {
			c.AbortWithStatusJSON(400, ErrorMsg("uuid required"))
			return
		}

		countStr, ok := c.GetQuery("count")
		if !ok {
			c.AbortWithStatusJSON(400, ErrorMsg("count required"))
			return
		}
		drawCount, err := strconv.Atoi(countStr)
		if err != nil {
			c.AbortWithStatusJSON(400, ErrorMsg("count needs to be a number"))
			return
		}

		deck := deckStore.GetDeck(uuid)
		if deck == nil {
			c.AbortWithStatusJSON(404, ErrorMsg("deck not found"))
			return
		}

		cards, err := deck.DrawCard(drawCount)
		if err != nil {
			c.AbortWithStatusJSON(400, ErrorMsg(err.Error()))
			return
		}

		c.JSON(200, gin.H{
			"cards": cards,
		})
	})

	r.GET("/deck/open", func(c *gin.Context) {
		uuid, ok := c.GetQuery("uuid")
		if !ok || uuid == "" {
			c.AbortWithStatusJSON(400, ErrorMsg("uuid required"))
			return
		}

		deck := deckStore.GetDeck(uuid)
		if deck == nil {
			c.AbortWithStatusJSON(404, ErrorMsg("deck not found"))
			return
		}

		c.JSON(200, deck.GetOpenDeck())
	})

	// listen and serve
	err := r.Run(":" + port)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// ErrorMsg returns error message format for json
func ErrorMsg(msg string) gin.H {
	return gin.H{
		"error": msg,
	}
}
