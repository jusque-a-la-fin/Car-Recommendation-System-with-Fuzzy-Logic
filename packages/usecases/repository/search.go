package repository

import (
	"vehicles/packages/domain/models"
)

type SearchRepository interface {
	// ScrapeSearchCars собирает данные автомобилей из интернета
	// Входной параметр: search - параметры поиска пользователя, которые он вводил на главное странице в большой форме сверху
	ScrapeSearchCars(search models.Search) ([]models.Car, error)
}
