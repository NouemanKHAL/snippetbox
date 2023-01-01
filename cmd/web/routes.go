package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// middleware chains
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamic := alice.New(app.sessionManager.LoadAndSave)
	auth := dynamic.Append(app.requireAuth)

	// serve static assets
	fileServer := http.FileServer(http.Dir(app.config.staticDir))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// regular routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", auth.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", auth.ThenFunc(app.snippetCreatePost))

	// auth routes
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", auth.ThenFunc(app.userLogoutPost))

	return standard.Then(router)
}
