package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func NewApiServer(g GoImportServicer) *ApiServer {
	return &ApiServer{g}
}

type ApiServer struct {
	giService GoImportServicer
}

func (s *ApiServer) HomeHandler(w http.ResponseWriter, r *http.Request) {

	index, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		log.Println("cant read home page:", err)
	}

	fmt.Fprintf(w, string(index))
}

func (s *ApiServer) ImportsHandler(w http.ResponseWriter, r *http.Request) {
	gis, _ := s.giService.FindAll()

	json.NewEncoder(w).Encode(gis)
}

func (s *ApiServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := s.giService.Fetch()
		if err != nil {
			//fmt.Fprint(w, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//fmt.Fprint(w, "Started")
		w.WriteHeader(http.StatusOK)
	}
}

func (s *ApiServer) Run() {

	http.HandleFunc("/", s.HomeHandler)
	http.HandleFunc("/api/fetch", s.UpdateHandler)
	http.HandleFunc("/api/imports", s.ImportsHandler)

	http.ListenAndServe(":8080", nil)

}
