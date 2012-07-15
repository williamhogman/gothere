package gothere

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func write_json(w io.Writer, obj interface{}) error {
	j, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err := w.Write(j); err != nil {
		return err
	}
	return nil
}

func write_xml(w io.Writer, obj interface{}) error {
	io.WriteString(w, xml.Header)
	x, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err := w.Write(x); err != nil {
		return err
	}
	return nil
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

func output_obj(w io.Writer, mediatype string, obj interface{}) error {
	switch mediatype {
	case "text/html":
		if err := write_html(w, obj); err != nil {
			return err
		}
	case "application/json":
		if err := write_json(w, obj); err != nil {
			return err
		}
	case "application/xml":
		if err := write_xml(w, obj); err != nil {
			return err
		}
	}
	return nil
}

func Output(w http.ResponseWriter, mediatype string, obj interface{}) error {
	w.Header().Set("Content-Type", mediatype)
	return output_obj(w, mediatype, obj)
}

func Created(w http.ResponseWriter, loc string) {
	w.Header().Set("Location", loc)
	w.WriteHeader(http.StatusCreated)
}

func write_html(w io.Writer, obj interface{}) error {
	switch targetobj := obj.(type) {
	case []HtmlRenderable:
		for _, rend := range targetobj {
			io.WriteString(w, rend.RenderHtml())
		}
	case HtmlRenderable:
		io.WriteString(w, targetobj.RenderHtml())
	case string:
		io.WriteString(w, targetobj)
	default:
		write_xml(w, targetobj)
	}
	return nil
}

type HtmlRenderable interface {
	RenderHtml() string
}

func ExtractData(r *http.Request, into interface{}) error {
	mt, err := GetContentType(r)
	if err != nil {
		return err
	}
	switch mt {
	case "application/json":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, into)
	default:
		return errors.New("Unknown type")
	}
	return nil
}
