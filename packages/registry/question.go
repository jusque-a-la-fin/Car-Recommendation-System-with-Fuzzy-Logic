package registry

import (
	"database/sql"
	"vehicles/packages/adapters/controller"
	"vehicles/packages/adapters/gateway"
	"vehicles/packages/adapters/presenter"
	usecase "vehicles/packages/usecases/usecases"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewQuestionController(ctx *gin.Context, rdb *redis.Client, pdb *sql.DB) controller.Question {
	nqr := gateway.NewQuestionRepository(ctx, pdb)
	nuc := usecase.NewUserUseCase(gateway.NewUserRepository(ctx))
	nsp := presenter.NewSearchPresenter(ctx)
	ncr := gateway.NewCarsRepository(ctx, rdb)
	nqu := usecase.NewQuestionUseCase(nqr, ncr, nuc, nsp)
	return controller.NewQuestionController(ctx, nqu)
}
