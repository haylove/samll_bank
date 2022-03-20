//  *@createTime    2022/3/20 16:48
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/haylove/small_bank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, ok := s.validateAccount(ctx, req.ID, "")
	if !ok {
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
