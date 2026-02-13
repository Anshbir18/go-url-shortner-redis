package ma

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Anshbir18/go-url-shortner-redis/routes"
)


func setupRoutes(router *gin.Engine){

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/:url", routes.ResolveURL)
	router.POST("/shorten", routes.ShortenURL)

}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	router :=gin.Default()

	log.Fatal(router.Run(os.Getenv("APP_PORT")))
}