package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

// https://medium.com/@nhw2n0lg/graylog-implementation-in-golang-gin-2a14176302e5
func LoggingMiddleware(gelfwriter *gelf.TCPWriter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime).Milliseconds()
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		gelfwriter.WriteMessage(
			&gelf.Message{
				Host:     "API",
				Short:    reqMethod,
				TimeUnix: float64(endTime.Unix()),
				Extra: map[string]interface{}{
					"_reqUri":   reqUri,
					"_clientIP": clientIP,
					"_status":   statusCode,
					"_latency":  latencyTime,
				},
			},
		)

		ctx.Next()
	}
}
