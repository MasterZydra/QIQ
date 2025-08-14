package request

import (
	"net/http"
	"strings"
	"time"
)

type Request struct {
	DocumentRoot   string
	Method         string
	QueryString    string
	Protocol       string
	RequestTime    time.Time
	RequestURI     string
	RemoteAddr     string
	RemotePort     string
	ScriptFilename string
	ServerAddr     string
	ServerPort     string
	Env            map[string]string
	Args           [][]string
	Post           string
	Cookie         string
	// Internal use
	UploadedFiles []string
}

func NewRequest() *Request {
	return &Request{
		RequestTime:   time.Now(),
		Args:          [][]string{},
		UploadedFiles: []string{},
	}
}

func NewRequestFromGoRequest(r *http.Request, documentRoot string, serverAddr string, scriptFilename string) *Request {
	request := NewRequest()
	if r != nil {
		request.DocumentRoot = documentRoot
		request.Method = r.Method
		request.QueryString = r.URL.RawQuery
		request.Cookie = strings.Join(r.Header["Cookie"], "")
		request.Protocol = r.Proto
		request.RequestURI = r.RequestURI
		request.RemoteAddr = strings.Split(r.RemoteAddr, ":")[0]
		request.RemotePort = strings.Split(r.RemoteAddr, ":")[1]
		request.ScriptFilename = scriptFilename
		request.ServerAddr = strings.Split(serverAddr, ":")[0]
		request.ServerPort = strings.Split(serverAddr, ":")[1]
	}
	// TODO HTTPS - Also in env init for $_SERVER
	// TODO ORIG_PATH_INFO - Also in env init for $_SERVER
	// TODO PATH_INFO - Also in env init for $_SERVER
	// TODO PATH_TRANSLATED - Also in env init for $_SERVER
	// TODO REDIRECT_REMOTE_USER - Also in env init for $_SERVER
	// TODO REMOTE_HOST - Also in env init for $_SERVER
	// TODO REMOTE_USER - Also in env init for $_SERVER
	// TODO SERVER_ADMIN - Also in env init for $_SERVER
	// TODO SERVER_NAME - Also in env init for $_SERVER
	return request
}
