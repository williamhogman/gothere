package gothere
import (
	"appengine"
	"appengine/user"
	"net/http")

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
