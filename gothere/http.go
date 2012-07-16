// My tiny kit of http extras for Go
package gothere
import ("net/http")

type FnHandler func(*ReqContext)

func HandleMethods(path string,fnget FnHandler,fnpost FnHandler) {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		x := NewReqContext(w,r)
		switch r.Method {
		case "GET":
			fnget(x)
		case "POST":
			fnpost(x)
		}
	}
	http.HandleFunc(path,wrapper)
}


func Output(w http.ResponseWriter, mediatype string, obj interface{}) error {
	w.Header().Set("Content-Type", mediatype)
	return output_obj(w, mediatype, obj)
}



type ReqContext struct {
	w http.ResponseWriter
	r *http.Request
	_outfmt string
}

func NewReqContext(w http.ResponseWriter,r *http.Request) (*ReqContext){
	return &ReqContext{w,r,""}
}

func (x ReqContext) Created(loc string){
	x.w.Header().Set("Location",loc)
	x.w.WriteHeader(http.StatusCreated)
}

func (x ReqContext) outputData(data interface{},code int) {
	mt := x.getOutputFormat()
	x.w.Header().Set("Content-Type", mt)
	x.w.WriteHeader(code)
	if err := output_obj(x.w, mt, data); err != nil {
		E500(x.w,err)
	}
}

func (x ReqContext) Ok(data interface{}) {
	x.outputData(data,http.StatusOK)
}


func (x ReqContext) ServerError(err interface{}) {
	x.outputData(EWrap(err),500)
}

func (x ReqContext) BadRequest(err interface{}) {
	x.outputData(EWrap(err),400)
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


func (x ReqContext) getOutputFormat() string {
	if x._outfmt == "" {
		x._outfmt = select_mediatype(x.r.Header.Get("Accept"))
	}
	return x._outfmt
}