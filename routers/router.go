package routers

import (
	"github.com/fonzie1006/shortlink/pkg/middlewares"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

// 路由的结构体,如果后续有添加删除直接在这里改好,然后RouterInit中的循环中修改一下即可
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// 初始化路由
func RouterInit() *mux.Router {
	router := mux.NewRouter()
	// 使用中间件
	m := alice.New(middlewares.MuxRecoverHandler, middlewares.MuxLoggingHandler)

	// 如果设置为false,会导致/index访问404,一定要/index/才能访问
	router.StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(m.ThenFunc(route.HandlerFunc))
	}

	return router
}
