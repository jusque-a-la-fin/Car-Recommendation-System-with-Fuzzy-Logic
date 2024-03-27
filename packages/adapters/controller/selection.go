package controller

import (
	"fmt"
	"vehicles/packages/adapters"
	usecase "vehicles/packages/usecases/usecases"
)

type selectionController struct {
	ctx              adapters.Context
	selectionUseCase usecase.SelectionInput
}

// Selection содержит методы, которые обслуживают сервис,
// использующий нечеткий алгоритм для ранжирования автомобилей
type Selection interface {
	ChoosePriorities()
	PutPriorities() error
	ChoosePrice()
	PutPrice() error
	ChooseManufacturers()
	PutManufacturers() error
	ChooseSource()
	GetSelectionFromDBCars() error
	GetSelectionFromInternetCars() error
	DisplaySelectionCarAd(sessionID string, carID int, choice bool) error
	TransferSelectionCarsData(sessionID string, choice bool) error
}

func NewSelectionController(ctx adapters.Context, slu usecase.SelectionInput) Selection {
	return &selectionController{ctx, slu}
}

// ChoosePriorities ответственен за формирование веб-страницы, предлагающей пользователю
// расставить приоритеты: "Комфорт", "Экономичность", "Безопасность", "Динамика", "Управляемость", по
// которым будут ранжироваться автомобили
func (slc *selectionController) ChoosePriorities() {
	slc.selectionUseCase.PickPriorities()
}

// PutPriorities ответственен за сбор приоритетов, расставленных пользователем, и сохранение их в cookie
func (slc *selectionController) PutPriorities() error {
	type priorities struct {
		Priorities []string `json:"priorities"`
		SessionID  string   `json:"sessionID"`
	}

	prs := new(priorities)
	var err error
	if err = slc.ctx.BindJSON(&prs); err != nil {
		return fmt.Errorf("error from `BindJSON` method, package `gin`: %#v", err)
	}

	slc.selectionUseCase.SelectPriorities(prs.SessionID, prs.Priorities)
	return nil
}

// ChoosePrice ответственен за формирование веб-страницы, предлагающей пользователю задать
// минимальную и максимальную цены автомобилей
func (slc *selectionController) ChoosePrice() {
	slc.selectionUseCase.PickPrice()
}

// PutPrice ответственен за сбор минимальной и максимальной цен, заданных пользователем, и сохранение их в cookies
func (slc *selectionController) PutPrice() error {
	type price struct {
		MinPrice string `json:"minPrice"`
		MaxPrice string `json:"maxPrice"`
	}

	pc := new(price)
	var err error
	if err = slc.ctx.BindJSON(&pc); err != nil {
		return fmt.Errorf("error from `BindJSON` method, package `gin`: %#v", err)
	}

	slc.selectionUseCase.SelectPrice(pc.MinPrice, pc.MaxPrice)
	return nil
}

// ChooseManufacturers ответственен за формирование веб-страницы, предлагающей пользователю
// выбрать страны-производители автомобилей
func (slc *selectionController) ChooseManufacturers() {
	slc.selectionUseCase.PickManufacturers()
}

// PutManufacturers ответственен за сбор названий стран-производителей, выбранных
// пользователем, и их сохранение в cookie
func (slc *selectionController) PutManufacturers() error {
	type manufacturers struct {
		Manufacturers []string `json:"manufacturers"`
	}

	mns := new(manufacturers)
	if err := slc.ctx.BindJSON(&mns); err != nil {
		return fmt.Errorf("error from `BindJSON` method, package `gin`: %#v", err)
	}

	slc.selectionUseCase.SelectManufacturers(mns.Manufacturers)
	return nil
}

// ChooseSource ответственен за формирование веб-страницы, предлагающей пользователю выбрать,
// из какого источника: базы данных или интернета он хочет получить ранжированный по его предпочтениям
// список автомобилей
func (slc *selectionController) ChooseSource() {
	slc.selectionUseCase.PickSource()
}

// GetSelectionFromDBCars ответственен за получение ранжированного списка автомобилей, чьи данные собраны
// из реляционной БД, и сохранение его в БД под управлением Redis
func (slc *selectionController) GetSelectionFromDBCars() error {
	var sessionID string
	if err := slc.ctx.BindJSON(&sessionID); err != nil {
		return fmt.Errorf("error from `BindJSON` method, package `gin`: %#v", err)
	}

	err := slc.selectionUseCase.MakeSelectionFromDBCars(sessionID)
	if err != nil {
		return fmt.Errorf("error from `MakeSelectionFromDBCars` method, package `usecase`: %#v", err)
	}
	return nil
}

// GetSelectionFromInternetCars ответственен за получение ранжированного списка автомобилей, чьи данные
// собраны из интернета, и сохранение его в БД под управлением Redis
func (slc *selectionController) GetSelectionFromInternetCars() error {
	var sessionID string
	if err := slc.ctx.BindJSON(&sessionID); err != nil {
		return fmt.Errorf("error from `BindJSON` method, package `gin`: %#v", err)
	}

	err := slc.selectionUseCase.MakeSelectionFromInternetCars(sessionID)
	if err != nil {
		return fmt.Errorf("error from `MakeSelectionFromInternetCars` method, package `usecase`: %#v", err)
	}
	return nil
}

// TransferSelectionCarsData ответственен за формирование веб-страницы, отображающей ранжированный список автомобилей
// Входной параметр: sessionID - идентификатор сессии
func (slc *selectionController) TransferSelectionCarsData(sessionID string, choice bool) error {
	err := slc.selectionUseCase.PassSelectionCarsData(sessionID, choice)
	if err != nil {
		return fmt.Errorf("error from `PassSelectionCarsData` method, package `usecase`: %#v", err)
	}
	return nil
}

// DisplaySelectionCarAd ответственен за формирование веб-страницы конкретного автомобиля
// Входные параметры: sessionID - идентификатор сессии, carID - идентификатор автомобиля
func (slc *selectionController) DisplaySelectionCarAd(sessionID string, carID int, choice bool) error {
	err := slc.selectionUseCase.PresentSelectionCarAd(sessionID, carID, choice)
	if err != nil {
		return fmt.Errorf("error from `PresentSelectionCarAd` method, package `usecase`: %#v", err)
	}
	return nil
}
