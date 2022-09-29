package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const HeaderAuthKey = "token"
const UserKey = "userKey"

func Auth(session SessionUsecase, ignorePath ...string) func(*gin.Context) {
	m := make(map[string]struct{})
	for _, v := range ignorePath {
		m[v] = struct{}{}
	}
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if _, ok := m[path]; ok {
			c.Next()
			return
		}
		token := c.Request.Header.Get(HeaderAuthKey)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusOK, NewResp(ErrCodeUnauthenticated, "invalid token"))
			return
		}
		account, ok := session.Get(token)
		if !ok {
			c.AbortWithStatusJSON(http.StatusOK, NewResp(ErrCodeUnauthenticated, "invalid token"))
			return
		}
		c.Set(UserKey, account)
		c.Next()
	}
}
