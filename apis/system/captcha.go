package system

import (
	"ferry/tools"
	"ferry/tools/app"
	"ferry/tools/captcha"

	"github.com/gin-gonic/gin"
)

func GenerateCaptchaHandler(c *gin.Context) {
	id, b64s, err := captcha.DriverDigitFunc()
	tools.HasError(err, "验证码获取失败", 500)
	app.Custum(c, gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}
