package handler

import (
	"encoding/json"
	"net/http"
	"rest1/internal/domain"
	"rest1/internal/usecases"
	"strconv"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)


type UserHandler struct {
	UseCase *usecases.UserUsecase
	conn *pgx.Conn
	logger *zap.Logger
}

func NewUserHandler(useCase *usecases.UserUsecase , conn *pgx.Conn , logger *zap.Logger) *UserHandler{
	return &UserHandler{
		UseCase: useCase,
		conn: conn,
		logger: logger,
	}
}


// Create User route
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request){
	var user domain.User

	/*
		var myCourse onlineCourses

		isValid := json.Valid(jsonFromWeb)
		if isValid {
			fmt.Println("Json is Valid")
			json.Unmarshal(jsonFromWeb, &myCourse)
			fmt.Printf("%#v", myCourse)
		} else {
			fmt.Println("Not a Valid JSON")
		}


		// Decode JSON data from the request body
		var requestData YourStruct
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Now, 'requestData' contains the mapped struct values
		fmt.Printf("Field1: %s, Field2: %d\n", requestData.Field1, requestData.Field2)
	*/

	err := json.NewDecoder(r.Body).Decode(&user)

	// check if user from req.body is valid
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Error("Invalid Data to create User", zap.Error(err))
		return 
	}

	newAccId, err := h.UseCase.CreateAccount(user.ID, h.conn)
	if err != nil {
		h.logger.Error("Error while creating account by user.go handler", zap.Error(err))
		return 
	}
	user.AccountNo = newAccId
	err = h.UseCase.CreateUser(&user, h.conn)
	if err != nil {
		h.logger.Error("Error in creating user at Handler Layer" , zap.Error(err))
		return 
	}
	// fmt.Println("no use", idd)
	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request){

	users, err := h.UseCase.GetAll(h.conn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error("Failed to get all users  at Handler layer", zap.Error(err))
		return 
	}

	respondWithJSON(w, http.StatusOK, users)
}


func (h *UserHandler) GetAccountById(w http.ResponseWriter, r *http.Request){
	// get ID from url parameters
	idStr := chi.URLParam(r, "userid")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Invalid ID while getting user at Handler Layer", zap.Error(err))

		return 
	}
	user, err := h.UseCase.GetAccountByID(id, h.conn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		h.logger.Error("Invalid ID to get get account by ID", zap.Error(err))
	}
	respondWithJSON(w, http.StatusOK, user)
}

// withdraw money
// localhost:3000/
func (h *UserHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request){

}


// Utitlity function to response in JSON
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}