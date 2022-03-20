//  *@createTime    2022/3/20 16:31
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/haylove/small_bank/db/sqlc"
)

//Server servers a http request for our bank service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

//NewServer creates a new Server
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	//todo add routers in router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)

	server.router = router
	return server
}

//Start run a Server in a special address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}