package api

import (
	"database/sql"
	"net/http"

	db "github.com/AlhajiTaibu/simplebank/sqlc"
	"github.com/gin-gonic/gin"

)

type createEntryRequest struct{
	AccountID int64 `json:"account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required"` 
}

func (server *Server) createEntry(ctx *gin.Context){
	var req createEntryRequest

	if err := ctx.ShouldBindJSON(&req); err !=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	account, err := server.store.GetAccount(ctx, req.AccountID)
	if err !=nil{
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}

	args := db.CreateEntryParams{
		AccountID: account.ID,
		Amount: req.Amount,
	}

	entry, err := server.store.CreateEntry(ctx, args)

	if err !=nil{
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.IndentedJSON(http.StatusCreated, entry)
}

type getEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getEntry(ctx *gin.Context){
	var req getEntryRequest

	if err := ctx.ShouldBindUri(&req); err!=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	entry, err := server.store.GetEntry(ctx, req.ID)

	if err != nil{
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.IndentedJSON(http.StatusOK, entry)
}

type EntriesQueryRequest struct{
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=5"`
}

func (server *Server) getEntries(ctx *gin.Context){
	var req EntriesQueryRequest

	if err := ctx.ShouldBindQuery(&req); err!=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	args := db.ListEntriesParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	entries, err := server.store.ListEntries(ctx, args)

	if err != nil{
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.IndentedJSON(http.StatusOK, entries)
}

type updateEntryRequest struct {
	AccountID int64 `json:"account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required"` 
}

type updateEntryRequestParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) updateEntry(ctx *gin.Context){
	var req updateEntryRequest
	var req1 updateEntryRequestParams
	if err := ctx.ShouldBindJSON(&req); err !=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	if err := ctx.ShouldBindUri(&req1); err !=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}

	entry1, err := server.store.GetEntry(ctx, req1.ID)

	if err !=nil{
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}

	account, err := server.store.GetAccount(ctx, req.AccountID)
	if err !=nil{
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}

	args := db.UpdateEntryParams{
		ID: entry1.ID,
		AccountID: account.ID,
		Amount: req.Amount,
	}

	entry, err := server.store.UpdateEntry(ctx, args)

	if err !=nil{
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.IndentedJSON(http.StatusOK, entry)
}


func (server *Server) deleteEntry(ctx *gin.Context){
	var req getEntryRequest

	if err := ctx.ShouldBindUri(&req); err !=nil{
		ctx.IndentedJSON(http.StatusBadRequest, errorResponse(err))
	}
	err := server.store.DeleteEntry(ctx, req.ID)

	if err !=nil{
		if err == sql.ErrNoRows{
			ctx.IndentedJSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.IndentedJSON(http.StatusNoContent, "Entry deleted")
}