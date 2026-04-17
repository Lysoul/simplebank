package accounts

import (
	"assignments/simplebank/accounts/http"
	"assignments/simplebank/accounts/internal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB, router gin.IRouter) {
	store := internal.NewStore(db)
	service := internal.NewService(store)
	http.Mount(router, service)
}
