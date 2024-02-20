package user

import (
	"backend/auth"
	"backend/db"
	"backend/endpoint"
	server_error "backend/error"
	"backend/validator"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
)

// API represents a struct for the user API
type API struct {
	q         *db.Queries
	validator *validator.Validate
}

// NewAPI initializes an API struct
func NewAPI(conn *db.DB, v *validator.Validate) API {
	q := db.New(conn)
	return API{
		q:         q,
		validator: v,
	}
}

// type GetUserInput struct {
// 	ID   `json:"id" validate:"required"`
// }

//HandleGetUser returns a user by ID

func (api *API) HandleGetUser(w http.ResponseWriter, r *http.Request) {


	// // var input GetUserInput
	// // err := endpoint.DecodeAndValidateJson(w, r, api.validator, &input)
	// if err != nil {
	// 	return
	// }

	id := chi.URLParam(r, "id");

	result, err:= strconv.ParseInt(id, 10, 32)
	var id1 int32

	if err!= nil {
		endpoint.WriteWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else{
		id1 = int32(result)
	}

	user, err := api.q.GetUser(r.Context(), id1)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			endpoint.WriteWithError(w, http.StatusNotFound, server_error.ErrUserNotFound)
			return
		}
		endpoint.WriteWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	endpoint.WriteWithStatus(w, http.StatusOK, user)
}

func (api *API) HandleGetUsers(w http.ResponseWriter, r *http.Request) {

	user, err := api.q.GetUsers(r.Context())
	if err != nil {
		endpoint.WriteWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	endpoint.WriteWithStatus(w, http.StatusOK, user)
}




func (api *API) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		endpoint.WriteWithError(w, http.StatusUnauthorized, server_error.ErrUnauthorized)
		return
	}

	user, err := api.q.GetUser(r.Context(), int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			endpoint.WriteWithError(w, http.StatusNotFound, server_error.ErrUserNotFound)
			return
		}
		endpoint.WriteWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	endpoint.WriteWithStatus(w, http.StatusOK, user)
}

func (api *API) RegisterHandlers(r chi.Router, auth_guard func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(auth_guard)
			r.Get("/me", api.HandleGetMe)
		})
		r.Get("/", api.HandleGetUsers)
		r.Get("/{id}", api.HandleGetUser)
	})
}
