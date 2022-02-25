package main

import (
	"bytes"
	"fmt"
	"net/http"

	// "os"
	"strconv"

	"time"

	"alexedwards.net/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// // for _, snippet := range s {
	// // 	fmt.Fprintf(w, "%v\n", snippet)
	// // }
	// data := &templateData{Snippets: s}

	// files := []string{
	// 	"/home/surdyhey/code/snippetbox/ui/html/home.page.tmpl",
	// 	"/home/surdyhey/code/snippetbox/ui/html/base.layout.tmpl",
	// 	"/home/surdyhey/code/snippetbox/ui/html/footer.partial.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// app.errorLog.Println(err.Error())
	// 	// http.Error(w, "Internal Server Error", 500)
	// 	app.serverError(w, err)
	// 	return
	// }

	// // We then use the Execute() method on the template set to write the templa
	// // content as the response body. The last parameter to Execute() represents
	// // dynamic data that we want to pass in, which for now we'll leave as nil.
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.errorLog.Println(err.Error())
	// 	http.Error(w, "Internal Server Error", 500)
	// }

	app.render(w, r, "home.page.tmpl", &templateData{Snippets: s})

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{Snippet: s})

	// data := &templateData{Snippet: s}

	// //fmt.Fprintf(w, "%v", s)
	// files := []string{
	// 	"/home/surdyhey/code/snippetbox/ui/html/show.page.tmpl",
	// 	"/home/surdyhey/code/snippetbox/ui/html/base.layout.tmpl",
	// 	"/home/surdyhey/code/snippetbox/ui/html/footer.partial.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", 405)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
	expires := "7"
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

	// 	w.Write([]byte("Create a new snippet..."))
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	} // Execute the template set, passing in any dynamic data.
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	buf.WriteTo(w)
}
