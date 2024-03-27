package registry

import (
	"database/sql"
	"vehicles/packages/adapters/controller"
	"vehicles/packages/adapters/gateway"
	"vehicles/packages/adapters/presenter"
	"vehicles/packages/domain/models"
	usecase "vehicles/packages/usecases/usecases"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewSelectionController(ctx *gin.Context, rdb *redis.Client, vehiclesDB *sql.DB) controller.Selection {
	nsu := usecase.NewSelectionUseCase(
		ctx,
		gateway.NewSelectionRepository(ctx, vehiclesDB),
		gateway.NewCarsRepository(ctx, rdb),
		usecase.NewUserUseCase(gateway.NewUserRepository(ctx)),
		presenter.NewSelectionPresenter(ctx),
		models.User{},
	)
	return controller.NewSelectionController(ctx, nsu)
}
