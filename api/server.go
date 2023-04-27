package api

import (
	db "simpledice/db/sqlc"
	"simpledice/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our wallet service
type Server struct {
	config util.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("asset", validAsset)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/start-game", server.startUserGame)
	router.PUT("/users/:sessionID/end-game", server.endSession)
	router.GET("/users/:sessionID/check-session", server.checkSession)
	router.POST("/users/roll-dice", server.rollDice)
	router.GET("/users/transaction-logs/:username", server.getTransactionLogs)

	router.POST("/transfers/fund-wallet", server.fundWallet)
	router.GET("wallets/:username/balance", server.getWalletBalance)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"success": false, "error": err.Error()}
}
