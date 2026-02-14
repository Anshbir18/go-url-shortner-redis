package routes

import(
	"github.com/Anshbir18/go-url-shortner-redis/database"
	"github.com/go-redis/redis/v8"
	"github.com/gin-gonic/gin"
)

// we will get the short url from the user
//once we get that we will check in our db if it exists or not and return it

func ResolveURL(c *gin.Context){
	shortURL := c.Param("url")

	r := database.CreateClient(0)
	defer r.Close()

	// check if the short url exists in our db
	val, err := r.Get(database.Ctx, shortURL).Result()
	if err == redis.Nil {
		c.JSON(404, gin.H{"error": "short URL not found"})
		return
		
	} else if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
		
	}

	rInrc := database.CreateClient(1)
	defer rInrc.Close()

	_= rInrc.Incr(database.Ctx, "counter")
	c.Redirect(301,val)
	return

}