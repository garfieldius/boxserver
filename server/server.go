package server

import (
	c "github.com/trenker/boxserver/conf"
	"github.com/trenker/boxserver/log"
	"net/http"
	"strings"
	"errors"
	"encoding/json"
)

type request struct {
	content interface {}
	omitContent bool
	status int
	parts []string
}

type Message struct {
	Message error `json:"msg,string"`
}

var errorMessage []byte

func init() {
	errorMessage, _ = json.Marshal(&Message{Message: errors.New("Unknown resource")})
}

func (r *request) Process(res http.ResponseWriter) {

	h := res.Header()
	h["Content-Type"] = []string{"application/json;charset=UTF-8"}
	res.WriteHeader(r.status)

	if !r.omitContent && r.content != nil {
		body, _ := json.MarshalIndent(r.content, "", "  ")
		res.Write(body)
	}
}

func NewRequest(req *http.Request) *request {
	hasContent := req.Method == "HEAD"
	path := strings.Split(strings.Trim(strings.TrimPrefix(req.URL.Path, c.Get().Proxy), "/"), "/")

	var content interface {}
	var status int = http.StatusNotFound

	log.Debug("Requested %s %s", req.Method, path)

	switch req.Method {
		default:
			content, status = handleGet(path)

	}

	var r request

	if content != nil {
		r = request{
			omitContent: hasContent,
			status: status,
			content: content,
		}
	} else {
		r = request{
			omitContent: hasContent,
			status: http.StatusBadRequest,
			content: errorMessage,
		}
	}

	return &r
}
