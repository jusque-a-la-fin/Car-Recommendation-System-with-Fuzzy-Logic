package controller

import (
	"fmt"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	usecase "vehicles/packages/usecases/usecases"

	"github.com/redis/go-redis/v9"
)

type searchController struct {
	ctx           adapters.Context
	rdb           *redis.Client
	searchUseCase usecase.SearchInput
}

// Search содержит методы для обслуживания сервиса,
// осуществляющего обычный поиск без ранжирования и собирающего
// ответы от пользователей для "обучения" нечеткого алгоритма
type Search interface {
	GetSeachCars() error
	TransferSearchCarsData(sessionID string) error
	DisplaySearchCarAd(sessionID string, carID int) error
	DisplayMainPage()
}

func NewSearchController(ctx adapters.Context, rdb *redis.Client, sri usecase.SearchInput) Search {
	return &searchController{ctx, rdb, sri}
}

// GetSeachCars ответственен за получение списка автомобилей, чьи данные
// собраны из интернета, и сохранение его в БД под управлением Redis
func (src *searchController) GetSeachCars() error {
	type search struct {
		SessionID string        `json:"sessionID"`
		Form      models.Search `json:"form"`
	}
	srch := new(search)
	if err := src.ctx.Bind(&srch); err != nil {
		return fmt.Errorf("error from `Bind` method, package `gin`: %#v", err)
	}

	err := src.searchUseCase.GetCars(srch.Form, srch.SessionID)
	if err != nil {
		return fmt.Errorf("error from `GetCars` method, package `usecase`: %#v", err)
	}
	return nil
}

// TransferSearchCarsData ответственен за формирование веб-страницы, отображающей список автомобилей с вопросом,
// ответ на который помогает "нечеткому алгоритму" обучиться
func (src *searchController) TransferSearchCarsData(sessionID string) error {
	err := src.searchUseCase.PassSearchCarsData(sessionID)
	if err != nil {
		return fmt.Errorf("error from `PassSearchCarsData` method, package `usecase`: %#v", err)
	}

	return nil
}

// DisplaySearchCarAd ответственен за формирование веб-страницы конкретного автомобиля
// Входные параметры: sessionID - идентификатор сессии, carID - идентификатор автомобиля
func (src *searchController) DisplaySearchCarAd(sessionID string, carID int) error {
	err := src.searchUseCase.PresentSearchCarAd(sessionID, carID)
	if err != nil {
		return fmt.Errorf("error from `PresentSearchCarAd` method, package `usecase`: %#v", err)
	}
	return nil
}

// DisplayMainPage ответственен за формирование главной веб-страницы
func (src *searchController) DisplayMainPage() {
	src.searchUseCase.PresentMainPage()
}
