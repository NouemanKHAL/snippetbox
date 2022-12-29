package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/NouemanKHAL/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

const (
	homeTemplateFilename   = "home.go.tmpl"
	viewTemplateFilename   = "view.go.tmpl"
	createTemplateFilename = "create.go.tmpl"
)

type snippetCreateForm struct {
	Title            string
	Content          string
	Expires          int
	ValidationErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, homeTemplateFilename, data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, viewTemplateFilename, data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, 200, createTemplateFilename, data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:            r.PostForm.Get("title"),
		Content:          r.PostForm.Get("content"),
		Expires:          expires,
		ValidationErrors: make(map[string]string),
	}

	if form.Title == "" {
		form.ValidationErrors["title"] = "this field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.ValidationErrors["title"] = "this field cannot be more than 100 characters long"
	}

	if form.Content == "" {
		form.ValidationErrors["content"] = "this field cannot be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.ValidationErrors["expires"] = "this field must equal 1, 7 or 365"
	}

	if len(form.ValidationErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, createTemplateFilename, data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
