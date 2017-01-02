package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/nosarthur/todobot/storage"
	"github.com/stretchr/testify/assert"
)

var (
	app    *application
	server *httptest.Server
)

func setup() {
	app = NewApplication("postgresql://localhost/travis_ci_test?sslmode=disable")
	app.db.MustDropTables()
	app.db.MustCreateTables()
	router := NewRouter(app)
	server = httptest.NewServer(router)
}

func teardown() {
	app.db.MustDropTables()
	server.Close()
}

func makeRequest(t *testing.T, url string, method string, body string) *http.Request {
	reader := &strings.Reader{}
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, reader)
	assert.Nil(t, err)
	req.Header.Set("Token", os.Getenv("Token"))
	return req
}

func runRequest(t *testing.T, req *http.Request, expectedCode int) {
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, expectedCode, resp.StatusCode)
}

func TestServer(t *testing.T) {
	setup()
	t.Run("AddTicket", func(t *testing.T) {
		url := server.URL + "/ticket/add"
		req := makeRequest(t, url, "POST", `{"id": "test1", "detail": "test1"}`)
		runRequest(t, req, http.StatusCreated)
		// add a ticket with same id
		req = makeRequest(t, url, "POST", `{"id": "test1"}`)
		runRequest(t, req, http.StatusInternalServerError)
		// add another ticket
		req = makeRequest(t, url, "POST", `{"id": "test2", "detail": "test2"}`)
		runRequest(t, req, http.StatusCreated)
	})
	t.Run("EndTicket", func(t *testing.T) {
		url := server.URL + "/ticket/end"
		req := makeRequest(t, url, "POST", `{"id":"test1"}`)
		runRequest(t, req, http.StatusAccepted)
		// end a non-existing ticket
		req = makeRequest(t, url, "POST", `{"id":"test100"}`)
		runRequest(t, req, http.StatusInternalServerError)
	})
	t.Run("AddTodo", func(t *testing.T) {
		url := server.URL + "/todo/add"
		req := makeRequest(t, url, "POST", `{"ticket_id": "test1", "item": "todo1"}`)
		runRequest(t, req, http.StatusCreated)
		// add it again
		req = makeRequest(t, url, "POST", `{"ticket_id": "test1", "item": "todo1"}`)
		runRequest(t, req, http.StatusCreated)
		// add another todo
		req = makeRequest(t, url, "POST", `{"ticket_id": "test2"}`)
		runRequest(t, req, http.StatusCreated)
		// add another todo with wrong authentication token
		req = makeRequest(t, url, "POST", `{"ticket_id": "test2"}`)
		req.Header.Set("Token", "Token")
		runRequest(t, req, http.StatusInternalServerError)
		// add a todo with invalid ticket_id
		req = makeRequest(t, url, "POST", `{"ticket_id": "test200"}`)
		runRequest(t, req, http.StatusInternalServerError)
	})
	t.Run("EndTodo", func(t *testing.T) {
		url := server.URL + "/todo/end"
		req := makeRequest(t, url, "POST", `{"ticket_id":"test1", "idx":"1"}`)
		runRequest(t, req, http.StatusAccepted)
		// end a non-existing todo
		req = makeRequest(t, url, "POST", `{"ticket_id":"test100", "idx":1}`)
		runRequest(t, req, http.StatusInternalServerError)
	})
	t.Run("Data", func(t *testing.T) {
		url := server.URL + "/data"
		resp, err := http.DefaultClient.Get(url)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"))
		type data struct {
			Tickets []storage.Ticket `json:"tickets"`
		}
		reply := data{}
		err = json.NewDecoder(resp.Body).Decode(&reply)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(reply.Tickets))
		// reply.Tickets[0] is test2
		// reply.Tickets[1] is test1
		// hopefully this ordering doesn't vary on machine
		assert.Nil(t, reply.Tickets[0].EndTime)
		assert.NotNil(t, reply.Tickets[1].EndTime)
		assert.Equal(t, "test1", reply.Tickets[1].Id)
		assert.Equal(t, "test2", reply.Tickets[0].Id)
		assert.Equal(t, 2, len(reply.Tickets[1].Todos))
		assert.Equal(t, 1, len(reply.Tickets[0].Todos))
		assert.Equal(t, "", (*reply.Tickets[0].Todos[0]).Item)
		assert.Equal(t, "todo1", (*reply.Tickets[1].Todos[0]).Item)
		assert.Equal(t, "todo1", (*reply.Tickets[1].Todos[1]).Item)
		assert.Equal(t, false, (*reply.Tickets[1].Todos[0]).Done)
		assert.Equal(t, true, (*reply.Tickets[1].Todos[1]).Done)
	})
	teardown()
}

func TestStr2JSON(t *testing.T) {
	testStrs := []struct {
		in  string
		out string
	}{
		{`a:b`, `{"a":"b"}`},
		{`id:test1, detail:test1`, `{"id":"test1", "detail":"test1"}`},
		{`id:household task, detail:do laundry`, `{"id":"household task", "detail":"do laundry"}`},
	}
	for _, tt := range testStrs {
		reader := str2reader(tt.in)
		bArray, err := ioutil.ReadAll(reader)
		assert.Nil(t, err)
		str := string(bArray)
		assert.Equal(t, tt.out, str)
	}
}
