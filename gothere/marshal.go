package gothere
import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// Outputs an object
func outputObj(w io.Writer, mediatype string, obj interface{}) error {
	switch mediatype {
	case "text/html":
		if err := outputHtml(w, obj); err != nil {
			return err
		}
	case "application/json":
		if err := outputJson(w, obj); err != nil {
			return err
		}
	case "application/xml":
		if err := outputXml(w, obj); err != nil {
			return err
		}
	}
	return nil
}

// Outputs the object as HTML
func outputHtml(w io.Writer, obj interface{}) error {
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
		outputXml(w, targetobj)
	}
	return nil
}

type HtmlRenderable interface {
	RenderHtml() string
}

// Outputs the object as json
func outputJson(w io.Writer, obj interface{}) error {
	j, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err := w.Write(j); err != nil {
		return err
	}
	return nil
}

// Outputs the object as xml
func outputXml(w io.Writer, obj interface{}) error {
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

func ReadData(mt string,from io.Reader,into interface{}) error {
	switch mt {
	case "application/json":
		data, err := ioutil.ReadAll(from)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, into)
	default:
		return errors.New("Unknown type")
	}
	return nil
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
