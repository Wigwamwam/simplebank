package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"github.com/wigwamwam/simplebank/apic"
	db "github.com/wigwamwam/simplebank/db/sqlc"
	"github.com/wigwamwam/simplebank/util"
)

// chi framework

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	chiStore := db.NewStore(conn)
	// Why can't I sort this out?
	// apic.NewServer(chiStore)

	handler := apic.NewHandler(chiStore)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/accounts", handler.CreateAccount())
	r.Get("/accounts/{id}", handler.GetAccount())
	r.Get("/accounts", handler.IndexAccount())
	r.Put("/accounts/{id}", handler.UpdateAccount())
	r.Delete("/accounts/{id}", handler.DeleteAccount())

	http.ListenAndServe(":3000", r)

}
