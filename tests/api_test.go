package tests

import (
	"decksframework/server"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	server.InitializeRoutes(r)
	return r
}

func TestApiDecks(t *testing.T) {
	t.Run("create deck", func(t *testing.T) {
		httptest.NewRequest()
	})
}
