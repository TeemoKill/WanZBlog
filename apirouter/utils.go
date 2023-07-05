package apirouter

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/TeemoKill/WanZBlog/constants"
	"github.com/gin-gonic/gin"
)

type NotFoundResponse struct{}

func Params(c *gin.Context) interface{} {
	data, _ := c.Get("_pp")
	return data
}

func GenerateLoginToken(userUUID string) (token string, err error) {
	buffer := make([]byte, constants.LoginTokenLength)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return hex.EncodeToString(buffer), err
}
