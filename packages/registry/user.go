package registry

import (
	"vehicles/packages/adapters/controller"
	"vehicles/packages/adapters/gateway"
	usecase "vehicles/packages/usecases/usecases"

	"github.com/gin-gonic/gin"
)

func NewUserController(ctx *gin.Context) controller.User {
	nuu := usecase.NewUserUseCase(
		gateway.NewUserRepository(ctx),
	)
	return controller.NewUserController(ctx, nuu)
}
