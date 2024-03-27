package repository

import (
	"vehicles/packages/domain/models"
)

type SelectionRepository interface {
	// SetPriorities сохраняет приоритеты, расставленные пользователем, в cookie
	// Приоритеты - это свойства: "Комфорт", "Экономичность", "Безопасность", "Динамика", "Управляемость",
	// по которым будут ранжироваться автомобили
	// Входной параметр: priorities - приоритеты, расставленные пользователем
	SetPriorities(priorities []string)

	// SetPrice сохраняет диапазон цен, заданный пользователем, в cookies
	// Входные параметры: minPrice - минимальная цена, maxPrice - максимальная цена
	SetPrice(minPrice, maxPrice string)

	// SetManufacturers сохраняет названия стран-производителей, выбранных пользователем, в cookie
	// Входной параметр: manufacturers - страны, выбранные пользователем
	SetManufacturers(manufacturers []string)

	// GetSelectionParams получает параметры, заданные ранее пользователем, из cookies
	GetSelectionParams() (*models.Selection, error)

	// SelectCars получает из реляционной БД под управлением PostgreSQL информацию об автомобилях
	// Входной параметр: sln - запрос пользователя
	SelectCars(slc models.Selection) ([]models.Car, error)

	// ScrapeSelectionCars собирает данные автомобилей из интернета
	// Входные параметры: minPrice  - минимальная цена, maxPrice - максимальная цена, makes - срез марок
	ScrapeSelectionCars(minPrice, maxPrice string, makes []models.Makes) ([]models.Car, error)
}
