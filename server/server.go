package server

import (
	"encoding/json"
	c "github.com/trenker/boxserver/conf"
	"github.com/trenker/boxserver/log"
	"github.com/trenker/boxserver/util"
	"net/http"
	"strings"
)

type request struct {
	content     interface{}
	omitContent bool
	status      int
	parts       []string
}

func (r *request) Process(res http.ResponseWriter) {

	h := res.Header()
	h["Content-Type"] = []string{"application/json;charset=UTF-8"}
	cors := c.Get().Cors

	if cors != "" {
		h["Access-Control-Allow-Origin"] = []string{cors}
	}

	res.WriteHeader(r.status)

	if !r.omitContent {
		body, _ := json.MarshalIndent(r.content, "", "  ")
		res.Write(body)
	}
}

func NewRequest(req *http.Request) *request {

	path := strings.Split(strings.Trim(strings.TrimPrefix(req.URL.Path, c.Get().Proxy), "/"), "/")

	var content interface{} = util.Str("Unsupported request type")
	status := http.StatusBadRequest

	log.Debug("Requested '%s' %s", req.Method, path)

	if req.Method == "PUT" {
		log.Debug("Handle PUT")
		content, status = handlePut(path, req)
	}

	if req.Method == "DELETE" {
		log.Debug("Handle DELETE")
		content, status = handleDelete(path)
	}

	if req.Method == "GET" || req.Method == "HEAD" || req.Method == "" {
		log.Debug("Handle GET")
		content, status = handleGet(path)
	}

	return &request{
		omitContent: req.Method == "HEAD",
		status:      status,
		content:     content,
	}
}
