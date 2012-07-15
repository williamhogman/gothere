package gothere

import (
	"appengine"
	"appengine/user"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/login", do_login)
	http.HandleFunc("/urls", urls)
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func urls(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accept := r.Header.Get("Accept")
	output_t := select_mediatype(accept)

	switch r.Method {
	case "GET":
		redirs, err := GetRedirs(c, 0, 100)

		if E500IfErr(w, err) {
			return
		}

		if err := Output(w, output_t, redirs); err != nil {
			E500(w, err)
			return
		}
	case "POST":
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
}

func do_login(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	} else {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
	}
}
