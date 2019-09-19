package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/fonzie1006/shortlink/pkg/e"
	"net/http"
	//httpfrom "github.com/smartwalle/form"
)

func BindAndValid(request *http.Request, form interface{}) (int, int) {
	//err := httpfrom.Bind(request.Form, form)
	//if err != nil {
	//	return http.StatusInternalServerError, e.INVALID_PARAMS
	//}
	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}
