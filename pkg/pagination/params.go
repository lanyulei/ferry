package pagination

import (
	"github.com/RichardKnop/machinery/v1/log"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RequestParams(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{}, 10)

	if c.Request.Form == nil {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			log.ERROR.Println(err)
		}
	}

	if len(c.Request.Form) > 0 {
		for key, value := range c.Request.Form {
			if key == "page" || key == "per_page" || key == "sort" {
				continue
			}
			params[key] = value[0]
		}
	}

	return params
}
