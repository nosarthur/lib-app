package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Slack handles request of http:post::/slack
// It parses the slack request for authentication, routing, and assembles a new request to replay to the next API URL
func (app *application) Slack(w http.ResponseWriter, req *http.Request) error {
	nextURL := req.FormValue("command")
	if handler, ok := app.slackRoutes[nextURL]; ok {
		// create new request and relay to the next url
		reader := str2reader(req.FormValue("text"))
		newReq, err := http.NewRequest("POST", nextURL, reader)
		if err != nil {
			return err
		}
		newReq.Header.Set("Token", os.Getenv(req.FormValue("token")))
		err = handler(w, newReq)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unknown URL from slack: " + nextURL)
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
	msg = strings.Replace(msg, ", ", `", "`, -1)
	reader := &strings.Reader{}
	if msg != "" {
		reader = strings.NewReader(msg)
	}
	return reader
}
