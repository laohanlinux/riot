package router

import (
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	http.ListenAndServe("127.0.0.1:8888", NewRouter())
}
