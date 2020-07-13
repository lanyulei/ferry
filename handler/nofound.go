package handler

import (
	jwt "ferry/pkg/jwtauth"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NoFound(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Printf("NoRoute claims: %#v\n", claims)
	c.JSON(http.StatusOK, gin.H{
		"code":    "404",
		"message": "not found",
	})
}
