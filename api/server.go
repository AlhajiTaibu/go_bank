package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/AlhajiTaibu/simplebank/sqlc"
)

type Server struct {
	store db.Store
	router  *gin.Engine
}

func NewServer(store db.Store) *Server{
	server := &Server{
		store: store,
	}
	router := gin.Default()
	
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAccounts)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	// entries routes
	router.POST("/entries", server.createEntry)
	router.GET("/entries/:id", server.getEntry)
	router.GET("/entries", server.getEntries)
	router.PUT("/entries/:id", server.updateEntry)
	router.DELETE("/entries/:id", server.deleteEntry)
	// transfer routes
	router.POST("/transfers", server.createTransfer)
	router.GET("/transfers/:id", server.getTransfer)
	router.GET("/transfers", server.getTransfers)
	router.DELETE("/transfers/:id", server.deleteTransfer)


	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H{
	return gin.H{"error": err.Error()}
}