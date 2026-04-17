package http

import (
	accounts "assignments/simplebank/accounts/api"
	"assignments/simplebank/shared"

	"github.com/gin-gonic/gin"
)

func Mount(router gin.IRouter, service accounts.Service) {
	r := router.Group("/accounts")
	r.POST("", CreateAccount(service))
}

func CreateAccount(service accounts.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request accounts.Account
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.AbortWithError(400, err)
		}

		_, err := service.CreateAccount(ctx, request.Owner, request.Balance, request.Currency)

		switch err {
		case nil:
			ctx.Status(200)
			return
		case accounts.ErrNegativeBalance:
			ctx.AbortWithError(400, err)
			return
		case accounts.ErrCurrencyRequired:
			ctx.AbortWithError(400, err)
			return
		case accounts.ErrCurrencyNotSupported:
			ctx.AbortWithError(400, err)
			return
		default:
			shared.ErrorToHTTP(ctx, err)
		}

	}
}
