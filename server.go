package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ditrit/sandbox_eav/eav"
	"github.com/ditrit/sandbox_eav/endpoints"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// were using sqlite now since it's easy to use but we will move to cockroachDb as soon as possible
	sqlit := sqlite.Open("test.db")
	db, err := gorm.Open(sqlit, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Connected to database")
	eav.Init(db)
	log.Println("Automigrate finished")
	eav.PopulateDatabase(db)

	router := mux.NewRouter()
	router.Use(endpoints.MiddlewareLogger)
	router.HandleFunc("/v1/object/{type}/{id}", endpoints.GetObject(db)).Methods("GET")
	router.HandleFunc("/v1/object/{type}/{id}", endpoints.DeleteObject(db)).Methods("DELETE")
	router.HandleFunc("/v1/object/{type}", endpoints.CreateObject(db)).Methods("POST")
	router.HandleFunc("/v1/object/{type}/{id}", endpoints.ModifyObject(db)).Methods("PUT")
	// It may be a good idea to choose the CORS options at the bare minimum level
	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Connection", "Host", "Origin", "User-Agent", "Referer", "Cache-Control", "X-header"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(0),
	)(router)
	fmt.Println("Ready !")
	http.ListenAndServe(":9999", cors)
}
