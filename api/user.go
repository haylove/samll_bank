//  *@createTime    2022/3/21 5:21
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/haylove/small_bank/db/sqlc"
	"github.com/haylove/small_bank/util"
	"github.com/lib/pq"
	"net/http"
	"time"
	"unsafe"
)

// userResponse prevent  password from being returned,it just copied form db.User
type userResponse struct {
	Username          string `json:"username"`
	_                 string
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// createUserRequest is the createUser form
type createUserRequest struct {
	Username string `json:"Username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// createUser is the api for POST: '/users'
func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	//just speculation
	userRes := *(*userResponse)(unsafe.Pointer(&user))
	ctx.JSON(http.StatusOK, userRes)
}

// loginUserRequest is the login form
type loginUserRequest struct {
	Username string `json:"Username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

// loginResponse contains user information and token
type loginResponse struct {
	AccessToken string        `json:"access_token"`
	User        *userResponse `json:"user"`
}

// loginUser is the api for POST: '/users/login'
func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	//get user by username
	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	//check password
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	//had passed,generate a token for user
	token, err := s.tokenGenerator.GenerateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	res := loginResponse{
		AccessToken: token,
		User:        (*userResponse)(unsafe.Pointer(&user)),
	}
	ctx.JSON(http.StatusOK, res)
}
