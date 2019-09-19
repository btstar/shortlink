package v1

import (
	"github.com/fonzie1006/shortlink/pkg/app"
	"github.com/fonzie1006/shortlink/pkg/e"
	"github.com/fonzie1006/shortlink/services/conversion"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type LToSLink struct {
	// 需要被转换成短链接的长链接
	Url string `form:"url" valid:"Required"`
	//
	Timeout int `form:"timeout" valid:"Required"`
}

// 长链接转短链接
func LongPostToShortLink(writer http.ResponseWriter, request *http.Request) {
	appMux := app.AppR{
		Write:   writer,
		Request: request,
	}
	form := LToSLink{}
	request.ParseForm()
	if v, ok := request.Form["url"]; ok {
		if len(v) >= 1 {
			form.Url = v[0]
		}
	}
	if v, ok := request.Form["timeout"]; ok {
		if len(v) >= 1 {
			var err error
			form.Timeout, err = strconv.Atoi(v[0])
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	httpCode, errCode := app.BindAndValid(request, &form)
	if errCode != e.SUCCESS {
		appMux.Response(httpCode, errCode, nil)
		return
	}

	shortLinkCon := conversion.ShortLinkConversion{
		LongLink: form.Url,
		Timeout:  form.Timeout,
	}

	err := shortLinkCon.GetShortLink()
	if err != nil {
		appMux.Response(http.StatusInternalServerError, e.CREATE_SHORTLINK_ERROR, nil)
		return
	}

	appMux.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"shortlink": "/" + shortLinkCon.Shortlink,
		"longlink":  shortLinkCon.LongLink,
		"md5":       shortLinkCon.Md5,
		"timeout":   shortLinkCon.Timeout,
	})

}

// 短链接转长链接
func ShortGetToLongLink(writer http.ResponseWriter, request *http.Request) {
	appMux := app.AppR{
		Write:   writer,
		Request: request,
	}

	vars := mux.Vars(request)
	shorlink := ""
	if v, ok := vars["shortlink"]; ok {
		shorlink = v
	} else {
		appMux.Response(http.StatusInternalServerError, e.SHORT_LINK_IS_MUST, nil)
		return
	}

	shortLinkCon := conversion.ShortLinkConversion{
		Shortlink: shorlink,
	}
	err := shortLinkCon.GetLongLink()
	if err != nil {
		appMux.Response(http.StatusBadRequest, e.SHORT_TO_LONG_LINK_ERROR, nil)
		return
	}

	appMux.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"shortlink": "/" + shortLinkCon.Shortlink,
		"longlink":  shortLinkCon.LongLink,
		"md5":       shortLinkCon.Md5,
		"timeout":   shortLinkCon.Timeout,
	})

}

// 短链接Redicre 302
func ShortRedirect(writer http.ResponseWriter, request *http.Request) {
	appMux := app.AppR{
		Write:   writer,
		Request: request,
	}

	vars := mux.Vars(request)
	shorlink := ""
	if v, ok := vars["shortlink"]; ok {
		shorlink = v
	} else {
		appMux.Response(http.StatusInternalServerError, e.SHORT_LINK_IS_MUST, nil)
		return
	}

	shortLinkCon := conversion.ShortLinkConversion{
		Shortlink: shorlink,
	}
	err := shortLinkCon.GetLongLink()
	if err != nil {
		appMux.Response(http.StatusBadRequest, e.SHORT_TO_LONG_LINK_ERROR, nil)
		return
	}
	// 记录短链接访问次数,如果是301就无法完成记录了
	shortLinkCon.CountShortLinkAccess()

	http.Redirect(writer, request, shortLinkCon.LongLink, 302)
}
