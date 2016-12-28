package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	server  *httptest.Server
	respRec *httptest.ResponseRecorder
	m       *mux.Router
	req     *http.Request
	err     error
)

func init() {
	server = httptest.NewServer(nil)
	fmt.Println(server.URL)
}
func setup() {
	m = mux.NewRouter()
	respRec = httptest.NewRecorder()
}

func TestGet(t *testing.T) {
	setup()
	req, err = http.NewRequest("GET", server.URL+"/data", nil)
	assert.Nil(t, err)
	m.ServeHTTP(respRec, req)

}

func TestAddTicket(t *testing.T) {
}
