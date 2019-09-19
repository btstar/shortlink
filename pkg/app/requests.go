package app

import (
	"github.com/astaxie/beego/validation"
	logging "github.com/sirupsen/logrus"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Error(err)
	}

	return
}
