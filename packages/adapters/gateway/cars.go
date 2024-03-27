package gateway

import (
	"encoding/json"
	"fmt"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecases/repository"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type carsRepository struct {
	// ctx - переменная контекста
	ctx adapters.Context
	// rdb - клиент Redis для подключения к NoSQL БД, хранящей данные выбранных пользователем автомобилей
	rdb *redis.Client
}

func NewCarsRepository(ctx adapters.Context, rdb *redis.Client) repository.CarsRepository {
	return &carsRepository{ctx, rdb}
}

// LoadCarsData загружает в БД под управлением Redis данные об автомобилях, полученные из реляционной БД под управлением PostgreSQL
// или собранные из интернета
// Входные параметры: sessionID - идентификатор сессии, cars - автомобили
func (slr *carsRepository) LoadCarsData(sessionID string, cars []models.Car) error {
	carsJSON, err := json.Marshal(cars)
	if err != nil {
		return fmt.Errorf("error from `Marshal` function, package `json`: %#v", err)
	}

	if err = slr.rdb.Set(slr.ctx.(*gin.Context), sessionID, string(carsJSON), 0).Err(); err != nil {
		return fmt.Errorf("error from `Set` method, package `redis`: %#v", err)
	}
	return nil
}

// GetCarsData получает из БД под управлением Redis ранее полученные данные из реляционной БД под управлением PostgreSQL или
// собранные из интернета данные об автомобилях
// Входные параметры: sessionID - идентификатор сессии
func (slr *carsRepository) GetCarsData(sessionID string) ([]models.Car, error) {
	carsJSON, err := slr.rdb.Get(slr.ctx.(*gin.Context), sessionID).Result()
	if err != nil {
		return nil, fmt.Errorf("error from `Get` method, package `redis`: %#v", err)
	}

	var cars []models.Car
	if err := json.Unmarshal([]byte(carsJSON), &cars); err != nil {
		return nil, fmt.Errorf("error from `Unmarshal` function, package `json`: %#v", err)
	}
	return cars, nil
}
