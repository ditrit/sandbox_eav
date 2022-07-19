package main

import (
	"flag"
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
	var port = flag.Int("p", 9999, "http port")
	flag.Parse()
	// were using sqlite now since it's easy to use but we will move to cockroachDb as soon as possible
	sqlit := sqlite.Open("test.db")
	db, err := gorm.Open(sqlit, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Connected to database")
	eav.Init(db) // drop and automigrate DB
	log.Println("Automigrate finished")

	// This function create the test schema
	eav.PopulateDatabase(db)

	router := mux.NewRouter()
	router.Use(endpoints.MiddlewareLogger)
	// Get whole collection
	router.HandleFunc("/v1/objects/{type}", endpoints.GetObjects(db)).Methods("GET")
	// 405 method not allowed
	router.HandleFunc("/v1/objects/{type}", MethodNotAllowed).Methods("PUT", "DELETE")

	//CRUD
	router.HandleFunc("/v1/objects/{type}/{id}", endpoints.GetObject(db)).Methods("GET")
	router.HandleFunc("/v1/objects/{type}/{id}", endpoints.DeleteObject(db)).Methods("DELETE")
	router.HandleFunc("/v1/objects/{type}", endpoints.CreateObject(db)).Methods("POST")
	router.HandleFunc("/v1/objects/{type}/{id}", endpoints.ModifyObject(db)).Methods("PUT")
	// It may be a good idea to choose the CORS options at the bare minimum level
	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Connection", "Host", "Origin", "User-Agent", "Referer", "Cache-Control", "X-header"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(0),
	)(router)
	fmt.Printf("Ready to handle requests at :%d !\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), cors)
}

// 405 method not allowed handler
// https://www.restapitutorial.com/lessons/httpmethods.html
func MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
