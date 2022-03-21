//  *@createTime    2022/3/21 1:57
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/haylove/small_bank/token"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/haylove/small_bank/db/sqlc"
)

type transferRequest struct {
	FromAccountID int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64   `json:"to_account_id" binding:"required,min=1"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Currency      string  `json:"currency" binding:"required,currency"`
}

//createTransfer is the api for POST:/transfer
func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	//validate account and currency
	if _, ok := s.validateAccount(ctx, req.FromAccountID, req.Currency, true); !ok {
		return
	}
	if _, ok := s.validateAccount(ctx, req.ToAccountID, req.Currency, false); !ok {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// validateAccount validates the accountID  currency ,returns the account and the currency is for the account.
// If the currency is "",it will return the account with accountID if existed.
// If false ,that means  the currency is not for the account.
func (s *Server) validateAccount(ctx *gin.Context, accountID int64, currency string, needAuth bool) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return db.Account{}, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return db.Account{}, false
	}
	if needAuth {
		payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		if account.Owner != payload.Username {
			err := errors.New("account doesn't belong to the authenticated user")
			ctx.JSON(http.StatusForbidden, errResponse(err))
			return db.Account{}, false
		}
	}
	if currency != "" {
		if account.Currency != currency {
			err = fmt.Errorf("account [%d] mismatch: %v vs %v", account.ID, account.Currency, currency)
			ctx.JSON(http.StatusBadRequest, errResponse(err))
			return db.Account{}, false
		}
	}
	return account, true
}
