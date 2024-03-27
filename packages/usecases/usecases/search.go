package usecase

import (
	"fmt"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecases/repository"
)

// SearchInput содержит методы для обслуживания сервиса,
// осуществляющего обычный поиск без ранжирования и собирающего
// ответы от пользователей для "обучения" нечеткого алгоритма
type SearchInput interface {
	GetCars(search models.Search, sessionID string) error
	PassSearchCarsData(sessionID string) error
	PresentSearchCarAd(sessionID string, id int) error
	PresentMainPage()
}

// SearchOutput содержит методы, которые рендерят html-шаблоны
type SearchOutput interface {
	ShowCarsWithSurvey(sessionID, question string, cars []models.Car, possibleAnswers []string)
	ShowCars(sessionID string, cars []models.Car)
	ShowSearchCarAd(sessionID string, car models.Car)
	ShowMainPage()
}

type searchUseCase struct {
	searchRepo      repository.SearchRepository
	carsRepo        repository.CarsRepository
	userUseCase     UserInput
	questionUseCase QuestionInput
	output          SearchOutput
}

func NewSearchUseCase(r repository.SearchRepository, cr repository.CarsRepository, u UserInput, q QuestionInput, o SearchOutput) SearchInput {
	return &searchUseCase{r, cr, u, q, o}
}

// GetCars ответственен за получение списка автомобилей, чьи данные
// собраны из интернета, и сохранение его в БД под управлением Redis
func (sru *searchUseCase) GetCars(search models.Search, sessionID string) error {
	cars, err := sru.searchRepo.ScrapeSearchCars(search)
	if err != nil {
		return fmt.Errorf("error from `ScrapeSearchCars` method, package `gateway`: %#v", err)
	}

	err = sru.carsRepo.LoadCarsData(sessionID, cars)
	if err != nil {
		return fmt.Errorf("error from `LoadCarsData` method, package `gateway`: %#v", err)
	}
	return nil
}

// PassSearchCarsData ответственен за формирование веб-страницы, отображающей список автомобилей с вопросом,
// ответ на который помогает "нечеткому алгоритму" обучиться
func (sru *searchUseCase) PassSearchCarsData(sessionID string) error {
	cars, err := sru.carsRepo.GetCarsData(sessionID)
	if err != nil {
		return fmt.Errorf("error from `GetCarsData` method, package `gateway`: %#v", err)
	}

	question, possibleAnswers, err := sru.questionUseCase.PickQuestion()
	if err != nil {
		return fmt.Errorf("error from `PickQuestion` method, package `usecase`: %#v", err)
	}

	sru.output.ShowCarsWithSurvey(sessionID, question, cars, possibleAnswers)
	return nil
}

// PresentSearchCarAd ответственен за формирование веб-страницы конкретного автомобиля
// Входные параметры: sessionID - идентификатор сессии, carID - идентификатор автомобиля
func (sru *searchUseCase) PresentSearchCarAd(sessionID string, carID int) error {
	cars, err := sru.carsRepo.GetCarsData(sessionID)
	if err != nil {
		return fmt.Errorf("error from `GetCarsData` method, package `gateway`: %#v", err)
	}

	sru.output.ShowSearchCarAd(sessionID, cars[carID-1])
	return nil
}

// PresentMainPage ответственен за формирование главной веб-страницы
func (sru *searchUseCase) PresentMainPage() {
	sru.output.ShowMainPage()
}
