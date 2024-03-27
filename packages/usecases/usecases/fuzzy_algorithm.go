package usecase

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"vehicles/packages/domain/models"
)

const (
	suspensionType1 = "Независимая, на двойных поперечных рычагах"
	suspensionType2 = "Многорычажная, независимая"
	suspensionType3 = "Пневматическая"
	suspensionType4 = "Независимая, амортизационная стойка типа МакФерсон"
	suspensionType5 = "Полузависимая, торсионная балка"
	suspensionType6 = "Зависимая, пружинная"
	suspensionType7 = "Листовая, пружинная"
	brakesType1     = "Дисковые"
	brakesType2     = "Дисковые вентилируемые"
	brakesType3     = "Барабанные"
)

type carRecommendation struct {
	// car - автомобиль
	car models.Car
	// recommendationValue - выходное значение нечеткого алгоритма
	recommendationValue float64
}

// generateResultOfFuzzyAlgorithm получает выходное значение нечеткого алгоритма для каждого автомобиля, ранжирует автомобили
// по убыванию выходного значения нечеткого алгоритма и возвращает срез из идентификаторов ранжированных автомобилей
// Входные параметры: cars - автомобили, priorities - приоритеты, расставленные пользователем
func generateResultOfFuzzyAlgorithm(cars []models.Car, priorities []string) ([]int, error) {
	ids := make([]int, len(cars))
	if len(priorities) == 0 {
		var errFlag error
		sort.Slice(cars, func(i, j int) bool {
			price1, err := strconv.Atoi(cars[i].Offering.Price)
			if err != nil {
				errFlag = fmt.Errorf("can't get price1: error from `Atoi` function, package `strconv`: %#v", err)
			}
			price2, err := strconv.Atoi(cars[j].Offering.Price)
			if err != nil {
				errFlag = fmt.Errorf("can't get price2: error from `Atoi` function, package `strconv`: %#v", err)
			}
			return price1 < price2
		})

		if errFlag != nil {
			return nil, errFlag
		}

		for idx := 0; idx < len(cars); idx++ {
			ids[idx] = cars[idx].ID
		}
		return ids, nil

	} else {
		carRecs := make([]carRecommendation, len(cars))
		for idx := 0; idx < len(cars); idx++ {
			value, err := performFuzzyAlgorithm(cars[idx], priorities)
			if err != nil {
				return nil, fmt.Errorf("error from `performFuzzyAlgorithm` function, package `usecase`: %#v", err)
			}
			carRecs[idx] = carRecommendation{car: cars[idx], recommendationValue: value}
		}

		sort.Slice(carRecs, func(idx, jdx int) bool {
			return carRecs[idx].recommendationValue > carRecs[jdx].recommendationValue
		})

		for idx := 0; idx < len(carRecs); idx++ {
			ids[idx] = carRecs[idx].car.ID
		}
		return ids, nil
	}
}

// performFuzzyAlgorithm выполняет нечеткий алгоритм(получает выходное значение нечеткого алгоритма) для одного автомобиля
// Входные параметры: car - автомобиль, priorities - приоритеты, расставленные пользователем
func performFuzzyAlgorithm(car models.Car, priorities []string) (float64, error) {
	// filePriorities - имя файла, который содержит все возможные расстановки приоритетов.
	// Набор этих расстановок является "размещением" (термин комбинаторики)
	filePriorities := "/home/daniel/Desktop/INDEPENDENCY/Fuzzy-Logic-Recommendation-System 4.8/packages/usecases/usecase/priorities.txt"
	file, err := os.Open(filePriorities)
	if err != nil {
		return -1, fmt.Errorf("error from `Open` function, package `os`: %#v", err)
	}
	defer file.Close()

	// rulesFilesMap содержит в качестве ключей - условие "ЕСЛИ", например,
	// "экономичность низкий безопасность низкий динамика высокий управляемость высокий",
	// а в качестве значения - значение, которое определяет, насколько сильно
	// будет рекомендоваться автомобиль, например 1,2,3, и т.д. Это значение принадлежит нечеткому множеству "рекомендация"
	rulesFilesMap := make(map[string]int)
	count := 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		rulesFilesMap[line] = count
		count++
	}

	if err = scanner.Err(); err != nil {
		return -1, fmt.Errorf("error from `Err` method, package `bufio`: %#v", err)
	}

	prioritiesStr := fmt.Sprintf("%v", priorities)
	runes := []rune(prioritiesStr)
	prioritiesSubstr := ""
	if len(runes) > 2 {
		substrRunes := runes[1 : len(runes)-1]
		prioritiesSubstr = string(substrRunes)
	}

	if prioritiesSubstr == "" {
		return -1, fmt.Errorf("error deleting square brackets [] from a string containing priorities set by the user")
	}

	// fileRules - название файла, содежащего нечеткие правила
	fileRules := fmt.Sprintf("/home/daniel/Desktop/INDEPENDENCY/Fuzzy-Logic-Recommendation-System 4.8/packages/usecases/usecase/rules/%d_rules.txt",
		rulesFilesMap[prioritiesSubstr])
	file, err = os.Open(fileRules)
	if err != nil {
		return -1, fmt.Errorf("error from `Open` function, package `os`: %#v", err)
	}
	defer file.Close()

	var currentSet []string
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		currentSet = append(currentSet, line)
	}

	if err := scanner.Err(); err != nil {
		return -1, fmt.Errorf("error from `Err` method, package `bufio`: %#v", err)
	}

	// rules - нечеткие правила. каждое правило хранится в срезе
	rules := make([][]string, 0, len(currentSet))
	var result []string
	for _, rule := range currentSet {
		parts := strings.Split(rule, " ")
		for idx := 0; idx < len(parts); idx += 2 {
			if idx+1 < len(parts) {
				result = append(result, parts[idx]+" "+parts[idx+1])
			} else {
				result = append(result, parts[idx])
			}
		}
		rules = append(rules, append([]string{}, result...))
		result = nil
	}

	// valuesRecommendation - значения, которые будут определять, насколько сильно
	// будет рекомендоваться автомобиль. Эти значения принадлежат нечеткому множеству "рекомендация"
	var valuesRecommendation = make([]int, 0, len(rules))
	for _, rule := range rules {
		value, err := strconv.Atoi(rule[len(rule)-1])
		if err != nil {
			return -1, fmt.Errorf("error from `Atoi` function, package `strconv`: %#v", err)
		}
		valuesRecommendation = append(valuesRecommendation, value)
	}

	// fuelConsumption - коэффициент экономичности или расход топлива в смешанном цикле в литрах на 100 км
	fuelConsumption := car.Specs.MixedFuelConsumption

	// timeOfAcceleration0To100kmh - коэффициент динамики или время разгона в секундах до 100 км/ч
	timeOfAcceleration0To100kmh := car.Specs.Acceleration0To100

	// handlingCoeff - коэффициент управляемости
	handlingCoeff := calculateHandlingCoefficient(car.Specs.Engine.MaxPower, car.Specs.FrontTrackWidth, car.Specs.BackTrackWidth,
		car.Specs.Drive, car.Specs.Suspension, car.Specs.Tires, car.Features.SafetyAndMotionControlSystem.ABS,
		car.Features.SafetyAndMotionControlSystem.ESP, car.Features.SafetyAndMotionControlSystem.EBD,
		car.Features.SafetyAndMotionControlSystem.BAS, car.Features.SafetyAndMotionControlSystem.TCS, car.Specs.Brakes.FrontBrakes,
		car.Specs.Brakes.BackBrakes, car.Specs.Mass, car.Specs.Wheelbase, car.Specs.Length, car.Specs.Width,
		car.Specs.Height, car.Specs.GroundClearance, car.Specs.DragCoefficient,
	)

	// comfortCoeff - коэффициент комфорта
	comfortCoeff := calculateComfortCoefficient(car.Specs.Suspension, car.Specs.Gearbox, car.Features.CabinMicroclimate,
		car.Features.Interior, car.Features.ElectricOptions, car.Features.MultimediaSystems, car.Features.Lights,
		car.Specs.SteeringWheel.PowerSteering, car.Features.CarAlarm, car.Specs.TrunkVolume)

	// safetyCoeff - коэффициент безопасности
	safetyCoeff := calculateSafetyCoefficient(car.Specs.CrashTestEstimate, car.Features.SafetyAndMotionControlSystem,
		car.Features.Airbags, car.Specs.Brakes)

	// valuesOfMemebershipFunction - значения функций принадлежности, соответствующих нечетким подмножествам "низкий", "средний" и
	// "высокий" нечетких множеств "экономичность", "динамика", "управляемость", "комфорт", "безопасность"
	valuesOfMemebershipFunction := make([][]float64, 0, len(rules))
	for _, rule := range rules {
		values := []float64{}
		for _, set := range rule {
			parts := strings.Split(set, " ")
			newValues := calculateMembershipFunctionValues(parts, safetyCoeff, handlingCoeff, comfortCoeff, timeOfAcceleration0To100kmh, fuelConsumption)
			values = append(values, newValues...)
		}
		valuesOfMemebershipFunction = append(valuesOfMemebershipFunction, values)
	}

	// minValuesOfMemebershipFunction - глобальные максимумы функций принадлежности нечетких подмножеств множества "рекомендация"
	// или ординаты вершин треугольников, которые образуются под графиками этих функций
	var minValuesOfMemebershipFunction = make([]float64, 0, len(valuesOfMemebershipFunction))
	for _, values := range valuesOfMemebershipFunction {
		minValuesOfMemebershipFunction = append(minValuesOfMemebershipFunction, findMin(values))
	}

	// recommendationValue - значение рекомендации и выходное значение нечеткого алгоритма для одного автомобиля
	var recommendationValue float64
	if len(minValuesOfMemebershipFunction) == len(valuesRecommendation) {
		recommendationValue = defuzzyficate(minValuesOfMemebershipFunction, valuesRecommendation)
	} else {
		return -1, fmt.Errorf("error, the number of degrees of recommendation should be equal to the number of rules")
	}

	return recommendationValue, nil
}

// calculateMembershipFunctionValues вычисляет значения функций принадлежности, соответствующих одному из нечетких подмножеств "низкий", "средний" и
// "высокий" всех нечетких множеств "экономичность", "динамика", "управляемость", "комфорт", "безопасность"
// Входные параметры: parts - срез, содержащий название нечеткого множества и название нечеткого подмножества из нечеткого правила
// safetyCoeff - коэффициент безопасности, handlingCoeff - коэффициент управляемости, comfortCoeff - коэффициент комфорта,
// dynamicsCoeff - коэффициент динамики, economyCoeff - коэффициент экономичности
func calculateMembershipFunctionValues(parts []string, safetyCoeff, handlingCoeff, comfortCoeff, dynamicsCoeff, economyCoeff float64) []float64 {
	values := []float64{}
	switch parts[0] {
	case "безопасность":
		if safetyCoeff == 0 {
			break
		}
		params := provideFunctionParameters(parts[0], parts[1])
		safety := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], safetyCoeff)
		values = append(values, safety)

	case "управляемость":
		if handlingCoeff == 0 {
			break
		}
		params := provideFunctionParameters(parts[0], parts[1])
		handling := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], handlingCoeff)
		values = append(values, handling)

	case "комфорт":
		if comfortCoeff == 0 {
			break
		}
		params := provideFunctionParameters(parts[0], parts[1])
		comfort := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], comfortCoeff)
		values = append(values, comfort)

	case "динамика":
		if dynamicsCoeff == 0 {
			break
		}
		params := provideFunctionParameters(parts[0], parts[1])
		dynamics := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], dynamicsCoeff)
		values = append(values, dynamics)

	case "экономичность":
		if economyCoeff == 0 {
			break
		}
		params := provideFunctionParameters(parts[0], parts[1])
		economy := calculateMembershipFunctionValueForLinguisticVariables(params, parts[1], economyCoeff)
		values = append(values, economy)
	}
	return values
}

// calculateHandlingCoefficient вычисляет коэффициент управляемости
// Входные параметры: power - мощность двигателя в л.с, frontTrackWidth - ширина передней колеи в мм,
// backTrackWidth - ширина задней колеи в мм, drive - тип привода, sps - информация о подвеске,
// trs - информация о шинах, abs, esp, ebd, bas, tcs - наличие систем соответственно
// ABS, ESP, EBD, BAS, TCS, frontBrakes - тип передних тормозов, backBrakes - тип задних тормозов,
// mass - масса в кг, wheelbase - колесная база в мм, length - длина в мм, width - ширина в мм, height - высота в мм,
// groundClearance - клиренс в мм, dragCoefficient - коэффициент лобового сопротивления
func calculateHandlingCoefficient(power, frontTrackWidth, backTrackWidth float64, drive string, sps models.Suspension,
	trs models.Tires, abs, esp, ebd, bas, tcs models.Availability, frontBrakes, backBrakes string, mass, wheelbase, length,
	width, height, groundClearance, dragCoefficient float64) float64 {

	// перевод из лошадиных сил в ватты
	var newPower = power * 735.5

	// перевод из мм в метры
	frontTrackWidth /= 1000
	backTrackWidth /= 1000
	frontTiresWidth := float64(trs.FrontTiresWidth) / 1000
	backTiresWidth := float64(trs.BackTiresWidth) / 1000
	wheelbase /= 1000
	length /= 1000
	width /= 1000
	height /= 1000
	groundClearance /= 1000

	// driveTypeCoefficient - коэффициент типа привода
	driveTypeCoefficient := calculateDriveTypeCoefficient(drive)

	// коэффициент наличия переднего стабилизатора
	var frontStabilizerCoefficient = 1.0

	// коэффициент наличия заднего стабилизатора
	var backStabilizerCoefficient = 1.0

	if sps.FrontStabilizer == models.YesValue {
		frontStabilizerCoefficient = 1.2
	}

	if sps.BackStabilizer == models.YesValue {
		backStabilizerCoefficient = 1.2
	}

	// frontSuspensionCoefficient - коэффициент типа передней подвески,
	// backSuspensionCoefficient - коэффициент типа задней подвески
	frontSuspensionCoefficient, backSuspensionCoefficient := calculateSuspensionCoeffsForHandlingCoeff(sps.FrontSuspension,
		sps.BackSuspension)

	frontTiresDiameter, backTiresDiameter := calculateTiresParamsAndTrackWidths(&frontTrackWidth, &backTrackWidth, &frontTiresWidth, &backTiresWidth,
		&trs.FrontTiresAspectRatio, &trs.BackTiresAspectRatio, trs.FrontTiresRimDiameter, trs.BackTiresRimDiameter)

	var absCoefficient float64 = 0
	var espCoefficient float64 = 0
	var ebdCoefficient float64 = 0
	var basCoefficient float64 = 0
	var tcsCoefficient float64 = 0

	if abs == models.YesValue {
		absCoefficient = 0.064
	}

	if esp == models.YesValue {
		espCoefficient = 0.07
	}

	if ebd == models.YesValue {
		ebdCoefficient = 0.056
	}

	if bas == models.YesValue {
		basCoefficient = 0.059
	}

	if tcs == models.YesValue {
		tcsCoefficient = 0.051
	}

	// коэффициент типа передних тормозов
	var frontBrakesCoefficient float64
	// коэффициент типа задних тормозов
	var backBrakesCoefficient float64

	switch frontBrakes {
	case brakesType1, brakesType2:
		frontBrakesCoefficient = 0.7
	case brakesType3:
		frontBrakesCoefficient = 0.5
	default:
		frontBrakesCoefficient = 0
	}

	switch backBrakes {
	case brakesType1, brakesType2:
		backBrakesCoefficient = 0.6
	case brakesType3:
		backBrakesCoefficient = 0.4
	default:
		backBrakesCoefficient = 0
	}

	// efficientFrontTrackWidth - оптимальная ширина передней колеи в метрах
	efficientFrontTrackWidth := frontTrackWidth + 0.5*math.Abs(frontTiresWidth-backTiresWidth)/float64(trs.FrontTiresAspectRatio)

	// efficientBackTrackWidth - оптимальная ширина задней колеи в метрах
	efficientBackTrackWidth := backTrackWidth + 0.5*math.Abs(backTiresWidth-frontTiresWidth)/float64(trs.BackTiresAspectRatio)

	// числитель
	numerator := ((newPower * (efficientFrontTrackWidth + efficientBackTrackWidth) / 2) * driveTypeCoefficient *
		(frontSuspensionCoefficient*frontStabilizerCoefficient + backSuspensionCoefficient*backStabilizerCoefficient) *
		(frontTiresWidth*frontTiresDiameter + backTiresWidth*backTiresDiameter) * (frontBrakesCoefficient + backBrakesCoefficient))

	// знаменатель
	denominator := mass * wheelbase * (length + width + height) * groundClearance * dragCoefficient

	if mass == 0 || wheelbase == 0 || length == 0 || width == 0 || height == 0 || groundClearance == 0 {
		denominator = 1
	}

	// коэффициент управляемости
	handlingCoefficient := numerator / denominator

	var sizeCoefficient float64
	if handlingCoefficient != 0 {
		sizeCoefficient = 30.0
	}

	// коэффициент наличия систем безопасности
	controlSystemsCoefficient := absCoefficient + espCoefficient + ebdCoefficient + basCoefficient + tcsCoefficient
	handlingCoefficient += handlingCoefficient * controlSystemsCoefficient
	handlingCoefficient = math.Abs(handlingCoefficient - sizeCoefficient)
	return handlingCoefficient
}

// calculateDriveTypeCoefficient вычисляет коэффициент типа привода
// Входной параметр: driveType - тип привода
func calculateDriveTypeCoefficient(driveType string) float64 {
	switch driveType {
	case "Передний(FF)", "Передний":
		return 0.9
	case "Полный (4WD)", "Полный":
		return 1
	case "Задний(FR)", "Задний":
		return 0.7
	default:
		return 0
	}
}

// calculateSuspensionCoeffsForHandlingCoeff вычисляет коэффициент типа передней подвески и коэффициент типа задней подвески
// Входные параметры: frontSuspension - тип передней подвески, backSuspension - тип задней подвески
func calculateSuspensionCoeffsForHandlingCoeff(frontSuspension, backSuspension string) (float64, float64) {
	// коэффициент типа передней подвески
	var frontSuspensionCoefficient = 1.0
	// коэффициент типа задней подвески
	var backSuspensionCoefficient = 1.0

	switch frontSuspension {
	case suspensionType1:
		frontSuspensionCoefficient = 1.9

	case suspensionType2:
		frontSuspensionCoefficient = 1.8

	case suspensionType3:
		frontSuspensionCoefficient = 1.7

	case suspensionType4:
		frontSuspensionCoefficient = 1.6

	case suspensionType5:
		frontSuspensionCoefficient = 1.5

	case suspensionType6:
		frontSuspensionCoefficient = 1.4

	case suspensionType7:
		frontSuspensionCoefficient = 1.3
	}

	switch backSuspension {
	case suspensionType2:
		backSuspensionCoefficient = 1.9

	case suspensionType1:
		backSuspensionCoefficient = 1.8

	case suspensionType3:
		backSuspensionCoefficient = 1.7

	case suspensionType4:
		backSuspensionCoefficient = 1.6

	case suspensionType5:
		backSuspensionCoefficient = 1.4

	case suspensionType6:
		backSuspensionCoefficient = 1.3

	case suspensionType7:
		backSuspensionCoefficient = 1.2
	}
	return frontSuspensionCoefficient, backSuspensionCoefficient
}

// calculateTiresParamsAndTrackWidths работает с параметрами шин и значениями передней колеи и задней колеи
// Входные параметры: frontTiresWidth - ширина передних шин в метрах, backTiresWidth - ширина задних шин в метрах,
// frontTiresAspectRatio - процентное соотношение высоты профиля передних шин к их ширине,
// backTiresAspectRatio - процентное соотношение высоты профиля задних шин к их ширине,
// frontTiresRimDiameter - диаметр обода  передних шин в дюймах,
// backTiresRimDiameter - диаметр обода задних шин в дюймах
// frontTrackWidth - ширина передней колеи в метрах, backTrackWidth - ширина задней колеи в метрах,
func calculateTiresParamsAndTrackWidths(frontTrackWidth, backTrackWidth, frontTiresWidth, backTiresWidth *float64,
	frontTiresAspectRatio, backTiresAspectRatio *int, frontTiresRimDiameter, backTiresRimDiameter int) (float64, float64) {

	// диаметр передних шин в метрах
	var frontTiresDiameter float64
	// диаметр задних шин в метрах
	var backTiresDiameter float64

	if *frontTiresWidth == 0 || *backTiresWidth == 0 || *frontTrackWidth == 0 || *backTrackWidth == 0 || *frontTiresAspectRatio == 0 ||
		*backTiresAspectRatio == 0 || frontTiresRimDiameter == 0 || backTiresRimDiameter == 0 {
		*frontTiresAspectRatio = 1
		*backTiresAspectRatio = 1
	} else {
		// высота профиля передних шин
		frontTiresProfile := *frontTiresWidth * (float64(*frontTiresAspectRatio) / 100)
		// высота профиля задних шин
		backTiresProfile := *backTiresWidth * (float64(*backTiresAspectRatio) / 100)
		// float64(ts.FrontTiresRimDiameter)*0.0254 - это перевод в метры из дюймов inches
		frontTiresDiameter = float64(frontTiresRimDiameter)*0.0254 + 2*frontTiresProfile
		// float64(ts.BackTiresRimDiameter)*0.0254 - это перевод в метры из дюймов inches
		backTiresDiameter = float64(backTiresRimDiameter)*0.0254 + 2*backTiresProfile
	}
	return frontTiresDiameter, backTiresDiameter
}

// calculateComfortCoefficient вычисляет коэффициент комфорта
// Входные параметры: sps - информация о подвеске, grb - информация о коробке передач, cmc - информация о микроклимате салона,
// idn - информация об отделке салона, seo - информация об электропакете салона, mts - информация о мультимедийных системах,
// lts - информация о фонарях, powerSteeringType - тип рулевого усилителя, carAlarm - информация о наличии сигнализации,
// trunkVolume - объем багажника в литрах
func calculateComfortCoefficient(sps models.Suspension, gearbox string, cmc models.CabinMicroclimate,
	idn models.Interior, seo models.SetOfElectricOptions, mts models.MultimediaSystems, lts models.Lights,
	powerSteering models.PowerSteering, carAlarm models.Availability, trunkVolume float64) float64 {

	frontSuspensionCoefficient := calculateFrontSuspensionCoeffForComfortCoeff(sps.FrontSuspension, sps.FrontStabilizer)
	backSuspensionCoefficient := calculateBackSuspensionCoeffForComfortCoeff(sps.BackSuspension, sps.BackStabilizer)

	var powerSteeringTypeCoefficient float64 = 0
	switch powerSteering {
	case "Электроусилитель руля", "Гидроусилитель руля", "Электрогидроусилитель руля":
		powerSteeringTypeCoefficient = 2
	}

	var gearboxCoefficient float64 = 0
	switch gearbox {
	case "АКПП 6", "АКПП 5", "Вариатор":
		gearboxCoefficient = 4
	}

	var climateCoefficient float64 = 0
	switch {
	case cmc.AirConditioner == models.YesValue && cmc.ClimateControl == models.YesValue:
		climateCoefficient = 3
	case cmc.AirConditioner == models.YesValue && cmc.ClimateControl == models.NoValue:
		climateCoefficient = 2
	}

	var interiorCoefficient float64 = 0
	if idn.Upholstery == "Кожаная" {
		interiorCoefficient = 0.2962962962962963
	}

	var lightsCoefficient float64 = 0
	if lts.Headlights != "Галогенные" {
		lightsCoefficient += 0.8888888888888888
	}
	if lts.LEDRunningLights == models.YesValue {
		lightsCoefficient += 0.1
	}
	if lts.LEDTailLights == models.YesValue {
		lightsCoefficient += 0.1
	}
	if lts.FrontFogLights == models.YesValue {
		lightsCoefficient += 0.2962962962962963
	}
	if lts.BackFogLights == models.YesValue {
		lightsCoefficient += 0.2962962962962963
	}
	if lts.LightSensor == models.YesValue {
		lightsCoefficient += 0.2962962962962963
	}

	electricOptionsCoefficient := calculateElectricOptionsCoefficient(seo)

	var trunkVolumeCoefficient float64 = 0
	if trunkVolume > 500 {
		trunkVolumeCoefficient = 0.8888888888888888
	}

	var carAlarmCoefficient float64 = 0
	if carAlarm == models.YesValue {
		carAlarmCoefficient = 0.5925925925925926
	}

	var multimediaCoefficient float64 = 0
	if mts.OnBoardComputer == models.YesValue {
		multimediaCoefficient += 0.2962962962962963
	}
	if mts.MP3Support == models.YesValue {
		multimediaCoefficient += 0.2962962962962963
	}
	if mts.HandsFreeSupport == models.YesValue {
		multimediaCoefficient += 0.2962962962962963
	}

	sizeCoefficient := 0.8
	comfortCoefficient := (frontSuspensionCoefficient + backSuspensionCoefficient + powerSteeringTypeCoefficient +
		gearboxCoefficient + climateCoefficient + interiorCoefficient + lightsCoefficient + electricOptionsCoefficient +
		trunkVolumeCoefficient + carAlarmCoefficient + multimediaCoefficient) * sizeCoefficient
	return comfortCoefficient
}

// calculateFrontSuspensionCoeffForComfortCoeff вычисляет коэффициент типа передней подвески для коэффициента комфорта
// Входные параметры: frontSuspension - тип передней подвески, frontStabilizer - наличие переднего стабилизатора
func calculateFrontSuspensionCoeffForComfortCoeff(frontSuspension string, frontStabilizer models.Availability) float64 {
	var frontSuspensionCoefficient float64 = 0
	switch frontSuspension {
	case suspensionType2:
		if frontStabilizer == models.YesValue {
			frontSuspensionCoefficient = 3.8
		} else {
			frontSuspensionCoefficient = 3
		}
	case suspensionType1:
		if frontStabilizer == models.YesValue {
			frontSuspensionCoefficient = 3.6
		} else {
			frontSuspensionCoefficient = 2.8
		}
	case suspensionType3:
		if frontStabilizer == models.YesValue {
			frontSuspensionCoefficient = 4
		} else {
			frontSuspensionCoefficient = 3.2
		}
	case suspensionType4:
		if frontStabilizer == models.YesValue {
			frontSuspensionCoefficient = 3.4
		} else {
			frontSuspensionCoefficient = 2.6
		}
	case suspensionType5:
		if frontStabilizer == models.YesValue {
			frontSuspensionCoefficient = 2.8
		} else {
			frontSuspensionCoefficient = 2
		}
	case suspensionType6:
		if frontStabilizer == models.YesValue {
			frontSuspensionCoefficient = 2.4
		} else {
			frontSuspensionCoefficient = 1.6
		}
	case suspensionType7:
		if frontStabilizer == models.YesValue {
			frontSuspensionCoefficient = 1.8
		} else {
			frontSuspensionCoefficient = 1
		}
	}
	return frontSuspensionCoefficient
}

// calculateBackSuspensionCoeffForComfortCoeff вычисляет коэффициент типа задней подвески для коэффициента комфорта
// Входные параметры: backSuspension - тип задней подвески, backStabilizer - наличие заднего стабилизатора
func calculateBackSuspensionCoeffForComfortCoeff(backSuspension string, backStabilizer models.Availability) float64 {
	var backSuspensionCoefficient float64 = 0
	switch backSuspension {
	case suspensionType1:
		if backStabilizer == models.YesValue {
			backSuspensionCoefficient = 3.8
		} else {
			backSuspensionCoefficient = 3
		}
	case suspensionType2:
		if backStabilizer == models.YesValue {
			backSuspensionCoefficient = 3.6
		} else {
			backSuspensionCoefficient = 2.8
		}
	case suspensionType3:
		if backStabilizer == models.YesValue {
			backSuspensionCoefficient = 4
		} else {
			backSuspensionCoefficient = 3.2
		}
	case suspensionType4:
		if backStabilizer == models.YesValue {
			backSuspensionCoefficient = 3.4
		} else {
			backSuspensionCoefficient = 2.6
		}
	case suspensionType5:
		if backStabilizer == models.YesValue {
			backSuspensionCoefficient = 2.8
		} else {
			backSuspensionCoefficient = 2
		}
	case suspensionType6:
		if backStabilizer == models.YesValue {
			backSuspensionCoefficient = 2.4
		} else {
			backSuspensionCoefficient = 1.6
		}
	case suspensionType7:
		if backStabilizer == models.YesValue {
			backSuspensionCoefficient = 1.8
		} else {
			backSuspensionCoefficient = 1
		}
	}
	return backSuspensionCoefficient
}

// calculateElectricOptionsCoefficient вычисляет коэффициент наличия электрических опиций
// Входные параметры: seo - информация об электропакете салона
func calculateElectricOptionsCoefficient(seo models.SetOfElectricOptions) float64 {
	var electricOptionsCoefficient float64 = 0
	if seo.ElectricFrontSideWindowsLifts == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricBackSideWindowsLifts == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricHeatingOfFrontSeats == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricHeatingOfBackSeats == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricHeatingOfSteeringWheel == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricHeatingOfWindshield == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricHeatingOfRearWindow == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricHeatingOfSideMirrors == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricDriveOfDriverSeat == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricDriveOfFrontSeats == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricDriveOfSideMirrors == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.ElectricTrunkOpener == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	if seo.RainSensor == models.YesValue {
		electricOptionsCoefficient += 0.2962962962962963
	}
	return electricOptionsCoefficient
}

// calculateSafetyCoefficient вычисляет коэффициент безопасности
// Входные параметры: crashTestEstimate - результат краш-теста, asmc - информация о наличии электронных систем безопасности и
// контроля движения, sab - информация о наличии подушек безопасности, bkt - информация о типах тормозов
func calculateSafetyCoefficient(crashTestEstimate float64, smc models.SafetyAndMotionControlSystems, sab models.SetOfAirbags,
	bkt models.Brakes) float64 {

	var controlSystemCoefficient float64 = 0
	if smc.ABS == models.YesValue {
		controlSystemCoefficient += 3
	}
	if smc.ESP == models.YesValue {
		controlSystemCoefficient += 1
	}
	if smc.EBD == models.YesValue {
		controlSystemCoefficient += 1
	}
	if smc.BAS == models.YesValue {
		controlSystemCoefficient += 1
	}
	if smc.TCS == models.YesValue {
		controlSystemCoefficient += 1
	}

	var airbagsCoefficient float64 = 0
	if sab.DriverAirbag == models.YesValue {
		airbagsCoefficient += 1
	}
	if sab.FrontPassengerAirbag == models.YesValue {
		airbagsCoefficient += 1
	}
	if sab.SideAirbags == models.YesValue {
		airbagsCoefficient += 1
	}
	if sab.CurtainAirbags == models.YesValue {
		airbagsCoefficient += 1
	}

	var frontBrakesCoefficient float64 = 0
	switch bkt.FrontBrakes {
	case brakesType1, brakesType2:
		frontBrakesCoefficient = 2
	}

	var backBrakesCoefficient float64 = 0
	switch bkt.BackBrakes {
	case brakesType1, brakesType2:
		backBrakesCoefficient = 2
	}

	safetyCoefficient := crashTestEstimate + controlSystemCoefficient + airbagsCoefficient + frontBrakesCoefficient + backBrakesCoefficient
	return safetyCoefficient
}

// sigmoidFunction вычисляет значение функции `сигмоида`. Является функцией принадлежности
// Аппроксимирует точки, абсциссами которых являются значения коэффициентов, например, "управляемости" и т.д.,
// а ординатами - значения принадлежности точек нечетким подмножествам, например, "Высокая безопасность"
// нечеткого множества "Безопасность" и т.д.
func sigmoidFunction(xParam float64, lParam float64, kParam float64, x0Param float64) float64 {
	return lParam / (1 + math.Exp(-kParam*(xParam-x0Param)))
}

// gaussianFunction вычисляет значение функции Гаусса. Также является функцией принадлежности
// Аппроксимирует точки, абсциссами которых являются значения коэффициентов, например, "управляемости" и т.д.,
// а ординатами - значения принадлежности точек нечетким подмножествам "низкий", "средний", "высокий" нечетких
// множеств "экономичность", "динамика", "управляемость", "комфорт", "безопасность"
func gaussianFunction(xParam, amp, cen, wid float64) float64 {
	return amp / (math.Sqrt(2*math.Pi) * wid) * math.Exp(-math.Pow(xParam-cen, 2)/(2*math.Pow(wid, 2)))
}

// calculateMembershipFunctionValueForLinguisticVariables вычисляет значение одной из трех функций принадлежности,
// соответствующих нечетким подмножествам "низкий", "средний" и "высокий" нечетких множеств "экономичность", "динамика",
// "управляемость", "комфорт", "безопасность"
// Входные параметры: params - параметры функции `сигмоида` или функции Гаусса, term - название нечеткого подмножества,
// value - значение коэффициента экономичности или динамики или управляемости или комфорта или безопасности
func calculateMembershipFunctionValueForLinguisticVariables(params [3]float64, term string, value float64) float64 {
	var valueOfMembershipFunction float64
	switch term {
	case "низкий":
		valueOfMembershipFunction = sigmoidFunction(value, params[0], params[1], params[2])
	case "средний":
		valueOfMembershipFunction = gaussianFunction(value, params[0], params[1], params[2])
	case "высокий":
		valueOfMembershipFunction = sigmoidFunction(value, params[0], params[1], params[2])
	}

	if valueOfMembershipFunction < 0 {
		valueOfMembershipFunction = 0
	} else if valueOfMembershipFunction > 1 {
		valueOfMembershipFunction = 1
	}
	return valueOfMembershipFunction
}

// Параметры функции `сигмоида` или функции Гаусса
type Parameters struct {
	// параметры для функции принадлежности, соответствующей нечеткому подмножеству "низкий"
	paramsLow [3]float64
	// параметры для функции принадлежности, соответствующей нечеткому подмножеству "средний"
	paramsAverage [3]float64
	// параметры для функции принадлежности, соответствующей нечеткому подмножеству "высокий"
	paramsHigh [3]float64
}

// provideFunctionParameters предоставляет параметры для функций принадлежности, соответствующих
// нечетким подмножествам "низкий", "средний" и "высокий" нечетких множеств "экономичность",
// "динамика", "управляемость", "комфорт", "безопасность"
func provideFunctionParameters(variable, term string) [3]float64 {
	// Значения параметров были получены путем аппроксимации точек методом Левенберга-Марквардта с помощью пакета "Lmfit" ЯП "Python"
	data := make(map[string]Parameters)
	data["экономичность"] = Parameters{[3]float64{1.043723139993038, 0.5194913435480255, 11.165188013054621},
		[3]float64{2.2900397063026374, 9.43414665796981, 2.4138470099365112},
		[3]float64{1.949834151590793, -0.3804532441502327, 5.188639378787266}}

	data["динамика"] = Parameters{[3]float64{1.0231319819933777, 0.5016231903133455, 13.44547910618538},
		[3]float64{4.59765854168931, 10.810654375352698, 3.529577571232097},
		[3]float64{1.1836613715914706, -0.3792245799359442, 7.418367289135995}}

	data["управляемость"] = Parameters{[3]float64{1.2027418825694678, -0.07501884950616336, 22.4088071782321},
		[3]float64{31.124751770295614, 44.305695848946904, 19.042293165507264},
		[3]float64{1.3245201627428753, 0.07214432798281176, 69.8908258450921}}

	data["комфорт"] = Parameters{[3]float64{1.1834768835495801, -0.2870468773928149, 5.7724209993240825},
		[3]float64{5.0422852289029185, 10.10928688292464, 4.210219040980836},
		[3]float64{1.2492396969207602, 0.27009941927484593, 14.702080359730674}}

	data["безопасность"] = Parameters{[3]float64{1.3490219429573107, -0.21934076866027877, 4.73473258614666},
		[3]float64{5.450821590257078, 10.048764235757659, 4.185552288427339},
		[3]float64{1.2799644032509998, 0.3119397892443973, 15.65152657451626}}

	switch variable {
	case "экономичность":
		return pickParams(term, "экономичность", data)

	case "динамика":
		return pickParams(term, "динамика", data)

	case "управляемость":
		return pickParams(term, "управляемость", data)

	case "комфорт":
		return pickParams(term, "комфорт", data)

	case "безопасность":
		return pickParams(term, "безопасность", data)
	}
	return [3]float64{}
}

// pickParams возвращает найденные параметры для функций принадлежности, соответствующих нечетким
// подмножествам  "низкий", "средний" и "высокий" нечетких множеств "экономичность", "динамика",
// "управляемость", "комфорт", "безопасность"
// Входные параметры: term - название нечеткого подмножества, value - значение коэффициента,
// data - значения параметров для нечетких множеств (ключ: название нечеткого множества:
// например, "экономичность", значение - объект структуры Parameters)
func pickParams(term, value string, data map[string]Parameters) [3]float64 {
	switch term {
	case "низкий":
		return data[value].paramsLow
	case "средний":
		return data[value].paramsAverage
	case "высокий":
		return data[value].paramsHigh
	}
	return [3]float64{}
}

// findMin находит минимальное значение в срезе
func findMin(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}

	min := numbers[0]
	for _, num := range numbers {
		if num < min {
			min = num
		}
	}
	return min
}

// defuzzyficate реализует метод дефаззификации: метод центра тяжести. Этот метод возвращает конкретное число или
// четкое значение (в какой степени будет рекомендоваться автомобиль)
// Входные параметры: minValuesOfMemebershipFunction - ординаты вершин треугольников, образуемых под графиками функций принадлежности
// нечетких подмножеств нечеткого множества "рекомендация" или глобальные максимумы этих функций
// valuesRecommendation - абсциссы центральных точек оснований треугольников или центральные точки нечетких подмножеств
// нечеткого множества "рекомендация"
func defuzzyficate(minValuesOfMemebershipFunction []float64, valuesRecommendation []int) float64 {
	var centersOfMassOfTheAreaUnderTheGraph []float64
	var areasOfTriangles []float64
	for idx := 0; idx < len(minValuesOfMemebershipFunction); idx++ {
		centersOfMassOfTheAreaUnderTheGraph = append(centersOfMassOfTheAreaUnderTheGraph, calculateCenterOfMassOfTheAreaUnderTheGraph(
			minValuesOfMemebershipFunction[idx], valuesRecommendation[idx]))

		areasOfTriangles = append(areasOfTriangles, calculateAreaOfTriangle(minValuesOfMemebershipFunction[idx]))
	}

	// расчет средневзвешенного значения, которое является значением рекомендации и выходным значением
	// нечеткого алгоритма для одного автомобиля
	var enumerator float64
	for idx := range areasOfTriangles {
		enumerator += areasOfTriangles[idx] * centersOfMassOfTheAreaUnderTheGraph[idx]
	}

	var denominator float64
	for idx := 0; idx < len(areasOfTriangles); idx++ {
		denominator += areasOfTriangles[idx]
	}

	result := enumerator / denominator
	return result
}

// calculateCenterOfMassOfTheAreaUnderTheGraph вычисляет абсциссу центра тяжести площади треугольника, расположенного
// под графиком функции принадлежности μG(r) нечеткого подмножества G множества "рекомендация"
// Входные параметры: minValueOfMemebershipFunction - ордината вершины треугольника или
// глобальный максимум функции μG(r), valueReсommendation - абсцисса центральной точки основания
// треугольника или центральная точка нечеткого подмножества G нечеткого множества "рекомендация"
func calculateCenterOfMassOfTheAreaUnderTheGraph(minValueOfMemebershipFunction float64, valueReсommendation int) float64 {
	xValues := setX(float64(valueReсommendation), 1)
	area := applyTrapezoidalRule(func(x float64) float64 {
		return setY(x, xValues[0], xValues[1], xValues[2], minValueOfMemebershipFunction)
	}, xValues[0], xValues[2], 10000)

	numerator := applyTrapezoidalRule(func(x float64) float64 {
		return float64(valueReсommendation) * setY(x, xValues[0], xValues[1], xValues[2], minValueOfMemebershipFunction)
	}, xValues[0], xValues[2], 10000)
	return numerator / area
}

// setX вычисляет абсциссы левой и правой точек основания треугольника и возвращает абсциссы точек
// основания треугольника: левую, центральную и правую точки
// Входные параметры: center - абсцисса центральной точки, deviation - отклонение от центральной точки
func setX(center, deviation float64) []float64 {
	xArg := []float64{}
	xArg = append(xArg, center-deviation)
	xArg = append(xArg, center)
	xArg = append(xArg, center+deviation)
	return xArg
}

// applyTrapezoidalRule реализует численный метод трапеций, вычисляющий определенный интеграл
// Входные параметры: fun - интегрируемая функция, aParam - нижний предел интегрирования,
// bParam - верхний предел интегрирования, nParam - количество интервалов, на которые
// разбивается область под графиком интегрируемой функции
func applyTrapezoidalRule(fun func(float64) float64, aParam, bParam float64, nParam int) float64 {
	hParam := (bParam - aParam) / float64(nParam)
	sParam := 0.5 * (fun(aParam) + fun(bParam))
	for idx := 1; idx < nParam; idx++ {
		sParam += fun(aParam + float64(idx)*hParam)
	}
	return hParam * sParam
}

// setY реализует кусочную функцию, графиком которой является треугольник. Эта кусочная функция является
// функцией принадлежности μG(r) нечеткого подмножества G множества "рекомендация"
// Входные параметры: xValue - аргумент кусочной функции, center - абсцисса центральной точки
// основания треугольника, lefbound - абсцисса левой точки основания треугольника,
// rightBound - абсцисса правой точки основания треугольника, yMax - ордината вершины треугольника или
// глобальный максимум функции μG(r)
func setY(xValue, leftBound, center, rightBound, yMax float64) float64 {
	if xValue <= leftBound || xValue >= rightBound {
		return 0
	}
	if xValue <= center && xValue >= leftBound {
		return xValue - leftBound - (1 - yMax)
	}
	if leftBound <= xValue && xValue <= rightBound {
		return rightBound - xValue - (1 - yMax)
	}

	return 0
}

// calculateAreaOfTriangle вычисляет площадь фигуры(треугольника), образуемой под графиком
// функции принадлежности μG(r) нечеткого подмножества G множества "рекомендация"
// Входные параметры: minValueOfMemebershipFunction - высота треугольника или
// ордината вершины или глобальный максимум функции μG(r)
func calculateAreaOfTriangle(minValueOfMemebershipFunction float64) float64 {
	// baseOfTriangle - основание треугольника
	var baseOfTriangle float64 = 2
	area := 0.5 * baseOfTriangle * minValueOfMemebershipFunction
	return area
}
