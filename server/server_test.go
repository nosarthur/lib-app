package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nosarthur/todobot/storage"
	"github.com/stretchr/testify/assert"
)

var (
	app    *application
	server *httptest.Server
	reader io.Reader
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

func TestServer(t *testing.T) {
	setup()
	t.Run("AddTicket", func(t *testing.T) {
		url := server.URL + "/ticket/add"
		reader = strings.NewReader(`{"id": "test1", "detail": "test1"}`)
		req, err := http.NewRequest("POST", url, reader)
		assert.Nil(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		// add a ticket with same id
		reader = strings.NewReader(`{"id": "test1"}`)
		req, err = http.NewRequest("POST", url, reader)
		assert.Nil(t, err)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		// add another ticket
		reader = strings.NewReader(`{"id": "test2", "detail": "test2"}`)
		req, err = http.NewRequest("POST", url, reader)
		assert.Nil(t, err)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
	t.Run("EndTicket", func(t *testing.T) {
		url := server.URL + "/ticket/end/test1"
		req, err := http.NewRequest("DELETE", url, nil)
		assert.Nil(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)
		// end a non-existing ticket
		url = server.URL + "/ticket/end/test100"
		req, err = http.NewRequest("DELETE", url, nil)
		assert.Nil(t, err)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("AddTodo", func(t *testing.T) {
		url := server.URL + "/todo/add"
		reader = strings.NewReader(`{"ticket_id": "test1", "item": "todo1"}`)
		req, err := http.NewRequest("POST", url, reader)
		assert.Nil(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		// add it again
		reader = strings.NewReader(`{"ticket_id": "test1", "item": "todo1"}`)
		req, err = http.NewRequest("POST", url, reader)
		assert.Nil(t, err)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		// add another todo
		reader = strings.NewReader(`{"ticket_id": "test2"}`)
		req, err = http.NewRequest("POST", url, reader)
		assert.Nil(t, err)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		// add a todo with invalid ticket_id
		reader = strings.NewReader(`{"ticket_id": "test200"}`)
		req, err = http.NewRequest("POST", url, reader)
		assert.Nil(t, err)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("EndTodo", func(t *testing.T) {
		url := server.URL + "/todo/end/test1/1"
		req, err := http.NewRequest("DELETE", url, nil)
		assert.Nil(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)
		// end a non-existing todo
		url = server.URL + "/todo/end/test100/1"
		req, err = http.NewRequest("DELETE", url, nil)
		assert.Nil(t, err)
		resp, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("Get", func(t *testing.T) {
		url := server.URL + "/data"
		resp, err := http.DefaultClient.Get(url)
		assert.Nil(t, err)
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
