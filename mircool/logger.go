package mircool

import (
	"log"
	"time"
)

func Logger() HandleFunc {
	return func(c *Context) {
		now := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(now))
	}
}
