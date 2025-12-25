package handlers

import "github.com/gin-gonic/gin"

func GetBody(c *gin.Context) interface{} {
	body, exist := c.Get("body")
	if !exist {
		return nil
	}
	return body
}
