package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type APIServer struct {
	listenAddr 	string
	store		Storage
}

func NewAPIServer(listenAddr string, s Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: s,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/", makeHTTPHandlerFunc(s.handleHello))
	router.HandleFunc("/signup", makeHTTPHandlerFunc(s.handleSignup))
	router.HandleFunc("/login", makeHTTPHandlerFunc(s.handleLogin))
	router.HandleFunc("/user/{id}", makeHTTPHandlerFunc(s.handleGetUserByID))
	fmt.Println("server running at", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))
}

func (s *APIServer) handleHello(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, &APIMessage{Message: "hello"})
}

func (s *APIServer) handleSignup(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
	var req CreateUserReq
	// err := json.NewDecoder(r.Body).Decode(&user)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "wrong format"})
	}
	user := NewUser(req.Username, req.Password, req.UserProfile)
	id, err := s.store.CreateUser(user)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "failed to create user"})
	}
	user.ID = id
	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "invalid credentials"})
	}
	user, err := s.store.GetUserByUsername(req.Username)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "invalid credentials"})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.EncPassword), []byte(req.Password))
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "invalid credentials"})
	}

	return WriteJSON(w, http.StatusOK, &APIMessage{Message: "login success"})
}

func (s *APIServer) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		id, err := getID(r)
		if err != nil {
			return fmt.Errorf("invalid id")
		}
		user, err := s.store.GetUserByID(id)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, user)
	}
	if r.Method == http.MethodPost {
		return s.handleUpdateUser(w, r)
	}

	return fmt.Errorf("method not allowed: %s", r.Method)
}

func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	var req UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "invalid credentials"})
	}
	user := &User{ID: id, EncPassword: req.Password, UserProfile: req.UserProfile}
	updated, err := s.store.UpdateUser(user)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "failed to update user information"})
	}

	return WriteJSON(w, http.StatusOK, updated)
}

type APIError struct {
	Error	string	`json:"error"`
}

type APIMessage struct {
	Message string	`json:"message"`
}

type APIFunc func(http.ResponseWriter, *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error{
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, &APIError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error){
	idStr := mux.Vars(r)["id"]
	return strconv.Atoi(idStr)
}