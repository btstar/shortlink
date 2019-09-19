package middlewares

import (
	logging "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// 请求日志处理中间件
func MuxLoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(write, request)
		t2 := time.Now()
		logging.Infof("[%s] %v %s Header=%s", request.Method, t2.Sub(t1), request.URL.String(), request.Header)
	})
}

// 程序奔溃恢复
func MuxRecoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logging.Errorf("Recover from panic : %+v", err)
			}
		}()

		next.ServeHTTP(write, request)
	})
}
