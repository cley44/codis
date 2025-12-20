package middleware

import (
	"codis/utils/slogger"

	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
)

func GinOopsRecovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := oops.Recoverf(func() {
			ctx.Next()
		}, "gin: panic recovered")
		//&& !isBrokenPipeError(err)
		if err != nil {
			slogger.Error(err)
			ctx.JSON(500, "Could not process request")
		}
	}
}

// // Check for a broken connection, as this is what Gin does already.
// func isBrokenPipeError(err error) bool {
// 	// @TODO we should not determine it is a broken pipe error by string matching
// 	msg := strings.ToLower(err.Error())
// 	return strings.Contains(msg, "broken pipe") || strings.Contains(msg, "connection reset by peer")

// 	// if netErr, ok := err.(*net.OpError); ok {
// 	// 	if sysErr, ok := netErr.Err.(*os.SyscallError); ok {
// 	// 		if strings.Contains(strings.ToLower(sysErr.Error()), "broken pipe") ||
// 	// 			strings.Contains(strings.ToLower(sysErr.Error()), "connection reset by peer") {
// 	// 			return true
// 	// 		}
// 	// 	}
// 	// }

// 	// return false
// }
