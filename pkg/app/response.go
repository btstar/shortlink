package app

import (
	"encoding/json"
	"github.com/fonzie1006/shortlink/pkg/e"
	"github.com/gorilla/mux"
	"net/http"
)

type AppR struct {
	Write   http.ResponseWriter
	Request *http.Request
	r       *mux.Router
}

func (app *AppR) Response(httpCode, errCode int, data interface{}) {
	content := map[string]interface{}{
		"Code": errCode,
		"Msg":  e.GetMsg(errCode),
		"Data": data,
	}

	js, err := json.Marshal(content)
	if err != nil {
		http.Error(app.Write, err.Error(), http.StatusInternalServerError)
		return
	}
	app.Write.Header().Set("Content-Type", "application/json")
	app.Write.WriteHeader(httpCode)
	app.Write.Write(js)
}
