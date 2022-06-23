package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	router.Handle(http.MethodGet, "/", app.home)
	router.Handle(http.MethodPost, "/create", app.CreateRecord)
	router.Handle(http.MethodPut, "/update", app.UpdateRecord)
	router.Handle(http.MethodGet, "/delete/:id", app.DeleteRecord)
	router.Handle(http.MethodGet, "/records", app.GetAllRecords)
	router.Handle(http.MethodGet, "/records/:id", app.GetRecordById)

	fileServer := http.Dir("./ui")
	router.ServeFiles("/ui/*filepath", fileServer)

	_cors := cors.Options{
		AllowedMethods: []string{"*"},
		AllowedOrigins: []string{"*"},
	}
	return cors.New(_cors).Handler(app.enableCors(router))
}

func (app *application) run() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Server.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	return srv.ListenAndServe()
}

func (app *application) home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "status: Running")
}

func (app *application) enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
