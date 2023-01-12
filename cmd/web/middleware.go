package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		authenticatedID, ok := app.sessionManager.Get(r.Context(), authenticatedUserIDKey).(string)
		if !ok {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		id, err := strconv.ParseInt(authenticatedID, 10, 0)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		exists, err := app.users.Exists(int(id))
		if err != nil {
			app.serverError(w, err)
		}

		if !exists {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "isAuthenticated", true)
		r = r.WithContext(ctx)
		// browsers should not cache pages that require auth
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)

	})
}
