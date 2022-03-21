//  *@createTime    2022/3/20 16:48
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"github.com/haylove/small_bank/token"
	"github.com/lib/pq"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/haylove/small_bank/db/sqlc"
)

// createAccountRequest is the form to createAccount
type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// createAccount is the api for POST:'/accounts'
func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	// it has been authenticated,so we get payload from context for create account
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    payload.Username,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

// getAccountRequest is the form to getAccount,it just needs ID
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount is the api for GET:"/account:id"
func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, ok := s.validateAccount(ctx, req.ID, "", true)
	if !ok {
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// listAccountsRequest is the form for listAccounts,it needs to be paginated data
type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAccounts is the api for GET:"/accounts"
func (s *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		Owner:  payload.Username,
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
