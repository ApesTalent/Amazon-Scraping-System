package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"primeprice.com/controller"

	_ "github.com/joho/godotenv/autoload"
	"primeprice.com/config"
	"primeprice.com/dal"
	"primeprice.com/pkg/logger"
	r "primeprice.com/router"
)

func main() {
	config.LoadConfig("config/config.json")
	dal.LoadDB()
	dal.MigrateDB()
	dal.CreateAdmin()
	logger.InitLogger()
	rand.Seed(time.Now().UnixNano())

	controller.StartsAllScraping()

	router := r.GetRouter()
	s := &http.Server{
		Addr:           r.GetPort(),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.SetKeepAlivesEnabled(false)
	log.Printf("Listening on port %s", r.GetPort())
	log.Fatal(s.ListenAndServe())
}
