package gothere

import (
	"net/http"
)





type ErrorWrap struct {
	Error string
}


func EWrap(obj interface{}) *ErrorWrap {
	return &ErrorWrap{Error: StrOrErr(obj)}
}

func StrOrErr(obj interface{}) string {
	switch obj.(type) {
	case string:
		return obj.(string)
	case error:
		return obj.(error).Error()
	}
	return ""
}

func E500(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func E400(w http.ResponseWriter, s interface{}) {
	http.Error(w, StrOrErr(s), http.StatusBadRequest)
}

func E500IfErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		E500(w, err)
		return true
	}
	return false
}
