package router

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"vehicles/packages/registry"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func MakeNewRouter(router *gin.Engine, redisSearchDB *redis.Client, redisSelectionDB *redis.Client, surveyDB *sql.DB, vehiclesDB *sql.DB) *gin.Engine {
	router.GET("main", func(ctx *gin.Context) {
		registry.NewSearchController(ctx, redisSearchDB, surveyDB).DisplayMainPage()
	})

	router.POST("main", func(ctx *gin.Context) {
		err := registry.NewUserController(ctx).SaveUserID()
		if err != nil {
			fmt.Printf("error from `SaveUserID` method, package `controller`: %#v", err)
			errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
			if errAbort != nil {
				fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
			}
		}
		err = registry.NewSearchController(ctx, redisSearchDB, surveyDB).GetSeachCars()
		if err != nil {
			fmt.Printf("error from `GetSeachCars` method, package `controller`: %#v", err)
			errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
			if errAbort != nil {
				fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
			}
		}
	})

	router.GET("search", func(ctx *gin.Context) {
		sessionID := ctx.Query("guest")
		thisCarID := ctx.Query("carID")
		if thisCarID != "" {
			carID, err := strconv.Atoi(thisCarID)
			if err != nil {
				fmt.Printf("error from `Atoi` function, package `strconv`: %#v", err)
			}
			err = registry.NewSearchController(ctx, redisSearchDB, surveyDB).DisplaySearchCarAd(sessionID, carID)
			if err != nil {
				fmt.Printf("error from `DisplaySearchCarAd` method, package `controller`: %#v", err)
				errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
				if errAbort != nil {
					fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
				}
			}
		} else {
			err := registry.NewSearchController(ctx, redisSearchDB, surveyDB).TransferSearchCarsData(sessionID)
			if err != nil {
				fmt.Printf("error from `TransferSearchCarsData` method, package `controller`: %#v", err)
				errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
				if errAbort != nil {
					fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
				}
			}
		}
	})

	router.POST("search", func(ctx *gin.Context) {
		err := registry.NewQuestionController(ctx, redisSearchDB, surveyDB).GetAnswer()
		if err != nil {
			fmt.Printf("error from `GetAnswer` method, package `controller`: %#v", err)
			errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
			if errAbort != nil {
				fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
			}
		}
	})

	ServeSelection(router, redisSelectionDB, vehiclesDB)

	return router
}

func ServeSelection(router *gin.Engine, redisSelectionDB *redis.Client, vehiclesDB *sql.DB) {
	selection := router.Group("/selection")
	{
		selection.GET("priorities", func(ctx *gin.Context) {
			registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).ChoosePriorities()
		})

		selection.POST("priorities", func(ctx *gin.Context) {
			err := registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).PutPriorities()
			if err != nil {
				fmt.Printf("error from `PutPriorities` method, package `controller`: %#v", err)
				errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
				if errAbort != nil {
					fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
				}
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "Данные успешно получены"})
		})

		selection.GET("price", func(ctx *gin.Context) {
			registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).ChoosePrice()
		})

		selection.POST("price", func(ctx *gin.Context) {
			err := registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).PutPrice()
			if err != nil {
				fmt.Printf("error from `PutPrice` method, package `controller`: %#v", err)
				errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
				if errAbort != nil {
					fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
				}
			}
		})

		selection.GET("manufacturers", func(ctx *gin.Context) {
			registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).ChooseManufacturers()
		})

		selection.POST("manufacturers", func(ctx *gin.Context) {
			err := registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).PutManufacturers()
			if err != nil {
				fmt.Printf("error from `PutManufacturers` method, package `controller`: %#v", err)
				errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
				if errAbort != nil {
					fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
				}
			}
		})

		selection.GET("choice", func(ctx *gin.Context) {
			registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).ChooseSource()
		})

		selection.POST("internet", func(ctx *gin.Context) {
			err := registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).GetSelectionFromInternetCars()
			if err != nil {
				fmt.Printf("error from `GetSelectionFromInternetCars` method, package `controller`: %#v", err)
				errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
				if errAbort != nil {
					fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
				}
			}
		})

		selection.GET("internet", func(ctx *gin.Context) {
			ServeSelectionCarList(ctx, redisSelectionDB, vehiclesDB, true)
		})

		selection.POST("internal_db", func(ctx *gin.Context) {
			err := registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).GetSelectionFromDBCars()
			if err != nil {
				fmt.Printf("error from `GetSelectionFromDBCars` method, package `controller`: %#v", err)
				errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
				if errAbort != nil {
					fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
				}
			}
		})

		selection.GET("internal_db", func(ctx *gin.Context) {
			ServeSelectionCarList(ctx, redisSelectionDB, vehiclesDB, false)
		})
	}
}

func ServeSelectionCarList(ctx *gin.Context, redisSelectionDB *redis.Client, vehiclesDB *sql.DB, choice bool) {
	sessionID := ctx.Query("guest")
	thisCarID := ctx.Query("carID")
	if thisCarID != "" {
		carID, err := strconv.Atoi(thisCarID)
		if err != nil {
			fmt.Printf("error from `Atoi` function, package `strconv`: %#v", err)
		}
		err = registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).DisplaySelectionCarAd(sessionID, carID, choice)
		if err != nil {
			fmt.Printf("error from `DisplaySelectionCarAd` method, package `controller`: %#v", err)
			errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
			if errAbort != nil {
				fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
			}
		}

	} else {
		err := registry.NewSelectionController(ctx, redisSelectionDB, vehiclesDB).TransferSelectionCarsData(sessionID, choice)
		if err != nil {
			fmt.Printf("error from `TransferSelectionCarsData` method, package `controller`: %#v", err)
			errAbort := ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("internal Server Error"))
			if errAbort != nil {
				fmt.Printf("error from `AbortWithError` method, package `gin`: %#v", err)
			}
		}
	}
}
