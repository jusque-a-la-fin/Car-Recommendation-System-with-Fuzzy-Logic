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

func NewSearchController(ctx *gin.Context, rdb *redis.Client, pdb *sql.DB) controller.Search {
	nur := usecase.NewUserUseCase(gateway.NewUserRepository(ctx))
	ncr := gateway.NewCarsRepository(ctx, rdb)
	nsp := presenter.NewSearchPresenter(ctx)
	nsr := gateway.NewSearchRepository(ctx, rdb)
	nsu := usecase.NewSearchUseCase(
		nsr,
		ncr,
		nur,
		usecase.NewQuestionUseCase(
			gateway.NewQuestionRepository(ctx, pdb), ncr, nur, nsp,
		),
		nsp,
	)
	return controller.NewSearchController(ctx, rdb, nsu)
}
