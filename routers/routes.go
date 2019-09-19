package routers

import "github.com/fonzie1006/shortlink/routers/api/v1"

// 这里顺序不能错,最后有点默认路由的意思,如果上面的都不匹配,就会都去最后一个
var routes = Routes{
	Route{
		"LongPostToShortLink",
		"Post",
		"/shortlink",
		v1.LongPostToShortLink,
	},
	Route{
		"ShortPostToLongLink",
		"GET",
		"/shortlink/{shortlink:[a-zA-Z0-9]{4,6}}",
		v1.ShortGetToLongLink,
	},
	Route{
		"ShortRedirect",
		"GET",
		"/{shortlink:[a-zA-Z0-9]{4,6}}",
		v1.ShortRedirect,
	},
}
