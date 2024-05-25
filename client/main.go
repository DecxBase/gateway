package main

import (
	"net/http"
	"time"

	"github.com/DecxBase/gateway/node"
	"github.com/gin-gonic/gin"
)

func main() {
	loc, _ := time.LoadLocation("Asia/Kathmandu")
	options := &node.ClientOptions{
		Timezone: loc,
	}

	client := node.Client("auth", options).
		SetTitle("Authentication").
		SetVersion("0.0.1").
		SetDescription("Handles authentication services on api").
		SetMountVersion("v1").
		SetCategory("auth").
		SetTags([]string{"auth"}).
		SetDevice("local1")

	client.AddRoute(node.MethodGet, "/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to login",
		})
	})

	client.AddRoute(node.MethodPost, "/login", func(c *gin.Context) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Login action restricted",
		})
	})

	client.AddRoute(node.MethodGet, "/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to register",
		})
	})

	client.AddRoute(node.MethodGet, "/logout", func(c *gin.Context) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Logout restricted",
		})
	})

	client.Run("http://localhost:3332", node.RegistryOptions{
		BindAddress: "localhost:3000",
	})
}
