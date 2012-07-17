// My tiny kit of http extras for Go
package gothere
import (
	"net/http"
)

type FnHandler func(*ReqContext)


type ReqContext struct {
	w http.ResponseWriter
	r *http.Request
	_outfmt string
}

func NewReqContext(w http.ResponseWriter,r *http.Request) (*ReqContext){
	return &ReqContext{w,r,""}
}

func (x ReqContext) outputData(data interface{},code int) {
	mt := x.getOutputFormat()
	x.setHeader("Content-Type", mt)
	x.w.WriteHeader(code)
	if err := outputObj(x.w, mt, data); err != nil {
		panic(err)
	}
}

func (x ReqContext) setHeader(header string,value string) {
	x.w.Header().Set(header,value)
}

func (x ReqContext) Created(loc string){
	x.setHeader("Location",loc)
	x.w.WriteHeader(http.StatusCreated)
}

func (x ReqContext) Ok(data interface{}) {
	x.outputData(data,http.StatusOK)
}

func (x ReqContext) SeeOther(url string){
	x.setHeader("Location",url)
	x.w.WriteHeader(http.StatusSeeOther)
}

func (x ReqContext) MovedPermanently(url string) {
	x.setHeader("Location",url)
	x.w.WriteHeader(http.StatusMovedPermanently)
}


func (x ReqContext) httpError(err interface{},code int) {
	x.outputData(EWrap(err),code)
}

func (x ReqContext) ServerError(err interface{}) {
	x.httpError(err,http.StatusInternalServerError)
}

func (x ReqContext) BadRequest(err interface{}) {
	x.httpError(err,http.StatusBadRequest)
}

func (x ReqContext) MethodNotAllowed(err interface{}) {
	x.httpError(err,http.StatusMethodNotAllowed)
}

func (x ReqContext) NotFound(err interface{}) {
	x.httpError(err,http.StatusNotFound)
}


func (x ReqContext) getReqHeader(header string) string{
	return x.r.Header.Get(header)
}

func (x ReqContext) reqBodyContentType() (string,error){
	return GetMediaType(x.getReqHeader("Content-Type"))
}

func (x ReqContext) Data(into interface{}) error {
	mt, err := x.reqBodyContentType()
	if err != nil {
		return err
	}
	ReadData(mt,x.r.Body,into)
	return nil
}

type validater interface {
	Validate() error
}

func (x ReqContext) ValidatedData(into validater) error {
	if err := x.Data(into); err != nil {
		return err
	}
	return into.Validate()
}

func (x ReqContext) getOutputFormat() string {
	if x._outfmt == "" {
		x._outfmt = select_mediatype(x.r.Header.Get("Accept"))
	}
	return x._outfmt
}


