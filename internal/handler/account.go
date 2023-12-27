package handler

import (
	"fmt"
	"net/http"
	"rest1/internal/domain"
	"rest1/internal/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AccountHandler struct {
	UseCase *usecases.AccountUsecase
	logger  *zap.Logger
}

func NewAccountHandler(useCase *usecases.AccountUsecase, logger *zap.Logger) *AccountHandler {
	return &AccountHandler{
		UseCase: useCase,
		logger:  logger,
	}
}
func (h *AccountHandler) DropAccountsTableHandler(w http.ResponseWriter, r *http.Request) {
	err := h.UseCase.DropAccountsTable()

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
	}

	respondWithJSON(w, http.StatusOK, nil)
}
func (h *AccountHandler) CreateAccountTableHandler(w http.ResponseWriter, r *http.Request) {

	err := h.UseCase.CreateAccountTable()

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, nil)
}

// Create Account route
// http://localhost:8000/account/create/{userid}
func (h *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userid")
	userId, err := uuid.Parse(idStr)
	fmt.Println(userId)

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
	}

	accid, err := h.UseCase.CreateAccount(userId)
	fmt.Println("handler account id", accid)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
	}
	var account *domain.Account
	account, err = h.UseCase.GetByAccountNo(accid)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, account)
}

// http://localhost:3000/account/{accoundId}
func (h *AccountHandler) GetByAccountNoHandler(w http.ResponseWriter, r *http.Request) {
	// get ID from url parameters
	idStr := chi.URLParam(r, "accoundId")
	accountId, err := uuid.Parse(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Failed to get account by ID at Handler layer", zap.Error(err))
		return
	}

	acc, err := h.UseCase.GetByAccountNo(accountId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	respondWithJSON(w, http.StatusOK, acc)
}
