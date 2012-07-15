package gothere

import (
	"mime"
	"net/http"
	"strconv"
	"strings"
)

func get_mt_qual(mediatype string) (string, float64, error) {
	tp, par, err := mime.ParseMediaType(mediatype)
	if err != nil {
		return "", 0, err
	}
	q, err := strconv.ParseFloat(par["q"], 32)
	if err != nil {
		q = 1
	}
	return tp, q, nil
}

func select_mediatype(mediatype string) string {
	var supported = [3]string{"application/json","text/html","application/xml"}
	accepted := strings.Split(mediatype, ",")

	bestmime := ""
	bestqual := 1.0
	for _, acc := range accepted {
		tp, q, err := get_mt_qual(acc)
		if err != nil {
			continue
		}

		if bestmime == "" || q > bestqual {
			for _, mt := range supported {
				if mt == acc {
					bestmime = tp
					bestqual = q
					break
				}
			}
		}
	}
	if bestmime != "" {
		return bestmime
	}
	return supported[0]
}

func GetMediaType(mediatype string) (string, error) {
	mt, _, err := mime.ParseMediaType(mediatype)
	return mt, err
}

func GetContentType(r *http.Request) (string, error) {
	return GetMediaType(r.Header.Get("Content-Type"))
}
