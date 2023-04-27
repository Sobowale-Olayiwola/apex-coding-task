package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "simpledice/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Asset    string `json:"asset" binding:"required,asset"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.CreateUserTx(ctx, req.Username, req.Asset)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

type startUserGameRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
}

func (server Server) startUserGame(ctx *gin.Context) {
	var req startUserGameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	rsp, err := server.store.StartGame(ctx, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrActiveSession), errors.Is(err, db.ErrInsufficientFund):
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		case errors.Is(err, db.ErrNoRecord):
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": rsp})
}

type roleDiceRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
}

func (server *Server) rollDice(ctx *gin.Context) {
	var req roleDiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserWithSession(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("kindly start a new game")))
			return
		}
	}
	switch user.NumOfDiceThrow {
	case 0:
		err := server.store.FirstDiceThrow(ctx, user)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrInsufficientFund):
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			default:
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"success": true})
		return
	case 1:
		_, err := server.store.SecondDiceThrow(ctx, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"success": true})
		return
	case 2:
		_, err := server.store.StartNewDiceThrow(ctx, user)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrInsufficientFund):
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			default:
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
}

type endSessionRequest struct {
	SessionID int64 `uri:"sessionID" binding:"required,min=1"`
}

func (server *Server) endSession(ctx *gin.Context) {
	var req endSessionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.EndUserGameSession(ctx, req.SessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("no active game in session")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
	return
}

type checkSessionRequest struct {
	SessionID int64 `uri:"sessionID" binding:"required,min=1"`
}

func (server *Server) checkSession(ctx *gin.Context) {
	var req checkSessionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	usersession, err := server.store.GetSession(ctx, req.SessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(db.ErrNoRecord))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"active": usersession.IsActive})
	return
}

type getTransactionLogsRequest struct {
	Username string `uri:"username" binding:"required"`
}

func (server *Server) getTransactionLogs(ctx *gin.Context) {
	var req getTransactionLogsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	txnLogs, err := server.store.GetUserTransactionLogs(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(db.ErrNoRecord))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": txnLogs})
}
