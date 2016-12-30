package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Slack handles request of http:post::/slack
// It parses the slack request for authentication, routing, and assembles the request for other API
func (app *application) Slack(w http.ResponseWriter, req *http.Request) error {
	if req.FormValue("token") != os.Getenv("Token") {
		return fmt.Errorf("Authentication failed.")
	}
	nextURL := req.FormValue("command")
	reader := str2reader(req.FormValue("text"))
	// create new request and relay to the next url
	newReq, err := http.NewRequest("POST", nextURL, reader)
	if err != nil {
		return err
	}
	newReq.Header.Set("Token", req.FormValue("token"))

	switch nextURL {
	case "/ticket/add":
		err = app.AddTicket(w, newReq)
	case "/ticket/end": // needs to be fixed
		req.Method = "DELETE"
		err = app.EndTicket(w, newReq)
	case "/todo/add":
		err = app.AddTodo(w, newReq)
	case "/todo/end": // needs to be fixed
		req.Method = "DELETE"
		err = app.EndTodo(w, newReq)
	default:
		err = fmt.Errorf("Unknown url from slack: " + nextURL)
	}
	if err != nil {
		return err
	}
	w.Write([]byte("Success!"))
	return nil
}

/* 	Str2reader converts slack msg to reader.
	Input string should be of the form `id:test1 detail:test1`.
 	It is then converted to JSON string `{"id": "test1", "detail": "test1"}`
*/
func str2reader(msg string) *strings.Reader {
	msg = `{"` + msg + `"}`
	msg = strings.Replace(msg, ":", `":"`, -1)
	msg = strings.Replace(msg, " ", `", "`, -1)
	reader := &strings.Reader{}
	if msg != "" {
		reader = strings.NewReader(msg)
	}
	return reader
}
