package gothere

import (
	"appengine"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/login", do_login)
	//http.HandleFunc("/urls/", urls)
	http.HandleFunc("/", handler)
	HandleRest("/urls/",UrlsGet,UrlsPost)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}


type FnHandler func(http.ResponseWriter,*http.Request)
func HandleRest(path string,fnget FnHandler,fnpost FnHandler) {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fnget(w,r)
		case "POST":
			fnpost(w,r)
		}
	}
	http.HandleFunc(path,wrapper)
}


func UrlsGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accept := r.Header.Get("Accept")
	output_t := select_mediatype(accept)

	redirs, err := GetRedirs(c, 0, 100)

	if E500IfErr(w, err) {
		return
	}

	if err := Output(w, output_t, redirs); err != nil {
		E500(w, err)
		return
	}
}

func UrlsPost(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var data NamedRedir
	if err := ExtractData(r, &data); err != nil {
		E400(w, err)
		return
	}

	if err := data.Validate(); err != nil {
		E400(w, err)
		return
	}
	if err := data.SaveNew(c); err != nil {
		E500(w, err)
	}

	Created(w, data.Location())
}


