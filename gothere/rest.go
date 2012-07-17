package gothere
import "net/http"

type methodHandler struct {
	Get,Head,Post,Put,Delete,Options FnHandler
}

func nyiMethodHandler() methodHandler {
	nyih := func (x *ReqContext) {
		x.MethodNotAllowed("Method not allowed (Not implemented)")
	}
	nyi := FnHandler(nyih)
	return methodHandler{nyi,nyi,nyi,nyi,nyi,nyi}
}

func NewMethodHandler() methodHandler {
	return nyiMethodHandler();
}

func RESTHandler(path string,col methodHandler,res methodHandler) {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		x := NewReqContext(w,r)

		var handler methodHandler
		if len(r.URL.Path) > len(path) {
			handler = res 
		} else {
			handler = col
		}
		
		switch r.Method {
		case "GET": handler.Get(x)
		case "HEAD": handler.Head(x)
		case "POST": handler.Post(x)
		case "PUT": handler.Put(x)
		case "DELETE": handler.Delete(x)
		case "OPTIONS": handler.Options(x)
		}
	}
	http.HandleFunc(path,wrapper)
}

func SimpleHandler(path string,fn FnHandler) {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		x := NewReqContext(w,r)
		fn(x)
	}
	http.HandleFunc(path,wrapper)
}