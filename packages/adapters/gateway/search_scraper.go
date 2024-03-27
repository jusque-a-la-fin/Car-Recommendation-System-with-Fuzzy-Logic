package gateway

import (
	"fmt"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecases/repository"

	"github.com/redis/go-redis/v9"
)

// ограничение количества автомобилей
var limitValue int = 10

const newWord = "new"

type searchRepository struct {
	// ctx - переменная контекста
	ctx adapters.Context
	// rdb - клиент Redis для подключения к NoSQL БД, хранящей данные выбранных пользователем автомобилей
	rdb *redis.Client
}

func NewSearchRepository(ctx adapters.Context, rdb *redis.Client) repository.SearchRepository {
	return &searchRepository{ctx, rdb}
}

// ScrapeSearchCars собирает данные автомобилей из интернета
// Входной параметр: search - параметры поиска пользователя, которые он вводил на главное странице в большой форме сверху
func (srp *searchRepository) ScrapeSearchCars(search models.Search) ([]models.Car, error) {
	// ссылка на страницу, содержащую объявления об автомобилях конкретной марки с названиями и ценами
	var link = prepareLinkForSearch(search)

	// получение ссылок на страницы автомобилей, их названий и цен
	links, names, prices, err := scrapeLinksNamesPrices(link, limitValue)
	if err != nil {
		return nil, fmt.Errorf("error from `scrapeLinksNamesPrices` function, package `gateway`: %#v", err)
	}

	cars := make([]models.Car, 0, limitValue)
	for i, link := range links {
		car := models.NewCar()
		car.ID = i
		car.FullName = names[i]

		car.Offering.Price = fmt.Sprintf("%s₽", prices[i])

		err := scrapeCharacteristics(&car, link, names[i])
		if err != nil {
			return nil, fmt.Errorf("error from `scrapeCharacteristics` function, package `gateway`): %#v", err)
		}

		cars = append(cars, car)
	}
	return cars, nil
}

// prepareLinkForSearch формирует и возвращает ссылку на веб-страницу согласно параметрам запроса
// Входной параметр: search - параметры поиска пользователя, которые он вводил на главное странице в большой форме сверху
func prepareLinkForSearch(search models.Search) string {
	var link = "https://auto.drom.ru/"
	if search.Mark == "" {
		if search.IsNewCar == newWord {
			link = fmt.Sprintf("%s%s/", link, search.IsNewCar)
		}
		link = fmt.Sprintf("%sall/?", link)
	} else {
		link = fmt.Sprintf("%s%s/", link, search.Mark)
	}

	if search.Model == "" {
		if search.IsNewCar == newWord && search.Mark != "" {
			link = fmt.Sprintf("%s%s/all/?", link, search.IsNewCar)
		}
		if search.IsNewCar != newWord && search.Mark != "" {
			link = fmt.Sprintf("%s?", link)
		}
	} else {
		if search.IsNewCar == newWord {
			link = fmt.Sprintf("%s%s/%s/?", link, search.Model, search.IsNewCar)
		} else {
			link = fmt.Sprintf("%s%s/?", link, search.Model)
		}
	}

	if search.LowPriceLimit != "" {
		link = fmt.Sprintf("%sminprice=%s", link, search.LowPriceLimit)
	}

	if search.HighPriceLimit != "" {
		if search.LowPriceLimit != "" {
			link = fmt.Sprintf("%s&", link)
		}
		link = fmt.Sprintf("%smaxprice=%s", link, search.HighPriceLimit)
	}

	if search.EarliestYear != "" {
		if search.LowPriceLimit != "" || search.HighPriceLimit != "" {
			link = fmt.Sprintf("%s&", link)
		}
		link = fmt.Sprintf("%sminyear=%s", link, search.EarliestYear)
	}

	if search.LatestYear != "" {
		if search.LowPriceLimit != "" || search.HighPriceLimit != "" || search.EarliestYear != "" {
			link = fmt.Sprintf("%s&", link)
		}
		link = fmt.Sprintf("%smaxyear=%s", link, search.LatestYear)
	}

	link = fmt.Sprintf("%s%s", link, prepareSecondPartOfLink(search, link))
	link = fmt.Sprintf("%s%s", link, prepareThirdPartOfLink(search, link))
	return link
}

/* Далее идут функции, которые формируют ссылку на web-страницу согласно другим параметрам запроса*/
// Входной параметр: search - параметры поиска пользователя, которые он вводил на главное странице в большой форме сверху
// link - формируемая ссылка на страницу марки
func prepareSecondPartOfLink(search models.Search, link string) string {
	if search.Gearbox == "AT" {
		if search.LowPriceLimit != "" || search.HighPriceLimit != "" ||
			search.EarliestYear != "" || search.LatestYear != "" {

			link = fmt.Sprintf("%s&", link)
		}
		link = fmt.Sprintf("%stransmission[]=2&transmission[]=5", link)

	} else if search.Gearbox != "" {
		if search.LowPriceLimit != "" || search.HighPriceLimit != "" ||
			search.EarliestYear != "" || search.LatestYear != "" {

			link = fmt.Sprintf("%s&", link)
		}
		link = fmt.Sprintf("%stransmission[]=%s", link, search.Gearbox)
	}

	if search.Fuel != "" {
		if search.LowPriceLimit != "" || search.HighPriceLimit != "" ||
			search.EarliestYear != "" || search.LatestYear != "" || search.Gearbox != "" {

			link = fmt.Sprintf("%s&", link)
		}
		link = fmt.Sprintf("%sfueltype=%s", link, search.Fuel)
	}
	return link
}

func prepareThirdPartOfLink(search models.Search, link string) string {
	if search.Drive != "" {
		if search.LowPriceLimit != "" || search.HighPriceLimit != "" || search.EarliestYear != "" ||
			search.LatestYear != "" || search.Gearbox != "" || search.Fuel != "" {
			link = fmt.Sprintf("%s&", link)
		}
		link = fmt.Sprintf("%sprivod=%s", link, search.Drive)
	}

	if search.LowPriceLimit != "" || search.HighPriceLimit != "" || search.EarliestYear != "" ||
		search.LatestYear != "" || search.Gearbox != "" || search.Fuel != "" || search.Drive != "" {
		link = fmt.Sprintf("%s&", link)
	}
	link = fmt.Sprintf("%sph=1&unsold=1", link)
	return link
}
