package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "simpledice/db/sqlc"
	"simpledice/util"

	"github.com/gin-gonic/gin"
)

type topupWalletRequest struct {
	ToWalletID int64  `json:"to_wallet_id" binding:"required,min=1"`
	Asset      string `json:"asset" binding:"required,asset"`
}

func (server *Server) fundWallet(ctx *gin.Context) {
	var req topupWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, req.ToWalletID, req.Asset) {
		return
	}

	arg := db.TransferTxParams{
		FromWalletID: util.HOLDING_ACCOUNT_ID,
		ToWalletID:   req.ToWalletID,
		Amount:       util.WALLET_FUND_AMOUNT,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, walletID int64, asset string) bool {
	wallet, err := server.store.GetWallet(ctx, walletID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if wallet.Asset != asset {
		err := fmt.Errorf("wallet [%d] asset mismatch: %s vs %s", wallet.ID, wallet.Asset, asset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	if wallet.Balance > 35 {
		err := fmt.Errorf("wallet is greater than the minimum balance of %d required to top up", util.MIN_ALLOWABLE_AMOUNT_TO_FUND_WALLET)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
