package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/anirudhp26/simpleurl-server/routes"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type handler struct {
	redisClient    *redis.Client
	postgresClient *pgxpool.Pool
}

type apiFunction func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandler(fn apiFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func NewHandler() *handler {
	rdb := NewRedisClient()
	pdb := NewDB()
	return &handler{
		redisClient:    rdb.conn,
		postgresClient: pdb.conn,
	}
}

func (h *handler) InitRoutes(router *mux.Router) {
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	router.HandleFunc("/users/create", makeHTTPHandler(h.CreateUser)).Methods("POST")
	router.HandleFunc("/users/{id}", makeHTTPHandler(h.GetUser)).Methods("GET")
	router.HandleFunc("/users/{id}/update", makeHTTPHandler(h.UpdateUser)).Methods("POST")
	router.HandleFunc("/users/{id}/delete", makeHTTPHandler(h.DeleteUser)).Methods("POST")

	router.HandleFunc("/links/create", makeHTTPHandler(h.CreateLink)).Methods("POST")
	router.HandleFunc("/links/{id}", makeHTTPHandler(h.GetLink)).Methods("GET")
	router.HandleFunc("/links/{id}/update", makeHTTPHandler(h.UpdateLink)).Methods("POST")
	router.HandleFunc("/links/{id}/delete", makeHTTPHandler(h.DeleteLink)).Methods("POST")
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var requestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Provider string `json:"provider"`
	}

	if err := ReadJson(r, &requestBody); err != nil {
		return err
	}

	username := requestBody.Username
	email := requestBody.Email
	provider := requestBody.Provider

	if username == "" || email == "" || provider == "" {
		return errors.New("missing required field(s): username, email, or provider")
	}
	log.Print(username, email, provider)
	user, err := routes.CreateUser(h.postgresClient, h.redisClient, r.Context(), map[string]string{
		"username": username,
		"email":    email,
		"provider": provider,
	})
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, user)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	user, err := routes.GetUser(h.postgresClient, h.redisClient, r.Context(), idInt)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, user)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	// Update user
	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	// Delete user
	return nil
}

func (h *handler) CreateLink(w http.ResponseWriter, r *http.Request) error {
	// Create link
	return nil
}

func (h *handler) GetLink(w http.ResponseWriter, r *http.Request) error {
	// Get link
	return nil
}

func (h *handler) UpdateLink(w http.ResponseWriter, r *http.Request) error {
	// Update link
	return nil
}

func (h *handler) DeleteLink(w http.ResponseWriter, r *http.Request) error {
	// Delete link
	return nil
}
