package repository

import "vehicles/packages/domain/models"

type CarsRepository interface {
	// LoadCarsData загружает в БД под управлением Redis данные об автомобилях, полученные из реляционной БД под управлением PostgreSQL
	// или собранные из интернета
	// Входные параметры: sessionID - идентификатор сессии, cars - автомобили
	LoadCarsData(sessionID string, cars []models.Car) error

	// GetCarsData получает из БД под управлением Redis ранее полученные данные из реляционной БД под управлением PostgreSQL или
	// собранные из интернета данные об автомобилях
	// Входные параметры: sessionID - идентификатор сессии
	GetCarsData(sessionID string) ([]models.Car, error)
}
