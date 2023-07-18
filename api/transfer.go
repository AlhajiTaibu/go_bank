package api

import (
	"database/sql"
	"net/http"

	db "github.com/AlhajiTaibu/simplebank/sqlc"
	"github.com/gin-gonic/gin"
)


type createTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID int64 `json:"to_account_id" binding:"required,min-1"`
	Amount int64 `json:"amount" binding:"required"`
}

func (server *Server) createTransfer(ctx *gin.Context){
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err !=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	args := db.TransferMoneyTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: sql.NullInt64{Int64: req.Amount, Valid: true},
	}
	transfer, err := server.store.TransferMoneyTx(ctx, args)

	if err != nil{
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.IndentedJSON(http.StatusCreated, transfer)
}


type getTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTransfer(ctx *gin.Context){
	var req getTransferRequest

	if err := ctx.ShouldBindUri(&req); err !=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	transfer, err := server.store.Queries.GetTransfer(ctx, req.ID)

	if err !=nil{
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.IndentedJSON(http.StatusOK, transfer)
}


func (server *Server) getTransfers(ctx *gin.Context){

	var req listAccountRequest

	if err := ctx.ShouldBindQuery(&req); err!=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListTransfersParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	transfers, err := server.store.Queries.ListTransfers(ctx, args)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	ctx.IndentedJSON(http.StatusFound, transfers)
}


func (server *Server) deleteTransfer(ctx *gin.Context){
	var req getTransferRequest
	
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.Queries.DeleteTransfer(ctx, req.ID)
	
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.IndentedJSON(http.StatusNoContent, "Transfer deleted")
}