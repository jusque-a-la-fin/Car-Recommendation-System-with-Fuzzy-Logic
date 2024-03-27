package presenter

import (
	"fmt"
	"net/http"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	usecase "vehicles/packages/usecases/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type searchPresenter struct {
	// ctx - переменная контекста
	ctx adapters.Context
}

func NewSearchPresenter(ctx adapters.Context) usecase.SearchOutput {
	return &searchPresenter{ctx}
}

// ShowCarsWithSurvey рендерит страницу, отображающую список автомобилей с вопросом для пользователя,
// ответ на который помогает нечеткому алгоритму "обучиться"
// Входные параметры: sessionID - идентификатор сессии, question - вопрос для пользователя,
// cars - автомобили, possibleAnswers - варианты ответа
func (s *searchPresenter) ShowCarsWithSurvey(sessionID, question string, cars []models.Car, possibleAnswers []string) {
	htmlFileName := "offer_for_search.html"

	indexes := make([]int, len(cars))
	for i := range cars {
		indexes[i] = i + 1
	}
	s.ctx.HTML(http.StatusOK, htmlFileName, gin.H{"Cars": cars, "Quantity": len(cars), "SessionID": sessionID,
		"Indexes": indexes, "NotAnswered": true, "Question": question, "PossibleAnswers": possibleAnswers})
}

// ShowCars рендерит страницу, отображающую список автомобилей без вопроса для пользователя
// Входные параметры: sessionID - идентификатор сессии, cars - автомобили
func (s *searchPresenter) ShowCars(sessionID string, cars []models.Car) {
	htmlFileName := "offer_for_search.html"
	indexes := make([]int, len(cars))
	for i := range cars {
		indexes[i] = i + 1
	}
	s.ctx.HTML(http.StatusOK, htmlFileName, gin.H{"Cars": cars, "Quantity": len(cars), "SessionID": sessionID, "Indexes": indexes})
}

// ShowSearchCarAd рендерит страницу конкретного автомобиля
// Входные параметры: sessionID - идентификатор сессии, car - автомобиль
func (s *searchPresenter) ShowSearchCarAd(sessionID string, car models.Car) {
	partOfLink := fmt.Sprintf("search?guest=%s", sessionID)
	s.ctx.HTML(http.StatusOK, "car_card.html", gin.H{"Car": car, "PartOfLink": partOfLink})
}

// ShowMainPage рендерит главную страницу
func (s *searchPresenter) ShowMainPage() {
	sessionID := uuid.New().String()
	htmlFileName := "main_page.html"
	s.ctx.HTML(http.StatusOK, htmlFileName, gin.H{"sessionID": sessionID})
}
