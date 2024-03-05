package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"snippetbox.olegmonabaka.net/internal/models"
	"snippetbox.olegmonabaka.net/internal/validator"
	"strconv"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	app.notFound(w)
	//	return
	//}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
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

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{Expires: 365}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//err := r.ParseForm()
	//if err != nil {
	//	app.clientError(w, http.StatusBadRequest)
	//	return
	//}

	var form snippetCreateForm
	err := app.decodePostForm(r, &form)

	//err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	//if err != nil {
	//	app.clientError(w, http.StatusBadRequest)
	//	return
	//}
	//
	//form := snippetCreateForm{
	//	Title:   r.PostForm.Get("title"),
	//	Content: r.PostForm.Get("content"),
	//	Expires: expires,
	//}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	//if strings.TrimSpace(form.Title) == "" {
	//	form.FieldErrors["title"] = "This field cannot be blank"
	//} else if utf8.RuneCountInString(form.Title) > 100 {
	//	form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	//}
	//
	//if strings.TrimSpace(form.Content) == "" {
	//	form.FieldErrors["content"] = "This field cannot be blank"
	//}
	//
	//if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
	//	form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	//}
	//
	//if len(form.FieldErrors) > 0 {
	//	data := app.newTemplateData(r)
	//	data.Form = form
	//	app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
	//	return
	//}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}