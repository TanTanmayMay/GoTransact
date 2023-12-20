package handler

import (
	"encoding/json"
	"net/http"
	"rest1/internal/domain"
	"rest1/internal/usecases"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)


type UserHandler struct {
	UseCase *usecases.UserUsecase
	conn *pgx.Conn
}

func NewUserHandler(useCase *usecases.UserUsecase , conn *pgx.Conn) *UserHandler{
	return &UserHandler{
		UseCase: useCase,
		conn: conn,
	}
}


// Create User route
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request){
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)

	// check if user from req.body is valid
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	err = h.UseCase.CreateUser(&user, h.conn)

	
	respondWithJSON(w, http.StatusOK, user)
}



func (h *UserHandler) GetAccountById(w http.ResponseWriter, r *http.Request){
	// get ID from url parameters
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return 
	}


	user, err := h.UseCase.GetAccountByID(id, h.conn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	respondWithJSON(w, http.StatusOK, user)
}

// withdraw money
func (h *UserHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request){
	
}


// Utitlity function to response in JSON
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}