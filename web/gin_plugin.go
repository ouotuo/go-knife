package web

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/cyfdecyf/bufio"
	"bytes"
	log "github.com/sirupsen/logrus"
	"fmt"
)

const (
	KEY_GIN_LOGGER_DEBUG="gin_logger_debug"
	KEY_GIN_LOGGER_BODY="gin_logger_body"
)

type bufferedWriter struct {
	gin.ResponseWriter
	out    *bufio.Writer
	Buffer bytes.Buffer
}

func (g *bufferedWriter) Write(data []byte) (int, error) {
	g.Buffer.Write(data)
	return g.out.Write(data)
}

func GetGinLogger(printBody bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.String()
		method := c.Request.Method
		clientIP := c.ClientIP()

		logId:=fmt.Sprintf("%d",start.UnixNano())

		if printBody{
			c.Set(KEY_GIN_LOGGER_DEBUG,"true")
			w := bufio.NewWriter(c.Writer)
			buff := bytes.Buffer{}
			newWriter := &bufferedWriter{c.Writer, w, buff}

			c.Writer = newWriter

			defer func() {
				w.Flush()
				body, existsBody := c.Get(KEY_GIN_LOGGER_BODY)
				responseBody:=newWriter.Buffer.Bytes()

				if existsBody {
					log.WithFields(log.Fields{
						"tag":"gin.req",
						"logId":logId,
						"contentType":c.ContentType(),
					}).Infoln(body)
				}

				if len(responseBody)>0 {
					log.WithFields(log.Fields{
						"tag":"gin.res",
						"logId":logId,
					}).Infoln(string(responseBody))
				}
			}()
		}

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()

		en:=log.WithFields(log.Fields{
			"tag":"gin.path",
			"logId":logId,
			"latency":latency,
			"ip":clientIP,
			"method":method,
			"path":path,
			"statusCode":statusCode,
		})

		if statusCode<500{
			en.Infof("%s %s %d %s",method,path,statusCode,latency)
		}else{
			en.Errorf("%s %s %d %s",method,path,statusCode,latency)
		}
	}
}
