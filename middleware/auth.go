package middleware

import (
	"ferry/handler"
	jwt "ferry/pkg/jwtauth"
	"ferry/tools/config"
	"time"
)

func AuthInit() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte(config.ApplicationConfig.JwtSecret),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		PayloadFunc:     handler.PayloadFunc,
		IdentityHandler: handler.IdentityHandler,
		Authenticator:   handler.Authenticator,
		Authorizator:    handler.Authorizator,
		Unauthorized:    handler.Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

}
