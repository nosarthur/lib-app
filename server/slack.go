package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func (app *application) Slack(w http.ResponseWriter, req *http.Request) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}
	data := req.Form
	if data["token"][0] != os.Getenv("Token") {
		return fmt.Errorf("Authentication failed.")
	}
	log.Println(data)

	w.Write([]byte("Success!"))
	return nil
}
