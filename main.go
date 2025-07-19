package main

import (
	"log"
	"net/http"
	"fitgenie/internal/api"
)

func main() {
	router := api.SetupRoutes()
	log.Println("Servidor FitGenie escuchando en :8080")
	http.ListenAndServe(":8080", router)
}
