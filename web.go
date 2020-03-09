package main

import (
	"html/template"
	"net/http"
)

func NewWebServer(gisDB GoImportDB, grsDB GithubRepoDB) *WebServer {
	return &WebServer{gisDB, grsDB}
}

type WebServer struct {
	GisDB GoImportDB
	GrsDB GithubRepoDB
}

func (ws *WebServer) GoImportsHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("imports.html")

	gis := ws.GisDB.GetAll()

	t.Execute(w, gis)
}

func (ws *WebServer) UpdateGithubHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		go func() {
			grs := GetGithubGoRepos()

			ws.GrsDB.SaveAll(grs)

			gis := ExtractGoImports(grs)

			ws.GisDB.SaveAll(gis)

		}()

	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (ws *WebServer) Run() {

	http.HandleFunc("/", ws.GoImportsHandler)
	http.HandleFunc("/update/github", ws.UpdateGithubHandler)

	http.ListenAndServe(":8080", nil)

}
