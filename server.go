package main

import (
	"decks/decks"
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

func ErrorMsg(msg string) gin.H {
	return gin.H{
		"error": msg,
	}
}

func httpServer(port string, debug bool) {
	if debug == false {
		gin.SetMode("release")
	}

	r := gin.Default()

	// deck
	deckStore := decks.NewDeckStore()

	r.POST("/deck/create", func(c *gin.Context) {
		shuffle := false
		shuffleQ := c.DefaultQuery("shuffle", "false")

		if strings.ToLower(shuffleQ) == "true" {
			shuffle = true
		}

		deck, err := decks.NewStandardDeck(shuffle)
		if err != nil {
			_ = c.AbortWithError(500, err)
			return
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
		count, err := strconv.Atoi(countStr)
		if err != nil {
			c.AbortWithStatusJSON(400, ErrorMsg("count needs to be a number"))
			return
		}

		deck := deckStore.GetDeck(uuid)
		if deck == nil {
			c.AbortWithStatusJSON(404, ErrorMsg("deck not found"))
			return
		}

		cards, err := deck.DrawCard(count)
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