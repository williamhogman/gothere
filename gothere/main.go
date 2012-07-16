package gothere

import (
	"appengine"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/login", do_login)
	http.HandleFunc("/", handler)
	HandleMethods("/urls/",UrlsGet,UrlsPost)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}




func UrlsGet(x *ReqContext) {
	c := appengine.NewContext(x.r)

	redirs, err := GetRedirs(c, 0, 100)

	if err != nil {
		x.ServerError(err)
		return
	}
	x.Ok(redirs)
}

func UrlsPost(x *ReqContext) {
	c := appengine.NewContext(x.r)
	var data NamedRedir
	if err := x.Data(&data); err != nil {
		x.BadRequest(err)
		return
	}

	if err := data.Validate(); err != nil {
		x.BadRequest(err)
		return
	}
	if err := data.SaveNew(c); err != nil {
		x.BadRequest(err)
		return
	}
	x.Created(data.Location())
}


