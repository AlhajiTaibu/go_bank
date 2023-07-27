package api

import (
	"database/sql"
	"net/http"

	db "github.com/AlhajiTaibu/simplebank/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var newAccount createAccountRequest
	if err := ctx.ShouldBindJSON(&newAccount); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.CreateAccountParams{
		Owner:    newAccount.Owner,
		Currency: newAccount.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, args)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.IndentedJSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	ctx.IndentedJSON(http.StatusOK, account)
}


type listAccountRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=5"`
}

func (server *Server) getAccounts(ctx *gin.Context){
	var req listAccountRequest

	if err := ctx.ShouldBindQuery(&req); err!=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, args)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	ctx.IndentedJSON(http.StatusOK, accounts)
}


type updateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (server *Server) updateAccount(ctx *gin.Context){
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqx updateAccountRequest
	if err := ctx.ShouldBindJSON(&reqx); err != nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account1, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	args := db.UpdateAccountParams{
		ID: account1.ID,
		Balance: reqx.Balance,
		Owner: reqx.Owner,
		Currency: reqx.Currency,
	}
	account, err := server.store.UpdateAccount(ctx, args)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.IndentedJSON(http.StatusOK, account)
}

func (server *Server) deleteAccount(ctx *gin.Context){
	var req getAccountRequest
	
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.IndentedJSON(http.StatusNoContent, "Account deleted")
}
