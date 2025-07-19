package api

import (
	"net/http"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/suggest-outfit", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementar lógica de sugerencia de outfit
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented"))
	}).Methods("GET")

	router.HandleFunc("/add-clothing", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementar lógica para añadir prenda
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented"))
	}).Methods("POST")

	router.HandleFunc("/clothes", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementar lógica para listar prendas
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented"))
	}).Methods("GET")

	return router
}
