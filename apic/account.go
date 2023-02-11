package apic

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/wigwamwam/simplebank/context"
	db "github.com/wigwamwam/simplebank/db/sqlc"
)

// context

type createAccountRequest struct {
	Owner    string `json:"owner" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

func (handler *Handler) createAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createAccountRequest

		ctx := context.WithHTTPRequest(r.Context(), r)

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			handleAppError(w, err)
			return
		}
		defer r.Body.Close()

		ctx = context.WithHTTPRequestBody(ctx, bytes)

		err = json.Unmarshal(bytes, &req)
		if err != nil {
			handleAppError(w, err)
			return
		}

		handler.validate = validator.New()
		err = handler.validate.Struct(req)
		if err != nil {
			handleAppError(w, err)
			return
		}

		arg := db.CreateAccountParams{
			Owner:    req.Owner,
			Currency: req.Currency,
			Balance:  0,
		}

		account, err := handler.store.CreateAccount(ctx, arg)
		if err != nil {
			handleAppError(w, err)
			return
		}

		js, err := json.Marshal(account)
		if err != nil {
			handleAppError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(js)
		if err != nil {
			handleAppError(w, err)
			return
		}

	}
}

type getAccountRequest struct {
	Id int64 `uri:"id" validate:"required,min=1"`
}

func (handler *Handler) getAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithHTTPRequest(r.Context(), r)
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			handleAppError(w, err)
		}
		

		// handler.validate = validator.New()
		// err = handler.validate.Struct(req)
		// if err != nil {
		// 	handleAppError(w, err)
		// 	return
		// }

		account, err := handler.store.GetAccount(ctx, int64(id))

		if err != nil {
			handleAppError(w, err)
			return
		}

		js, err := json.Marshal(account)
		if err != nil {
			handleAppError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(js)
		if err != nil {
			handleAppError(w, err)
			return
		}
	}
}
