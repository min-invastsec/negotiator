package negotiator

import "net/http"

var processors = []ResponseProcessor{&JSONProcessor{}, &XMLProcessor{}}

//New sets up response processors. By default XML and JSON are created
func New(responseProcessors ...*ResponseProcessor) {
	for _, proc := range responseProcessors {
		processors = append(processors, *proc)
	}
}

//Negotiate your model based on HTTP Accept header
func Negotiate(w http.ResponseWriter, req *http.Request, model interface{}) {

	accept := new(Accept)
	//TODO:test should not be case sensitive
	accept.Header = req.Header.Get("Accept")

	for _, mr := range accept.MediaRanges() {
		for _, processor := range processors {
			if processor.CanProcess(mr.Value) {
				processor.Process(w, model)
				return
			}
		}
	}

	//rfc2616-sec14.1
	//If an Accept header field is present, and if the
	//server cannot send a response which is acceptable according to the combined
	//Accept field value, then the server SHOULD send a 406 (not acceptable)
	//response.
	http.Error(w, "", http.StatusNotAcceptable)
}
