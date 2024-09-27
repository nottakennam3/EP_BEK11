package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gosocial/store"
	"gosocial/types"

	"github.com/gorilla/mux"
)

type apiServer struct {
	addr  string
	store *store.MySQLStorage
}

func NewAPIServer(addr string, store *store.MySQLStorage) *apiServer {
	return &apiServer{
		addr:  addr,
		store: store,
	}
}

func (s *apiServer) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/signup", makeHTTPHandlerFunc(s.handleUserSignup)).Methods(http.MethodPost)
	router.HandleFunc("/login", makeHTTPHandlerFunc(s.handleLogin)).Methods(http.MethodPost)
	router.HandleFunc("/profile", WithJWTAuth(makeHTTPHandlerFunc(s.handleGetUser), s.store)).Methods(http.MethodGet)
	router.HandleFunc("/profile", WithJWTAuth(makeHTTPHandlerFunc(s.handleUpdateUser), s.store)).Methods(http.MethodPut)
	router.HandleFunc("/posts", WithJWTAuth(makeHTTPHandlerFunc(s.handleCreatePost), s.store)).Methods(http.MethodPost)
	router.HandleFunc("/posts/{id}", WithJWTAuth(makeHTTPHandlerFunc(s.handleUpdatePost), s.store)).Methods(http.MethodPut)
	router.HandleFunc("/posts/{id}/like", WithJWTAuth(makeHTTPHandlerFunc(s.handleLikePost), s.store)).Methods(http.MethodPost)
	router.HandleFunc("/posts/{id}/comment", WithJWTAuth(makeHTTPHandlerFunc(s.handleCommentPost), s.store)).Methods(http.MethodPost)

	log.Println("server running at", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func (s *apiServer) handleUserSignup(w http.ResponseWriter, r *http.Request) error {
	var userSignupReq types.UserSignupRequest
	err := json.NewDecoder(r.Body).Decode(&userSignupReq)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	_, err = s.store.GetUserByUsername(userSignupReq.Username)
	if err == nil {
		return fmt.Errorf("username %s already exists", userSignupReq.Username)
	}

	hashed, err := hashPassword(userSignupReq.Password)
	if err != nil {
		return err
	}

	user := types.NewUser(userSignupReq.Username, hashed, userSignupReq.UserProfile)
	if err = s.store.CreateUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, map[string]string{"msg": "user created"})
}

func (s *apiServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var userLoginReq types.UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&userLoginReq)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	user, err := s.store.GetUserByUsername(userLoginReq.Username)
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}

	if !comparePasswords(user.Password, userLoginReq.Password) {
		return fmt.Errorf("invalid credentials")
	}

	token, err := CreateJWT(user.ID)
	if err != nil {
		return ServerError(w)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (s *apiServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	userID := GetUserIDFromContext(r.Context())
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, &apiError{Error: "internal server error"})
	}
	return WriteJSON(w, http.StatusOK, user)
}

func (s *apiServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	var userUpdateReq types.UserUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&userUpdateReq)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	if userUpdateReq.Password == "" && userUpdateReq.UserProfile == "" {
		return fmt.Errorf("no info to update")
	}
	user := &types.User{
		ID:          GetUserIDFromContext(r.Context()),
		UserProfile: userUpdateReq.UserProfile,
	}

	if userUpdateReq.Password != "" {
		hashed, err := hashPassword(userUpdateReq.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
	}

	if err := s.store.UpdateUser(user); err != nil {
		return ServerError(w)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"msg": "user info updated"})
}

func (s *apiServer) handleCreatePost(w http.ResponseWriter, r *http.Request) error {
	var postCreateReq types.PostCreateRequest
	if err := decodeRequest(r, &postCreateReq); err != nil {
		return err
	}

	userID := GetUserIDFromContext(r.Context())
	post := types.NewPost(userID, postCreateReq.Content)
	if err := s.store.CreatePost(post); err != nil {
		return ServerError(w)
	}

	return WriteJSON(w, http.StatusCreated, map[string]string{"msg": "post created"})
}

func (s *apiServer) handleUpdatePost(w http.ResponseWriter, r *http.Request) error {
	postID, err := getID(r)
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	post, err := s.store.GetPostByID(postID)
	if err != nil {
		return ServerError(w)
	}
	if post.ID == 0 {
		return fmt.Errorf("post not found")
	}

	currentUserID := GetUserIDFromContext(r.Context())
	if post.UserID != currentUserID {
		return WriteJSON(w, http.StatusForbidden, &apiError{Error: "permission denied"})
	}

	var postUpdateRequest types.PostUpdateRequest
	if err := decodeRequest(r, &postUpdateRequest); err != nil {
		return err
	}

	post.Content = postUpdateRequest.Content
	if err := s.store.UpdatePost(post); err != nil {
		return ServerError(w)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"msg": "post updated"})
}

func (s *apiServer) handleLikePost(w http.ResponseWriter, r *http.Request) error {
	postID, err := getID(r)
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	post, err := s.store.GetPostByID(postID)
	if err != nil {
		return ServerError(w)
	}

	if post.ID == 0 {
		return fmt.Errorf("post not found")
	}

	userID := GetUserIDFromContext(r.Context())

	like, err := s.store.GetPostLikeByUserID(postID, userID)
	if err != nil {
		return ServerError(w)
	}
	var msg string
	if like.ID == 0 {
		err = s.store.LikePost(postID, userID)
		msg = "post unliked"
	} else {
		err = s.store.UnlikePost(postID, userID)
		msg = "post liked"
	}

	if err != nil {
		return ServerError(w)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"msg": msg})
}

func (s *apiServer) handleCommentPost(w http.ResponseWriter, r *http.Request) error {
	postID, err := getID(r)
	if err != nil {
		return fmt.Errorf("invalid id format")
	}

	post, err := s.store.GetPostByID(postID)
	if err != nil {
		return ServerError(w)
	}

	if post.ID == 0 {
		return fmt.Errorf("post not found")
	}

	var postCommentReq types.PostCommentRequest
	if err := decodeRequest(r, &postCommentReq); err != nil {
		return ServerError(w)
	}

	userID := GetUserIDFromContext(r.Context())
	postComment := types.NewPostComment(post.ID, userID, postCommentReq.Content)
	if err := s.store.CommentPost(postComment); err != nil {
		return ServerError(w)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"msg": "comment submitted"})
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type apiError struct {
	Error string `json:"error"`
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, &apiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func ServerError(w http.ResponseWriter) error {
	return WriteJSON(w, http.StatusInternalServerError, &apiError{Error: "internal server error"})
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	return strconv.Atoi(idStr)
}

func decodeRequest(r *http.Request, payload any) error {
	err := json.NewDecoder(r.Body).Decode(payload)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
