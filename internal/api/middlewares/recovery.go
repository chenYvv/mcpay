package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			//if err := recover(); err != nil {
			//	logger.Error("recovery from panic",
			//		zap.Time("time", time.Now()), // 记录时间
			//		zap.Any("error", err),        // 记录错误信息
			//		zap.Stack("stacktrace"),      // 调用堆栈信息
			//	)
			//	//response.JSONAbort(c, -1, fmt.Sprintf("%v", err), nil, 500)
			//	c.JSON(http.StatusInternalServerError, "err")
			//	return
			//}
		}()
		c.Next()
	}
}
