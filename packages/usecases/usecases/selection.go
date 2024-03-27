package usecase

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecases/repository"
)

// количество анализируемых автомобилей из интернета
const NumberOfCars = 10

// SelectionInput содержит методы, которые обслуживают сервис,
// использующий нечеткий алгоритм для ранжирования автомобилей
type SelectionInput interface {
	PickPriorities()
	SelectPriorities(sessionID string, priorities []string)
	PickPrice()
	SelectPrice(minPrice, maxPrice string)
	PickManufacturers()
	SelectManufacturers(manufacturers []string)
	PickSource()
	MakeSelectionFromDBCars(sessionID string) error
	MakeSelectionFromInternetCars(sessionID string) error
	PassSelectionCarsData(sessionID string, choice bool) error
	PresentSelectionCarAd(sessionID string, carID int, choice bool) error
}

// SelectionOutput содержит методы, которые рендерят html-шаблоны
type SelectionOutput interface {
	ShowPriorities()
	ShowPrice()
	ShowManufacturers()
	ShowSources()
	ShowResultOfFuzzyAlgorithm(sessionID string, cars []models.Car, choice bool)
	ShowSelectionCarAd(sessionID string, car models.Car, choice bool)
}

type selectionUseCase struct {
	ctx           adapters.Context
	selectionRepo repository.SelectionRepository
	carsRepo      repository.CarsRepository
	userUseCase   UserInput
	output        SelectionOutput
	User          models.User
}

func NewSelectionUseCase(ctx adapters.Context, sr repository.SelectionRepository, cr repository.CarsRepository, ut UserInput, ot SelectionOutput, ur models.User) SelectionInput {
	return &selectionUseCase{ctx, sr, cr, ut, ot, ur}
}

// PickPriorities ответственен за формирование веб-страницы, предлагающей пользователю
// расставить приоритеты: "Комфорт", "Экономичность", "Безопасность", "Динамика", "Управляемость", по
// которым будут ранжироваться автомобили
func (slu *selectionUseCase) PickPriorities() {
	slu.output.ShowPriorities()
}

// SelectPriorities ответственен за сохранение приоритетов, расставленных пользователем, в cookie
// Входной параметр: priorities - приоритеты, расставленные пользователем
func (slu *selectionUseCase) SelectPriorities(sessionID string, priorities []string) {
	slu.selectionRepo.SetPriorities(priorities)
}

// PickPrice ответственен за формирование веб-страницы, предлагающей пользователю задать
// минимальную и максимальную цены автомобилей
func (slu *selectionUseCase) PickPrice() {
	slu.output.ShowPrice()
}

// SelectPrice ответственен за сохранение минимальной и максимальной цен, заданных пользователем, в cookies
// Входные параметры: minPrice - минимальная цена, maxPrice - максимальная цена
func (slu *selectionUseCase) SelectPrice(minPrice, maxPrice string) {
	slu.selectionRepo.SetPrice(minPrice, maxPrice)
}

// PickManufacturers ответственен за формирование веб-страницы, предлагающей пользователю
// выбрать страны-производители автомобилей
func (slu *selectionUseCase) PickManufacturers() {
	slu.output.ShowManufacturers()
}

// SelectManufacturers ответственен за сохранение названий стран-производителей, выбранных
// пользователем, в cookies
// Входной параметр: manufacturers - страны, выбранные пользователем
func (slu *selectionUseCase) SelectManufacturers(manufacturers []string) {
	slu.selectionRepo.SetManufacturers(manufacturers)
}

// PickSource ответственен за формирование веб-страницы, предлагающей пользователю выбрать,
// из какого источника: базы данных или интернета он хочет получить ранжированный по его предпочтениям
// список автомобилей
func (slu *selectionUseCase) PickSource() {
	slu.output.ShowSources()
}

// MakeSelectionFromDBCars ответственен за получение списка автомобилей из реляционной БД,
// его ранжирование и сохранение в БД под управлением Redis
// Входной параметр: sessionID - идентификатор сессии
func (slu *selectionUseCase) MakeSelectionFromDBCars(sessionID string) error {
	selection, err := slu.selectionRepo.GetSelectionParams()
	if err != nil {
		return fmt.Errorf("error from `GetSelectionParams` method, package `gateway`: %#v", err)
	}

	cars, err := slu.selectionRepo.SelectCars(*selection)
	if err != nil {
		return fmt.Errorf("error from `SelectCars` method, package `gateway`: %#v", err)
	}

	ids, err := generateResultOfFuzzyAlgorithm(cars, selection.Priorities)
	if err != nil {
		return fmt.Errorf("error from `generateResultOfFuzzyAlgorithm` function, package `usecase`: %#v", err)
	}

	indexMap := make(map[int]models.Car)
	for _, car := range cars {
		indexMap[car.ID] = car
	}

	var sortedCars = make([]models.Car, 0, len(ids))
	for _, id := range ids {
		sortedCars = append(sortedCars, indexMap[id])
	}

	err = slu.carsRepo.LoadCarsData(sessionID, sortedCars)
	if err != nil {
		return fmt.Errorf("error from `LoadDBCarsData` method, package `gateway`: %#v", err)
	}
	return nil
}

// MakeSelectionFromInternetCars ответственен за получение списка автомобилей из интернета,
// его ранжирование и сохранение в БД под управлением Redis
// Входной параметр: sessionID - идентификатор сессии
func (slu *selectionUseCase) MakeSelectionFromInternetCars(sessionID string) error {
	selection, err := slu.selectionRepo.GetSelectionParams()
	if err != nil {
		return fmt.Errorf("error from `GetSelectionParams` method, package `gateway`: %#v", err)
	}

	makes, err := chooseRandomMakes(selection.Manufacturers)
	if err != nil {
		return fmt.Errorf("error from `chooseRandomMakes` function, package `usecase`: %#v", err)
	}
	cars, err := slu.selectionRepo.ScrapeSelectionCars(selection.MinPrice, selection.MaxPrice, makes)
	if err != nil {
		return fmt.Errorf("error from `ScrapeSelectionCars` method, package `gateway`: %#v", err)
	}

	ids, err := generateResultOfFuzzyAlgorithm(cars, selection.Priorities)
	if err != nil {
		return fmt.Errorf("error from `generateResultOfFuzzyAlgorithm` function, package `usecase`: %#v", err)
	}

	cars = getCarsForRendering(cars, ids)
	err = slu.carsRepo.LoadCarsData(sessionID, cars)
	if err != nil {
		return fmt.Errorf("error from `LoadCarsData` method, package `gateway`: %#v", err)
	}
	return nil
}

// chooseRandomMakes ответственна за выбор рандомных марок из списка доступных и определение количества автомобилей для каждой марки
// Входной параметр: initialCountries - страны-производители, выбранные пользователем
func chooseRandomMakes(initialCountries []string) ([]models.Makes, error) {
	if len(initialCountries) == 0 {
		initialCountries = []string{"Германия", "Япония", "США", "Китай", "Южная_Корея", "Франция", "Великобритания", "Россия", "Другие"}
	}

	countries := chooseCountries(initialCountries)
	makes, err := chooseMakes(countries)
	if err != nil {
		return nil, fmt.Errorf("error from `chooseMakes` function, package `usecase`: %#v", err)
	}
	return makes, nil
}

// chooseCountries распределяет количество автомобилей на страну
// Входной параметр: countries - страны-производители, выбранные пользователем
func chooseCountries(countries []string) map[string]int {
	coutriesSet := make(map[string]int)
	lenCountries := len(countries)
	if NumberOfCars == lenCountries {
		coutriesSet = map[string]int{
			"Германия":       1,
			"Япония":         1,
			"США":            1,
			"Китай":          1,
			"Южная_Корея":    1,
			"Франция":        1,
			"Великобритания": 1,
			"Россия":         1,
			"Другие":         1,
		}
	}

	limit := 0
	if NumberOfCars < lenCountries {
		limit = NumberOfCars
	} else if NumberOfCars > lenCountries {
		limit = lenCountries
	}

	for len(coutriesSet) < limit {
		countryIndex := rand.Intn(len(countries))
		randomCountry := countries[countryIndex]
		_, value := coutriesSet[randomCountry]
		if !value {
			coutriesSet[randomCountry] = 1
		}

	}

	if NumberOfCars > lenCountries {
		for i := lenCountries; i < NumberOfCars; i++ {
			countryIndex := rand.Intn(len(countries))
			randomCountry := countries[countryIndex]
			coutriesSet[randomCountry]++
		}
	}
	return coutriesSet
}

// chooseMakes выбирает рандомные марки из списка доступных и определяет количество автомобилей для каждой марки
// Входной параметр: countries - хэш-таблица, где ключ - страна-производитель, выбранный пользователем,
// значение - количество автомобилей для этой страны
func chooseMakes(countries map[string]int) ([]models.Makes, error) {
	allMakes, err := getMakes()
	if err != nil {
		return nil, fmt.Errorf("error from `getMakes` function, package `usecase`: %#v", err)
	}

	var makes []models.Makes
	match := false
	for country, quantity := range countries {
		for quantity > 0 {
			match = false
			randomIndex := rand.Intn(len(allMakes[country]))
			for index, thisMake := range makes {
				if allMakes[country][randomIndex] == thisMake.Make {
					makes[index].NumberOfCars++
					match = true
					break
				}
			}

			if !match {
				numberOfCars := 1

				make := models.Makes{Make: allMakes[country][randomIndex], NumberOfCars: numberOfCars}
				makes = append(makes, make)
			}
			quantity--
		}
	}
	return makes, nil
}

// getMakes собирает из файла страны и марки, которые к ним относятся
func getMakes() (map[string][]string, error) {
	jsonData, err := os.ReadFile("../packages/usecases/usecases/carMakes.json")
	if err != nil {
		return nil, fmt.Errorf("error from `ReadFile` function, package `os`: %#v", err)
	}

	carMakes := make(map[string][]string)
	err = json.Unmarshal(jsonData, &carMakes)
	if err != nil {
		return nil, fmt.Errorf("error from `Unmarshal` function, package `json`: %#v", err)
	}
	return carMakes, nil
}

// getCarsForRendering получает набор ранжированных автомобилей на основании набора
// ранжированных идентификаторов
// Входные параметры: cars - набор автомобилей, ids - набор ранжированных
// идентификаторов автомобилей
func getCarsForRendering(cars []models.Car, ids []int) []models.Car {
	indexMap := make(map[int]models.Car)
	for _, car := range cars {
		indexMap[car.ID] = car
	}

	var sortedCars = make([]models.Car, 0, len(ids))
	for _, id := range ids {
		sortedCars = append(sortedCars, indexMap[id])
	}
	return sortedCars
}

// PassSelectionCarsData ответственен за формирование веб-страницы, отображающей ранжированный список автомобилей
// Входной параметр: sessionID - идентификатор сесии
func (slu *selectionUseCase) PassSelectionCarsData(sessionID string, choice bool) error {
	cars, err := slu.carsRepo.GetCarsData(sessionID)
	if err != nil {
		return fmt.Errorf("error from `GetCarsData` method, package `gateway`: %#v", err)
	}

	slu.output.ShowResultOfFuzzyAlgorithm(sessionID, cars, choice)
	return nil
}

// PresentSelectionCarAd ответственен за формирование веб-страницы конкретного автомобиля
// Входные параметры: sessionID - идентификатор сесии, carID - идентификатор автомобиля
func (slu *selectionUseCase) PresentSelectionCarAd(sessionID string, carID int, choice bool) error {
	cars, err := slu.carsRepo.GetCarsData(sessionID)
	if err != nil {
		return fmt.Errorf("error from `GetCarsData` method, package `gateway`: %#v", err)
	}
	slu.output.ShowSelectionCarAd(sessionID, cars[carID-1], choice)
	return nil
}
