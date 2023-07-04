package apirouter

import "github.com/gin-gonic/gin"

func Params(c *gin.Context) interface{} {
	data, _ := c.Get("_pp")
	return data
}
