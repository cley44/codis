package utils

import (
	"github.com/gin-gonic/gin"
)

func AbortRequest(c *gin.Context, statusCode int, err error, msg string) {
	output := buildHTTPErrorBody(c, err, gin.H{
		"message": msg,
	})

	c.AbortWithStatusJSON(statusCode, output)
}

func AbortRequestWithCode(c *gin.Context, statusCode int, err error, code string, msg string) {
	output := buildHTTPErrorBody(c, err, gin.H{
		"code":    code,
		"message": msg,
	})

	c.AbortWithStatusJSON(statusCode, output)
}

func AbortRequestWithCodeAndData(c *gin.Context, statusCode int, err error, code string, msg string, data map[string]interface{}) {
	output := buildHTTPErrorBody(c, err, gin.H{
		"code":    code,
		"message": msg,
		"data":    data,
	})

	c.AbortWithStatusJSON(statusCode, output)
}

// format body.
func buildHTTPErrorBody(c *gin.Context, err error, body gin.H) gin.H {
	if err != nil {
		c.Error(err) //nolint:errcheck
	}

	return body
}
