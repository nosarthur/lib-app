package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nosarthur/todoslacker/storage"
)

// handlers in this file are all of the type
// func(http.ResponseWriter, *http.Request) error
// without ServeHTTP method

/*
	Data handles request of http:get::/data

	Example response:

	{ "tickets" :[ {"id":"grocery",
		              "detail":"vegi",
					  "todos": []
					  "start_time": 2016-12-27T05:30:34.645428Z",
					  "end_time":null,
					  "priority":false}
		] }
*/
func (app *application) Data(w http.ResponseWriter, req *http.Request) error {
	tickets, err := app.db.All()
	if err != nil {
		return err
	}
	for i, t := range tickets {
		todos, err := app.db.ReadTodos(t.Id)
		if err != nil {
			return err
		}
		tickets[i].Todos = make([]*storage.Todo, len(todos))
		copy(tickets[i].Todos, todos)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reply := map[string]interface{}{
		"tickets":     tickets,
		"last_update": app.db.LastUpdate.String(),
	}
	if err := json.NewEncoder(w).Encode(reply); err != nil {
		return err
	}
	return nil
}

// AddTicket handles request of http:post::/ticket/add
func (app *application) AddTicket(w http.ResponseWriter, req *http.Request) error {
	t := storage.Ticket{StartTime: time.Now()}
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return err
	}
	if err := app.db.CreateTicket(t); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

// EndTicket handles request of http:post::/ticket/end
func (app *application) EndTicket(w http.ResponseWriter, req *http.Request) error {
	t := storage.Ticket{}
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return err
	}
	t, err := app.db.ReadTicket(t.Id)
	errMsg := fmt.Sprintf("Cannot end Ticket with id=%v.", t.Id)
	if err != nil {
		return fmt.Errorf("%v %v", errMsg, err)
	}
	if t.EndTime != nil {
		return fmt.Errorf("%v It has ended already.", errMsg)
	}
	now := time.Now()
	if now.Before(t.StartTime) {
		return fmt.Errorf("%v Causality broken.", errMsg)
	}
	t.EndTime = &now
	if err = app.db.UpdateTicket(t); err != nil {
		return fmt.Errorf("%v %v", errMsg, err)
	}
	w.WriteHeader(http.StatusAccepted)
	return nil
}

// Addtodo handles request of http:post::/todo/add
func (app *application) AddTodo(w http.ResponseWriter, req *http.Request) error {
	t := storage.Todo{}
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return err
	}
	if err := app.db.CreateTodo(t); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

// EndTodo handles request of http:post::/todo/end
func (app *application) EndTodo(w http.ResponseWriter, req *http.Request) error {
	t, err := app.verifyTodoReq(req)
	if err != nil {
		return err
	}
	t, err = app.db.ReadTodo(t.TicketId, t.Idx)
	if err != nil {
		return err
	}
	t.Done = true
	if err = app.db.UpdateTodo(t); err != nil {
		return err
	}
	w.WriteHeader(http.StatusAccepted)
	return nil
}

// verifyTodoReq checks if the request contains a valid Todo
func (app *application) verifyTodoReq(req *http.Request) (storage.Todo, error) {
	t := storage.Todo{}
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return t, err
	}
	if _, err := app.db.ReadTicket(t.TicketId); err != nil {
		return t, err
	}
	return t, nil
}

// DeleteTodo handles request of http:delete::/todo/delete
func (app *application) DeleteTodo(w http.ResponseWriter, req *http.Request) error {
	t, err := app.verifyTodoReq(req)
	if err != nil {
		return err
	}
	if err := app.db.DeleteTodo(t); err != nil {
		return err
	}
	w.WriteHeader(http.StatusAccepted)
	return nil
}

// DeleteTicket handles request of http:delete::/ticket/delete
func (app *application) DeleteTicket(w http.ResponseWriter, req *http.Request) error {
	t := storage.Ticket{}
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return err
	}
	if err := app.db.DeleteTicket(t.Id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusAccepted)
	return nil
}
