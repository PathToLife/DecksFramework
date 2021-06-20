package main

import (
	"decksframework/decks"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

/*
A Gin Api Server :)

Notes: Please branch out server routes into separate handler files when adding new functionality
Single file for now
*/

var router *gin.Engine

func startApiServer(port string, debug bool) {
	if debug == false {
		gin.SetMode("release")
	}

	router = gin.Default()

	initializeRoutes()

	// listen and serve
	err := router.Run(":" + port)
	if err != nil {
		fmt.Println(err)
	}
}

func initializeRoutes() {
	router.POST("/deck/create", func(c *gin.Context) {
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
				c.AbortWithStatusJSON(http.StatusBadRequest, ErrorMsg("cards requires comma separated list of card codes"))
				return
			}
			cardCodes := strings.Split(cardCodesStr, ",")
			customDeck, err := decks.NewStandardPartialDeck(cardCodes, shouldShuffleDeck)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, ErrorMsg(err.Error()))
				return
			}
			deck = customDeck
		} else { // standard deck if ?cards= is not specified
			standardDeck, err := decks.NewStandardDeck(shouldShuffleDeck)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			deck = standardDeck
		}

		deckStore.AddDeck(deck)
		c.JSON(http.StatusCreated, deck.GetClosedDeck())
	})

	router.POST("/deck/draw", func(c *gin.Context) {
		uuid, ok := c.GetQuery("uuid")
		if !ok || uuid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorMsg("uuid required"))
			return
		}

		countStr, ok := c.GetQuery("count")
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorMsg("count required"))
			return
		}
		drawCount, err := strconv.Atoi(countStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorMsg("count needs to be a number"))
			return
		}

		deck := deckStore.GetDeck(uuid)
		if deck == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorMsg("deck not found"))
			return
		}

		cards, err := deck.DrawCard(drawCount)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorMsg(err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"cards": cards,
		})
	})

	router.GET("/deck/open", func(c *gin.Context) {
		uuid, ok := c.GetQuery("uuid")
		if !ok || uuid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorMsg("uuid required"))
			return
		}

		deck := deckStore.GetDeck(uuid)
		if deck == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorMsg("deck not found"))
			return
		}

		c.JSON(http.StatusOK, deck.GetOpenDeck())
	})
}

// ErrorMsg returns error message format for json
func ErrorMsg(msg string) gin.H {
	return gin.H{
		"error": msg,
	}
}
