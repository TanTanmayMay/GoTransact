package handler

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"rest1/internal/domain"
	"rest1/internal/usecases"

	// "strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	UseCase *usecases.UserUsecase
	logger  *zap.Logger
}

func NewUserHandler(useCase *usecases.UserUsecase, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		UseCase: useCase,
		logger:  logger,
	}
}

func (h *UserHandler) DropUserTableHandler(w http.ResponseWriter, r *http.Request) {
	err := h.UseCase.DropUserTable()
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, nil)
}

func (h *UserHandler) CreateUsersTableHandler(w http.ResponseWriter, r *http.Request) {
	err := h.UseCase.CreateUserTable()

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, nil)
}

// Create User route
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	// err := json.NewDecoder(r.Body).Decode(&user)

	/*
		if err := render.Decode(r.Body, &user); err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}
	*/

	if err := render.Decode(r, &user); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New()

	err := h.UseCase.CreateUser(&user)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	render.JSON(w, r, user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.UseCase.GetAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error("Failed to get all users  at Handler layer", zap.Error(err))
		return
	}

	render.JSON(w, r, users)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	// get ID from url parameters
	idStr := chi.URLParam(r, "userid")
	id, err := uuid.Parse(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Invalid ID while getting user at Handler Layer", zap.Error(err))
		return
	}
	user, err := h.UseCase.GetUserById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		h.logger.Error("Invalid ID to get get account by ID", zap.Error(err))
	}
	render.JSON(w, r, user)

}

// withdraw money
// localhost:3000/
/*
func (h *UserHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	/*
		//Tanmay
		{
			"userid" : "-----",
			"amount" : 12
		}

		var req
		err := json.NewDecoder(r.Body).Decode(&user)
*/ /*
	idStr := chi.URLParam(r, "userid")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Invalid ID while getting user at Handler Layer", zap.Error(err))
		return
	}
	user, err := h.UseCase.GetUserById(id)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Cannot fetch user", zap.Error(err))
		return
	}

	amountStr := chi.URLParam(r, "amount")
	var amount int
	amount, err = strconv.Atoi(amountStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Cannot convert amount to integer", zap.Error(err))
		return
	}

	// func (a *UserUsecase) Withdraw(user *domain.User, amount int, conn *pgxpool.Conn) error {
	err = h.UseCase.Withdraw(user, amount)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Cannot Withdraw as Amount goes below minimum balance", zap.Error(err))
		return
	}
}

func (h *UserHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userid")
	userid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Invalid ID while getting user at Handler Layer", zap.Error(err))
		return
	}

	user, err := h.UseCase.GetUserById(userid) // NO ISSUE
	fmt.Println("user from handler", user.Name)
	amountStr := chi.URLParam(r, "amount")
	var amount int
	amount, err = strconv.Atoi(amountStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Cannot fetch user", zap.Error(err))
		return
	}

	// func (a *UserUsecase) Deposit(user *domain.User, amount int, conn *pgxpool.Conn) error {
	err = h.UseCase.Deposit(user, amount)

	if err != nil {
		fmt.Println("ERror from handler user", err)
		respondWithJSON(w, http.StatusBadRequest, err)
	}
}
*/

// Utitlity function to response in JSON
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
