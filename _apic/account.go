package apic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/wigwamwam/simplebank/context"
	db "github.com/wigwamwam/simplebank/db/sqlc"
)

// what does this have to be here and not in the routes section?

type Handler struct {
	Store db.Store
}

func NewHandler(store db.Store) Handler {
	return Handler{Store: store}
}

type createAccountRequest struct {
	Owner    string `json:"owner" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

func (handler *Handler) CreateAccount() http.HandlerFunc {
	var validate *validator.Validate
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithHTTPRequest(r.Context(), r)

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			handleAppError(w, err)
			return
		}

		defer r.Body.Close()

		var req createAccountRequest
		ctx = context.WithHTTPRequestBody(ctx, bytes)
		err = json.Unmarshal(bytes, &req)
		if err != nil {
			handleAppError(w, err)
			return
		}

		validate = validator.New()
		err = validate.Struct(req)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			handleAppError(w, errs)
			return
		}

		arg := db.CreateAccountParams{
			Owner:    req.Owner,
			Currency: req.Currency,
			Balance:  0,
		}

		account, err := handler.Store.CreateAccount(ctx, arg)
		if err != nil {
			handleAppError(w, err)
			return
		}

		js, err := json.Marshal(account)
		if err != nil {
			handleAppError(w, err)
			return
		}

		respondWithJSON(w, http.StatusOK, js)
	}
}

func (handler *Handler) GetAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		fmt.Println("here ", idStr)
		if idStr == "" {
			err := errors.New("missing or empty 'id' parameter")
			handleAppError(w, err)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			err = fmt.Errorf("invalid 'id' parameter '%s': %s", idStr, err)
			handleAppError(w, err)
			return
		}

		ctx := context.WithHTTPRequest(r.Context(), r)
		account, err := handler.Store.GetAccount(ctx, int64(id))
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

type ListAccountsRequest struct {
	Limit  int32 `form:"limit" validate:"required,min=1"`
	Offset int32 `form:"offset" validate:"required,min=5"`
}

func (handler *Handler) IndexAccount() http.HandlerFunc {
	var validate *validator.Validate
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithHTTPRequest(r.Context(), r)

		limitParam := r.URL.Query().Get("limit")
		offsetParam := r.URL.Query().Get("offset")

		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			handleAppError(w, err)
		}

		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			handleAppError(w, err)
		}

		// validate:
		// not sure how to shorten this:
		validate = validator.New()
		err = validate.Struct(
			ListAccountsRequest{
				Limit:  int32(limit),
				Offset: int32(offset),
			},
		)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			handleAppError(w, errs)
			return
		}

		arg := db.ListAccountsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		allAccounts, err := handler.Store.ListAccounts(ctx, arg)
		if err != nil {
			handleAppError(w, err)
			return
		}

		js, err := json.Marshal(allAccounts)
		if err != nil {
			handleAppError(w, err)
			return
		}

		respondWithJSON(w, http.StatusOK, js)
	}
}

func (handler *Handler) DeleteAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithHTTPRequest(r.Context(), r)
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			err := errors.New("missing or empty 'id' parameter")
			handleAppError(w, err)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			err = fmt.Errorf("invalid 'id' parameter '%s': %s", idStr, err)
			handleAppError(w, err)
			return
		}

		err = handler.Store.DeleteAccount(ctx, int64(id))
		if err != nil {
			handleAppError(w, err)
			return
		}

		respondWithJSON(w, http.StatusOK, nil)
	}
}

type updateAccountRequest struct {
	Balance int `json:"balance"`
}

func (handler *Handler) UpdateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithHTTPRequest(r.Context(), r)
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			err := errors.New("Mising or empty 'id' paraemeter")
			handleAppError(w, err)
			return
		}

		var req updateAccountRequest

		id, err := strconv.Atoi(idStr)
		if err != nil {
			err = fmt.Errorf("invalid 'id' parameter '%s': %s", idStr, err)
			handleAppError(w, err)
			return
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			handleAppError(w, err)
			return
		}

		err = json.Unmarshal(bytes, &req)
		if err != nil {
			handleAppError(w, err)
			return
		}

		fmt.Println(req)

		defer r.Body.Close()

		arg := db.UpdateAccountParams{
			ID: int64(id),
			Balance: int64(req.Balance),
		}

		updatedAccount, err := handler.Store.UpdateAccount(ctx, arg)
		if err != nil {
			handleAppError(w, err)
			return
		}

		js, err := json.Marshal(updatedAccount)
		if err != nil {
			handleAppError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(js)
		if err != nil {
			handleAppError(w, err)
			return
		}
	}
}
