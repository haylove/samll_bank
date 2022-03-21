//  *@createTime    2022/3/20 16:31
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/haylove/small_bank/db/sqlc"
	"github.com/haylove/small_bank/token"
	"github.com/haylove/small_bank/util"
)

//Server servers a http request for our bank service
type Server struct {
	config         util.Config
	store          db.Store
	router         *gin.Engine
	tokenGenerator token.Generator
}

//NewServer creates a new Server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	generator, err := token.NewPasetoGenerator(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token generator:%w", err)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}

	server := &Server{
		config:         config,
		store:          store,
		tokenGenerator: generator,
	}
	server.setRouter()
	return server, nil
}

//register router
func (s *Server) setRouter() {
	router := gin.Default()
	{
		router.POST("/users", s.createUser)
		router.POST("/users/login", s.loginUser)
	}

	authentication := router.Group("/").Use(authMiddleware(s.tokenGenerator))
	{
		authentication.POST("/accounts", s.createAccount)
		authentication.GET("/accounts", s.listAccounts)
		authentication.GET("/accounts/:id", s.getAccount)

		authentication.POST("/transfers", s.createTransfer)
	}

	s.router = router
}

//Start run a Server in a special address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
