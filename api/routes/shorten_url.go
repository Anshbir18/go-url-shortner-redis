package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/Anshbir18/go-url-shortner-redis/database"
	"github.com/Anshbir18/go-url-shortner-redis/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// define request
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

// define response
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *gin.Context) {
	var body request

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// ---------------- RATE LIMITING ----------------

	r2 := database.CreateClient(1)
	defer r2.Close()

	ip := c.ClientIP()

	val, err := r2.Get(database.Ctx, ip).Result()

	if err == redis.Nil {
		// first request from this IP
		if err := r2.Set(database.Ctx, ip, os.Getenv("API_QUOTA"), 30*time.Minute).Err(); err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}
	} else if err != nil {
		c.JSON(500, gin.H{"error": "redis error"})
		return
	} else {
		valInt, convErr := r2.Get(database.Ctx, ip).Int()
		if convErr != nil {
			c.JSON(500, gin.H{"error": "redis parse error"})
			return
		}

		if valInt <= 0 {
			ttl, _ := r2.TTL(database.Ctx, ip).Result()

			c.JSON(429, gin.H{
				"error":            "rate limit exceeded",
				"rate_limit_reset": ttl / time.Minute,
			})
			return
		}
	}

	// ---------------- VALIDATIONS ----------------

	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
		return
	}

	if !helpers.RemoveDomainError(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot shorten the domain itself"})
		return
	}

	// enforce https
	body.URL = helpers.EnforceHTTP(body.URL)

	// decrement quota
	_ = r2.Decr(database.Ctx, ip).Err()

	// -------- your remaining logic continues below --------
}
