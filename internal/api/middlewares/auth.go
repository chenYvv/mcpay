package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Claims struct {
	Aud    string `json:"aud"`
	Iat    int64  `json:"iat"`
	Iss    string `json:"iss"`
	UserId int    `json:"user_id"`
	// jwt.StandardClaims
}

func (c Claims) Valid() error {
	return nil
}

type RedisSSOUser struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		Authorization := c.Request.Header["Authorization"]
		//jwtKey := viper.GetString("server.key")
		// 本地测试
		//debugAuthorization := config.Get("app.debug_authorization")
		//if debugAuthorization != "" {
		//	Authorization = []string{debugAuthorization}
		//}
		if len(Authorization) == 0 {
			authErr(c)
			return
		}
		//
		//token := Authorization[0]
		//tokenClaims, err1 := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		//	return []byte(jwtKey), nil
		//})
		//if err1 != nil {
		//	authErr(c)
		//	return
		//}

		//uid := tokenClaims.Claims.(*Claims).UserId
		//uid := 1
		//user := models.User{}
		//user.GetUsersById(uid)
		//if user.Id == 0 {
		//	authErr(c)
		//	return
		//}

		//user := models.GetUsersByToken(token)
		//if user.Id == 0 {
		//	authErr(c)
		//	return
		//}

		//c.Set("current_user", *user)

		c.Next()
		return
	}
}

func authErr(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  -1,
		"message": "auth",
		"data":    gin.H{},
	})
}
