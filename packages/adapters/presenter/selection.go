package presenter

import (
	"fmt"
	"net/http"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	usecase "vehicles/packages/usecases/usecases"

	"github.com/gin-gonic/gin"
)

type selectionPresenter struct {
	// ctx - переменная контекста
	ctx adapters.Context
}

func NewSelectionPresenter(ctx adapters.Context) usecase.SelectionOutput {
	return &selectionPresenter{ctx}
}

// следующие методы рендерят различные веб-страницы
func (s *selectionPresenter) ShowPriorities() {
	s.ctx.HTML(http.StatusOK, "priorities.html", nil)
}

func (s *selectionPresenter) ShowPrice() {
	s.ctx.HTML(http.StatusOK, "price.html", nil)
}

func (s *selectionPresenter) ShowManufacturers() {
	s.ctx.HTML(http.StatusOK, "manufacturers.html", nil)
}

func (s *selectionPresenter) ShowSources() {
	s.ctx.HTML(http.StatusOK, "choice.html", nil)
}

// ShowResultOfFuzzyAlgorithm рендерит страницу, отображающую ранжированный с помощью нечеткого алгоритма список автомобилей
// Входные параметры: sessionID - идентификатор сессии, cars - автомобили
func (s *selectionPresenter) ShowResultOfFuzzyAlgorithm(sessionID string, cars []models.Car, choice bool) {
	indexes := make([]int, len(cars))
	for i := range cars {
		indexes[i] = i + 1
	}

	var Link string
	if choice {
		Link = fmt.Sprintf("http://localhost:8080/selection/internet?guest=%s&carID=", sessionID)
	} else {
		Link = fmt.Sprintf("http://localhost:8080/selection/internal_db?guest=%s&carID=", sessionID)
	}

	s.ctx.HTML(http.StatusOK, "offer_for_selection.html", gin.H{
		"Cars": cars, "Quantity": len(cars), "SessionID": sessionID, "Indexes": indexes, "Link": Link})
}

// ShowSelectionCarAd рендерит страницу конкретного автомобиля
// Входные параметры: sessionID - идентификатор сессии, car - автомобиль
func (s *selectionPresenter) ShowSelectionCarAd(sessionID string, car models.Car, choice bool) {
	var partOfLink string
	if choice {
		partOfLink = fmt.Sprintf("selection/internet?guest=%s", sessionID)
	} else {
		partOfLink = fmt.Sprintf("selection/internal_db?guest=%s", sessionID)
	}
	s.ctx.HTML(http.StatusOK, "car_card.html", gin.H{"Car": car, "PartOfLink": partOfLink})
}
