package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// serve static assets
	fileServer := http.FileServer(http.Dir(app.config.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// regular routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return secureHeaders(mux)
}
