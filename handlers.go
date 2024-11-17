package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping_func(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "PONG",
	})
}
