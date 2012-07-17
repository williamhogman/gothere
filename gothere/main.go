package gothere

import (
	"appengine"
	"appengine/datastore"
	"net/http"
)

func init() {
	http.HandleFunc("/login", do_login)

	collection := NewMethodHandler()
	collection.Get = UrlsGet
	collection.Post = UrlsPost

	resource := NewMethodHandler()
	resource.Get = UrlGet

	RESTHandler("/urls/",collection,resource)
	SimpleHandler("/",Index)
}

func Index(x *ReqContext) {
	c := appengine.NewContext(x.r)
	
	if x.r.URL.Path == "/" {
		x.Ok("Index page")
		return
	}

	redir, err := GetRedir(c, x.r.URL.Path[1:])
	
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			x.NotFound(err)
		} else {
			x.ServerError(err)
		}
		return
	}
	
	x.SeeOther(redir.Target)
	return
}

func UrlGet(x *ReqContext) {
	c := appengine.NewContext(x.r)

	path := x.r.URL.Path
	ident := path[len("/urls/"):]
	
	redir, err := GetRedir(c, ident)
	
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			x.NotFound(err)
		} else {
			x.ServerError(err)
		}
		return
	}

	x.Ok(redir)
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

	if err := x.ValidatedData(&data); err != nil {
		x.BadRequest(err)
		return
		
	}

	if err := data.SaveNew(c); err != nil {
		x.BadRequest(err)
		return
	}
	x.Created(data.Location())
}


