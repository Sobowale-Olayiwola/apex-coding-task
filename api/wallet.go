package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "simpledice/db/sqlc"

	"github.com/gin-gonic/gin"
)

type getWalletBalanceRequest struct {
	Username string `uri:"username" binding:"required"`
}

func (server *Server) getWalletBalance(ctx *gin.Context) {
	var req getTransactionLogsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	wallet, err := server.store.GetWalletByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(db.ErrNoRecord))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"balance":      fmt.Sprintf("%d", wallet.Balance),
		"wallet_asset": wallet.Asset,
	})
}
