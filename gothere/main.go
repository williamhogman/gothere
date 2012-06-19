package gothere

import (
	"io"
	"appengine"
	"appengine/user"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"
)



func init() {
	http.HandleFunc("/login", do_login)
	http.HandleFunc("/urls", urls)
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func get_mt_qual(mediatype string) (string,float64,error) {
	tp,par,err := mime.ParseMediaType(mediatype)
	if err != nil {
		return "",0,err
	}
	q,err := strconv.ParseFloat(par["q"],32)
	if err != nil {
		q = 0
	}
	return tp,q,nil
}

func select_mediatype(mediatype string) string {
	var supported = [3]string{"application/json","text/xml","text/html"}
	accepted := strings.Split(mediatype,",")

	bestmime := ""
	bestqual := 0.0
	for _, acc := range accepted {
		tp,q,err := get_mt_qual(acc)
		if err != nil {
			continue
		}
		
		if bestmime == "" ||  q > bestqual {
			for _, mt := range supported {
				if mt == acc {
					bestmime = tp
					bestqual = q
					break
				}
			}
		}
	}
	return bestmime
}

func write_json(w io.Writer,obj interface{}) error{
	j, err := json.Marshal(obj);
	if err != nil {
		return err
	}
	if _, err := w.Write(j); err != nil {
		return err
	}
	return nil
}

func urls(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	switch r.Method {
	case "GET":
		accept := r.Header.Get("Accept")
		switch select_mediatype(accept) {
		case "text/html":
			redirs,err := get_redirs(c,0,100)
			
			dt, err := json.Marshal(redirs)
			if err != nil {
				//todo: fix error here
			}
			w.Write(dt)
			
		case "application/json":
			
		}
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
