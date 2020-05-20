package topgomods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func NewApiServer(g GoModuleServicer) *ApiServer {
	return &ApiServer{g}
}

type ApiServer struct {
	gmService GoModuleServicer
}

func (s *ApiServer) HomeHandler(w http.ResponseWriter, r *http.Request) {

	index, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		log.Println("home handler error:", err)
	}

	fmt.Fprintf(w, string(index))
}

func (s *ApiServer) ModulesHandler(w http.ResponseWriter, r *http.Request) {
	moduleRanks, err := s.gmService.TopModules()
	if err != nil {
		log.Println("modules hadler error:", err)
	}

	json.NewEncoder(w).Encode(moduleRanks)
}

func (s *ApiServer) FetchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := s.gmService.Fetch()
		if err != nil {
			//fmt.Fprint(w, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		//fmt.Fprint(w, "Started")
		w.WriteHeader(http.StatusOK)
	}
}

func (s *ApiServer) MetaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"version\": \"%s\"}", VERSION)
}

func (s *ApiServer) Run() {

	http.HandleFunc("/", s.HomeHandler)
	http.HandleFunc("/api/fetch", s.FetchHandler)
	http.HandleFunc("/api/modules", s.ModulesHandler)

	http.HandleFunc("/api/meta", s.MetaHandler)

	http.ListenAndServe(":8080", nil)

}
