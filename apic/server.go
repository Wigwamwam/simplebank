package apic

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	db "github.com/wigwamwam/simplebank/db/sqlc"
)

type Handler struct {
	store    db.Store
	validate *validator.Validate
}

func NewServer(store db.Store) {

	handler := &Handler{store: store}
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/accounts/{id}", handler.getAccount())
	r.Post("/accounts", handler.createAccount())

	http.ListenAndServe(":3000", r)

}
