package handler

import (
	"net/http"
	"rest1/internal/domain"
	"rest1/internal/usecases"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)


type AccountHandler struct {
	UseCase *usecases.AccountUsecase
	conn *pgx.Conn
	logger *zap.Logger
}

func NewAccountHandler(useCase *usecases.AccountUsecase , conn *pgx.Conn, logger *zap.Logger) *AccountHandler{
	return &AccountHandler{
		UseCase: useCase,
		conn: conn,
		logger: logger,
	}
}
func (h *AccountHandler) DropAccountsTableHandler(w http.ResponseWriter, r *http.Request){
	err := h.UseCase.DropAccountsTable()

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
	}

	respondWithJSON(w, http.StatusOK, nil)
}
func (h *AccountHandler) CreateAccountTableHandler(w http.ResponseWriter, r *http.Request){

	err := h.UseCase.CreateAccountTable()

	if(err != nil){
		respondWithJSON(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, nil)
}

// Create Account route
// http://localhost:8000/account/create/{userid}
func (h *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request){
	var account domain.Account
	idStr := chi.URLParam(r, "userid")
	userId , err := uuid.Parse(idStr)

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err)
	}

	_, err = h.UseCase.CreateAccount(userId, h.conn)
	
	if err != nil{
		respondWithJSON(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, account)
}


// http://localhost:3000/account/{accoundId}
func (h *AccountHandler) GetByAccountNoHandler(w http.ResponseWriter, r *http.Request){
	// get ID from url parameters
	idStr := chi.URLParam(r, "accoundId")
	accountId, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.logger.Error("Failed to get account by ID at Handler layer", zap.Error(err))
		return 
	}


	acc, err := h.UseCase.GetByAccountNo(accountId, h.conn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	respondWithJSON(w, http.StatusOK, acc)
}
