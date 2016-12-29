package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Slack handles request of http:post::/slack
// It is a middleware for the urls that require authentication
func (app *application) Slack(w http.ResponseWriter, req *http.Request) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}
	if req.FormValue("token") != os.Getenv("Token") {
		return fmt.Errorf("Authentication failed.")
	}
	log.Println(req.FormValue("token"))

	url := req.FormValue("command")
	log.Println(url)

	reader := str2reader(req.FormValue("text"))
	req, err = http.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Token", req.FormValue("token"))
	//resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	//defer resp.Body.Close()

	w.Write([]byte("Success!"))
	return nil
}

func str2reader(s string) *strings.Reader {
	reader := strings.NewReader(`{"id": "test1", "detail": "test1"}`)
	return reader
}
