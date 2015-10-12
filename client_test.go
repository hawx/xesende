package xesende

import (
	"io"
	"io/ioutil"
	"net/http"
)

type recordingHandler struct {
	Request     *http.Request
	RequestBody string

	body    string
	code    int
	headers map[string]string
}

func newRecordingHandler(body string, code int, headers map[string]string) *recordingHandler {
	return &recordingHandler{
		Request: new(http.Request),
		body:    body,
		code:    code,
		headers: headers,
	}
}

func (h *recordingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Request = r
	h.RequestBody = readAll(r.Body)
	defer r.Body.Close()

	for key, value := range h.headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(h.code)
	w.Write([]byte(h.body))
}

func readAll(r io.Reader) string {
	if b, err := ioutil.ReadAll(r); err == nil {
		return string(b)
	}

	return ""
}
