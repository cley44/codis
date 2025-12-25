package middleware

import (
	"codis/utils"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ValidateBodyMiddleware(body any) gin.HandlerFunc {
	t := reflect.TypeOf(body)

	return func(ctx *gin.Context) {
		ptr := reflect.New(t)
		err := ctx.ShouldBindBodyWith(ptr.Interface(), binding.JSON)
		ctx.Set("body", ptr.Interface())
		if err != nil {
			utils.AbortRequest(ctx, http.StatusUnprocessableEntity, err, "")
			return
		}
	}
}
