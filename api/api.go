package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"vpeel/internal/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var routerInitList = make([]func(*gin.Engine), 0)

func AddRouter(f func(*gin.Engine)) {
	routerInitList = append(routerInitList, f)
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func logger(c *gin.Context) {
	var reqBodyBytes []byte
	var respBodyString string
	w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = w
	if strings.Contains(c.Request.Header.Get("Content-Type"), "json") && c.Request.Body != nil {
		reqBodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
	}
	c.Next()

	if strings.Contains(c.Writer.Header().Get("Content-Type"), "json") {
		respBodyString = w.body.String()
	}
	
	log.LoggerAccess.Info("access",
		zap.Int("status",c.Writer.Status()),
		zap.String("uri", c.Request.RequestURI),
		zap.String("addr", c.ClientIP()),
		zap.String("client_ip", c.Request.RemoteAddr),
		zap.String("args", c.Request.URL.RawQuery),
		zap.String("method", c.Request.Method),
		zap.Any("req_header", c.Request.Header),
		zap.ByteString("req_body", reqBodyBytes),
		zap.String("resp_body", respBodyString),
		zap.String("resp_header", fmt.Sprintf("%v", c.Writer.Header())),
	)
}

func Run(addr string) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	
	router := gin.Default()
	router.Use(logger, gin.Recovery())

	for _, rfunc := range routerInitList {
		rfunc(router)
	}
	fmt.Println("api listen on", addr)
	router.Run(addr) 
}
