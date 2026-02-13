package routes

import (
	"net/http"
	"time"

	"github.com/Anshbir18/go-url-shortner-redis/api/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// define request
type request struct {
	URL         string        `json:"url`
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

	// implement rate limiting // we want a user to only be able to make 10 requests per 30 minutes if they exceed that we want to return an error message

	//check if the input is valid

	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
	}

	//check for domain erros

	if !helpers.RemoveDomainError(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot shorten the domain itself"})
	}

	//enforece https ssl
	body.URL = helpers.EnforceHTTP(body.URL)
}
