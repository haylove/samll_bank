//  *@createTime    2022/3/21 1:57
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"database/sql"
	"fmt"
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

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	if _, ok := s.validateAccount(ctx, req.FromAccountID, req.Currency); !ok {
		return
	}
	if _, ok := s.validateAccount(ctx, req.ToAccountID, req.Currency); !ok {
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

func (s *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return db.Account{}, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return db.Account{}, false
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
