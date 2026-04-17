package http

import (
	accounts "assignments/simplebank/accounts/api"
	"assignments/simplebank/shared"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Mount(router gin.IRouter, service accounts.Service) {
	r := router.Group("/accounts")
	r.POST("", CreateAccount(service))
	r.GET("/:id", GetAccount(service))
	r.GET("", ListAccounts(service))
	r.DELETE("/:id", DeleteAccount(service))
	r.POST("/transfers", CreateTransfer(service))

}

func CreateAccount(service accounts.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request accounts.Account
		if err := ctx.ShouldBindJSON(&request); err != nil {
			_ = ctx.AbortWithError(400, err)
			return
		}

		account, err := service.CreateAccount(ctx, request.Owner, request.Balance, request.Currency)
		switch err {
		case nil:
			ctx.JSON(200, account)
			return
		case accounts.ErrNegativeBalance:
			_ = ctx.AbortWithError(400, err)
			return
		case accounts.ErrCurrencyRequired:
			_ = ctx.AbortWithError(400, err)
			return
		case accounts.ErrCurrencyNotSupported:
			_ = ctx.AbortWithError(400, err)
			return
		default:
			shared.ErrorToHTTP(ctx, err)
		}

	}
}

func GetAccount(service accounts.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("id")
		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			_ = ctx.AbortWithError(400, accounts.ErrInvalidID)
			return
		}

		result, err := service.GetAccount(ctx, id)
		switch err {
		case nil:
			ctx.JSON(200, result)
			return
		case accounts.ErrAccountNotFound:
			_ = ctx.AbortWithError(404, err)
			return
		default:
			shared.ErrorToHTTP(ctx, err)
		}
	}
}

func ListAccounts(service accounts.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//todo: add pagination support
		results, err := service.ListAccounts(ctx)
		if err != nil {
			shared.ErrorToHTTP(ctx, err)
			return
		}

		ctx.JSON(200, results)
	}
}

func DeleteAccount(service accounts.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("id")
		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			_ = ctx.AbortWithError(400, accounts.ErrInvalidID)
			return
		}

		err = service.DeleteAccount(ctx, id)
		switch err {
		case nil:
			ctx.Status(204)
			return
		case accounts.ErrAccountNotFound:
			_ = ctx.AbortWithError(404, err)
			return
		default:
			shared.ErrorToHTTP(ctx, err)
		}
	}
}

type TransferRequest struct {
	FromAccountID int64 `json:"fromAccountId" binding:"required,gt=0"`
	ToAccountID   int64 `json:"toAccountId" binding:"required,gt=0"`
	Amount        int64 `json:"amount" binding:"required,gt=0,lte=1000000"`
}

func CreateTransfer(service accounts.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request TransferRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			_ = ctx.AbortWithError(400, err)
			return
		}

		transfer, err := service.TransferTx(ctx, request.FromAccountID, request.ToAccountID, request.Amount)
		switch err {
		case nil:
			ctx.JSON(200, transfer)
			return
		case accounts.ErrInvalidTransferAmount:
			_ = ctx.AbortWithError(400, err)
			return
		case accounts.ErrSameAccountTransfer:
			_ = ctx.AbortWithError(400, err)
			return
		case accounts.ErrAccountNotFound:
			_ = ctx.AbortWithError(404, err)
			return
		case accounts.ErrInsufficientBalance:
			_ = ctx.AbortWithError(400, err)
			return
		default:
			shared.ErrorToHTTP(ctx, err)
		}

	}
}
