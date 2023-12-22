package handler

import (
	"encoding/json"
	"net/http"
	"rest1/internal/domain"
	"rest1/internal/usecases"
	"strconv"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (h *UserHandler) DropUserTableHandler(w http.ResponseWriter, r *http.Request){
	err := h.UseCase.DropUserTable()
	if(err != nil){
		respondWithJSON(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, nil)
}

func (h *UserHandler) CreateUsersTableHandler(w http.ResponseWriter, r *http.Request){
	err := h.UseCase.CreateUserTable()

	if(err != nil){
		respondWithJSON(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, nil)
}

// Create User route
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request){
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	user.ID = uuid.New()
	// check if user from req.body is valid
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Error("Invalid Data to create User")
		return 
	}
	// newAccId, err := h.UseCase.CreateAccount(user.ID, h.conn)
	// if err != nil {
	// 	h.logger.Error("Error while creating account by user.go handler")
	// 	return 
	// }

	err = h.UseCase.CreateUser(&user, h.conn)
	if err != nil {
		// http.Error(w, "Error while creating user", http.StatusInternalServerError)
		// h.logger.Error("Error while creating user at creatUser", zap.Error(err))
		// return 
		respondWithJSON(w, http.StatusInternalServerError, err)
		return 
	}
	
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


func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request){
	// get ID from url parameters
	idStr := chi.URLParam(r, "userid")
	id, err := uuid.Parse(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Invalid ID while getting user at Handler Layer", zap.Error(err))
		return 
	}
	user, err := h.UseCase.GetUserById(id, h.conn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		h.logger.Error("Invalid ID to get get account by ID", zap.Error(err))
	}
	respondWithJSON(w, http.StatusOK, user)
}

// withdraw money
// localhost:3000/
func (h *UserHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request){
	/* 
		//Tanmay
		{
			"userid" : "-----",
			"amount" : 12
		}

		var req
		err := json.NewDecoder(r.Body).Decode(&user)
	*/
	idStr := chi.URLParam(r, "userid")
	id, err := uuid.Parse(idStr)
	if(err != nil){
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Invalid ID while getting user at Handler Layer", zap.Error(err))
		return 
	} 
	user, err := h.UseCase.GetUserById(id, h.conn)

	amountStr := chi.URLParam(r, "amount")
	var amount int
	amount , err = strconv.Atoi(amountStr)

	if(err != nil){
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Cannot fetch user", zap.Error(err))
		return 
	} 
	
	// func (a *UserUsecase) Withdraw(user *domain.User, amount int, conn *pgx.Conn) error {
	err = h.UseCase.Withdraw(user, amount, h.conn)

	if(err != nil){
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Cannot Withdraw as Amount goes below minimum balance", zap.Error(err))
		return 
	} 
}

func (h *UserHandler) DepositHandler(w http.ResponseWriter, r *http.Request){
	idStr := chi.URLParam(r, "userid")
	id, err := uuid.Parse(idStr)
	if(err != nil){
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Invalid ID while getting user at Handler Layer", zap.Error(err))
		return 
	} 

	user, err := h.UseCase.GetUserById(id, h.conn)

	amountStr := chi.URLParam(r, "amount")
	var amount int
	amount , err = strconv.Atoi(amountStr)

	if(err != nil){
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Cannot fetch user", zap.Error(err))
		return 
	} 
	
	// func (a *UserUsecase) Deposit(user *domain.User, amount int, conn *pgx.Conn) error {
	err = h.UseCase.Deposit(user, amount, h.conn)

	if(err != nil){
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Cannot Deposit", zap.Error(err))
		return 
	} 
}
	


// Utitlity function to response in JSON
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}