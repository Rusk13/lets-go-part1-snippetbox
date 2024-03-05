package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	//mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.PathPrefix("/static/*filepath").Handler(http.StripPrefix("/static", fileServer)).Methods(http.MethodGet)

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handle("/", dynamic.ThenFunc(app.home)).Methods(http.MethodGet)
	router.Handle("/snippet/view/{id}", dynamic.ThenFunc(app.snippetView)).Methods(http.MethodGet)
	router.Handle("/snippet/create", dynamic.ThenFunc(app.snippetCreate)).Methods(http.MethodGet)
	router.Handle("/snippet/create", dynamic.ThenFunc(app.snippetCreatePost)).Methods(http.MethodPost)
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet/view", app.snippetView)
	//mux.HandleFunc("/snippet/create", app.snippetCreate)
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
