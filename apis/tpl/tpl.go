package tpl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func Tpl(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
