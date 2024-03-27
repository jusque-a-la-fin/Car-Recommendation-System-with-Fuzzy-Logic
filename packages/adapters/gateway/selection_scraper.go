package gateway

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"vehicles/packages/domain/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

/***************************************************************************************************************************************************************/
/*   Здесь представлен алгоритм, собирающий данные об автомобилях различных марок с интернет-портала объявлений о продаже автомобилей. Эти данные необходимы   *
* для ранжирования автомобилей и далее для отображения на страницах данного веб-приложения. На примере сбора данных автомобилей одной марки можно понять       *
* работу алгоритма. Данные собираются следующим образом: функция prepareLink формирует url, который ведет на страницу марки. Эта страница содержит множество   *
* объявлений, каждое из которых включает название автомобиля с годом выпуска, цену, краткий перечень характеристик (тип двигателя: например, "дизель 3.0",     *
* "бензин 2.0", количество л.с., тип коробки передач, тип привода) и фотографию. HTML-код каждого объявления содержит ссылку на страницу конкретного           *
* автомобиля. Будем далее считать ссылки на страницы автомобилей ссылками на автомобили. Функция scrapeLinksNamesPrices собирает с страницы марки              *
* указанное количество ссылок на автомобили(в метод GetCarsUsingScraping передается срез makes, который имеет поле NumberOfCars, задающее ограничение          *
* по количеству ссылок на автомобили данной марки). Вместе с ссылками на автомобили функция scrapeLinksNamesPrices собирает также из объявлений названия       *
* автомобилей с годом выпуска и их цены. Далее алгоритм вызывает функцию scrapeCharacteristics, которая переходит по ссылкам на страницы автомобилей и         *
* собирает там информацию, такую как описание автомобиля, его фотографии, пробег в км, год выпуска, расположение руля, цвет, тип кузова, краткий перечень      *
* характеристик (тип двигателя: например, "дизель 3.0", "бензин 2.0"; количество л.с.; тип коробки передач; тип привода), а также ссылку на страницу           *
* комплектации и ссылку на страницу поколения, если они есть. Название же комплектации и название поколения автомобиля перечислены на странице автомобиля и их *
* HTML-коды содержат ссылки на страницу комплектации и на страницу поколения. Если есть название комплектации, то алгоритм вызывает функцию                    *
* scrapePageOfComplectationLink, которая переходит по ссылке, закрепленной за названием комплектации, и попадает на страницу комплектации, где перечислены     *
* подробно разичные характеристики данного автомобиля, и собирает некоторые из их. Если названия комплектации нет, а название поколения есть, то алгоритм      *
* вызывает функцию scrapePageOfGenerationLink, которая переходит по ссылке, закрепленной за названием поколения, и попадает на страницу, где представлены      *
* названия различных комплектаций, каждая из которых характеризуется кратким перечнем характеристик (тип двигателя: например, "дизель 3.0", "бензин 2.0";      *
* количество л.с.; тип коробки передач; тип привода). Если есть совпадение по всем характеристикам между кратким перечнем характеристик какой-либо             *
* комплектации и кратким перечнем характеристик на странице автомобиля, то алгоритм считает, что автомобиль имеет такую комплектацию и вызывает функцию        *
* scrapePageOfComplectationLink, которая переходит по ссылке, закрепленной за названием пододшедшей комплектации, на страницу данной комлпектации и оттуда     *
* собирает требуемые характеристики. Если подошедшая комплектация не была найдена, либо на странице автомобиля нет ни названия комплектации, ни названия       *
* поколения(либо названия есть, но ссылок нет), то в качестве характеристик автомобиля остаются некоторые из следующих: название, цена, пробег в км,           *
* расположение руля, цвет, тип кузова, краткий перечень характеристик (тип двигателя: например, "дизель 3.0", "бензин 2.0"; количество л.с.;                   *
* тип коробки передач; тип привода). Также остаются фотографии автомобиля, если они есть. Некоторые из перечисленных только что характеристик берутся в        *
* качестве дополнительных характеристик и в случае наличия ссылок на страницу поколения и страницу комплектации, где содержатся основные характеристики        *
* В то же время какие-то характеристики и из дополнительных и основных могут отсутствовать.                                                                    *                                                                 *
****************************************************************************************************************************************************************/

const (
	checkMark  = "#yes"
	option     = "#option"
	notOnBoard = "нет"
	newCarWord = "новый автомобиль"
)

// ScrapeSelectionCars собирает данные автомобилей из интернета
// Входные параметры: minPrice  - минимальная цена, maxPrice - максимальная цена, makes - срез марок
func (slr *selectionRepository) ScrapeSelectionCars(minPrice, maxPrice string, makes []models.Makes) ([]models.Car, error) {
	// ссылки на страницы марок
	linksToMakePage := make([]string, len(makes))
	for idx, make := range makes {
		linksToMakePage[idx] = prepareLink(minPrice, maxPrice, make.Make)
	}

	// 1-я размерность - марки, 2-я размерность - автомобили этих марок
	// rawLinks - ссылки, представляющие конкретные автомобили разных марок
	rawLinks := make([][]string, len(makes))
	// rawNames - названия
	rawNames := make([][]string, len(makes))
	// rawPrices - цены
	rawPrices := make([][]string, len(makes))

	var newLinks, newNames, newPrices []string
	var err error
	for index, thisMake := range makes {
		newLinks, newNames, newPrices, err = scrapeLinksNamesPrices(linksToMakePage[index], thisMake.NumberOfCars)
		if err != nil {
			return nil, err
		}
		if len(newLinks) != 0 {
			rawLinks[index], rawNames[index], rawPrices[index] = newLinks, newNames, newPrices
		}
	}

	links, names, prices, numberOfCars := getLinksNamesPrices(rawLinks, rawNames, rawPrices)

	cars := make([]models.Car, numberOfCars)
	uniqueID := 0
	for idx := 0; idx < len(links); idx++ {
		for jx := 0; jx < len(links[idx]); jx++ {
			car := models.NewCar()
			car.ID = uniqueID
			car.FullName = names[idx][jx]
			car.Offering.Price = fmt.Sprintf("%s₽", prices[idx][jx])
			err := scrapeCharacteristics(&car, links[idx][jx], car.FullName)
			if err != nil {
				return nil, fmt.Errorf("error from `scrapeCharacteristics` function, package `gateway`: %#v", err)
			}

			cars[uniqueID] = car
			uniqueID++
		}
	}
	return cars, nil
}

// prepareLink формирует и возвращает ссылку на страницу марки автомобиля, откуда будут собираться данные
// Входные параметры: minPrice  - минимальная цена, maxPrice - максимальная цена, make - название марки
func prepareLink(minPrice, maxPrice, make string) string {
	var link = fmt.Sprintf("https://auto.drom.ru/%s/all/?", make)
	if minPrice != "" {
		link = fmt.Sprintf("%sminprice=%s&", link, maxPrice)
	}

	if maxPrice != "" {
		link = fmt.Sprintf("%smaxprice=%s&", link, maxPrice)
	}

	link = fmt.Sprintf("%sph=1&unsold=1", link)
	return link
}

// scrapeLinksNamesPrices собирает ссылки на страницы, содержащие сведения о автомобилях, а также
// их названия и цены.
// Входные параметры: link - ссылка на страницу марки, quantity - количество автомобилей для поиска
func scrapeLinksNamesPrices(link string, quantity int) ([]string, []string, []string, error) {
	document, err := getWebPage(link)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error from `getWebPage` function, package `gateway`: %#v", err)
	}

	links := make([]string, 0, quantity)
	names := make([]string, 0, quantity)
	prices := make([]string, 0, quantity)

	div := document.Find("div[data-bulletin-list=true]")
	limit := 0
	div.Find("a").EachWithBreak(func(i int, a *goquery.Selection) bool {
		if limit == quantity {
			return false
		}

		href, exists := a.Attr("href")
		if exists {
			links = append(links, href)
		} else {
			return false
		}

		span := a.Find("span[data-ftid=bull_title]")
		content := span.Text()
		if content != "" {
			names = append(names, content)
		}

		spanPrice := a.Find("span[data-ftid=bull_price]")
		contentPrice := spanPrice.Text()
		if contentPrice != "" {
			prices = append(prices, contentPrice)
		}

		limit++
		return true
	})

	return links, names, prices, nil
}

// getWebPage получает какую-либо веб-страницу
// Входной параметр: link - ссылка на веб-страницу
func getWebPage(link string) (*goquery.Document, error) {
	response, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("error from `Get` function, package `http`: error while sending GET request: %#v", err)
	}

	// смена кодировки страницы с windows-1251 на utf-8
	utfBody, err := iconv.NewReader(response.Body, "windows-1251", "utf-8")
	if err != nil {
		return nil, fmt.Errorf("error from `NewReader` function, package `iconv`: error while converting charset from windows-1251 to utf-8: %#v", err)
	}

	// создание объекта структуры, представляющего HTML документ
	document, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return nil, fmt.Errorf("error from `NewDocumentFromReader` function, package `goquery`: %#v", err)
	}

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error from `Close` method, package `io`): %#v", err)
	}

	return document, nil
}

// getLinksNamesPrices удаляет пустые срезы из срезов срезов и получает точное количество автомобилей,
// которое будет показано пользователю
// Входные параметры: rawLinks - ссылки на страницы автомобилей, rawNames - названия и rawPrices - цены
func getLinksNamesPrices(rawLinks, rawNames, rawPrices [][]string) ([][]string, [][]string, [][]string, int) {
	count := 0
	numberOfCars := 0
	for _, slice := range rawLinks {
		if len(slice) > 0 {
			count++
			numberOfCars += len(slice)
		}
	}

	links := make([][]string, 0, count)
	names := make([][]string, 0, count)
	prices := make([][]string, 0, count)

	for idx := range rawLinks {
		if len(rawLinks[idx]) > 0 {
			links = append(links, append([]string{}, rawLinks[idx]...))
		}

		if len(rawNames[idx]) > 0 {
			names = append(names, append([]string{}, rawNames[idx]...))
		}

		if len(rawPrices[idx]) > 0 {
			prices = append(prices, append([]string{}, rawPrices[idx]...))
		}
	}
	return links, names, prices, numberOfCars
}

// scrapeCharacteristics собирает характеристики автомобиля с его страницы
// Входные параметры: car - автомобиль, link - ссылка на страницу автомобиля, carName - название автомобиля
// Цикломатическая сложность игнорируется в целях оптимизации
//
//gocyclo:ignore
func scrapeCharacteristics(car *models.Car, link, carName string) error {
	document, err := getWebPage(link)
	if err != nil {
		return fmt.Errorf("error from `getWebPage` function, package `gateway`: %#v", err)
	}

	span := document.Find("span.css-1kb7l9z.e162wx9x0").Eq(1)
	if span.Text() != "" {
		car.Description = span.Text()
	}

	div := document.Find("div[data-ftid='bull-page_bull-gallery_thumbnails']")
	div.Find("a").Each(func(j int, a *goquery.Selection) {
		href, ok := a.Attr("href")
		if !ok {
			return
		}

		if strings.HasSuffix(href, ".jpg") || strings.HasSuffix(href, ".jpeg") || strings.HasSuffix(href, ".png") {
			car.Offering.PhotoURLs = append(car.Offering.PhotoURLs, href)
		}
	})

	// additionalParams - краткий перечень характеристик автомобиля, собранный с страницы автомобиля. В случае отсутствия
	// на странице автомобиля ссылки на страницу комплектации, но при наличии ссылки на страницу поколения additionalParams
	// сравнивается с краткими перечнями характеристик, соответствующих различным комплектациям данного поколения и
	// размещенных на странице поколения. И в случае совпадения одного из этих перечней с additionalParams, считается,
	// что автомобиль имеет ту комплектацию, которая соответствуюет совпавшему перечню. Далее собирается информация о найденной
	// комплектации и присваивается рассматриваемому автомобилю
	additionalParams := make(map[string]string)

	// rexp - шаблон регулярного выражения, которому должна соответствовать строка, содержащая сведения о двигателе.
	rexp := regexp.MustCompile(`(\W+\D\d.\d\D\W)|(\W+)`)
	additionalParams["Двигатель"] = rexp.FindString(document.Find("span.css-1jygg09.e162wx9x0").Text())

	car.Offering.Year, err = findYearOfManufacture(carName)
	if err != nil {
		return fmt.Errorf("error from `findYearOfManufacture` function, package `gateway`: %#v", err)
	}

	// complectationLink - ссылка на страницу комплектации
	var complectationLink string

	document.Find("td.css-1la7f7n.ezjvm5n0").EachWithBreak(func(index int, element *goquery.Selection) bool {
		switch element.Prev().Text() {
		case "Мощность":
			// rexp - шаблон регулярного выражения, которому должна соответствовать строка, содержащая сведения о мощности.
			rexp = regexp.MustCompile(`(\d{3}|\d{2}|\d{4})\Dл\.с\.`)
			additionalParams["Мощность"] = rexp.FindString(element.Text())
		case "Коробка передач":
			additionalParams["Коробка передач"] = element.Text()
		case "Привод":
			additionalParams["Привод"] = element.Text()
		case "Тип кузова":
			car.Specs.Body = element.Text()
		case "Цвет":
			car.Features.Color = element.Text()
		case "Руль":
			switch element.Text() {
			case "левый":
				car.Specs.SteeringWheel.SteeringWheelPosition = models.LeftPos
			case "правый":
				car.Specs.SteeringWheel.SteeringWheelPosition = models.RightPos
			}

		case "Комплектация":
			complectationLink, _ = element.Children().Attr("href")
		}
		return true
	})

	car.Offering.Kilometerage = document.Find("span.css-1osyw3j.ei6iaw00").Text()

	// если автомобиль новый
	if car.Offering.Kilometerage == "" {
		newCar := document.Find("span.css-ytyb35.e162wx9x0").Text()
		if newCar == newCarWord {
			car.Offering.Kilometerage = newCarWord
		}
	}

	// generationLink - ссылка на страницу поколения автомобиля, которая содержит ссылки на
	// страницы комплектаций, одна из которых подходит данному автомобилю.
	var generationLink string
	document.Find("a").EachWithBreak(func(index int, element *goquery.Selection) bool {
		item, exists := element.Attr("data-ga-stats-name")
		if exists {
			if item == "generation_link" {
				generationLink, _ = element.Attr("href")
				car.Generation = element.Text()
			}
		}
		return true
	})

	if complectationLink != "" {
		err = scrapePageOfComplectationLink(car, complectationLink)
		if err != nil {
			return fmt.Errorf("error from `scrapePageOfComplectationLink` function, package `gateway`: %#v", err)
		}

	} else if complectationLink == "" && generationLink != "" {
		err := scrapePageOfGenerationLink(car, additionalParams, generationLink)
		if err != nil {
			return fmt.Errorf("error from `scrapePageOfGenerationLink` function, package `gateway`: %#v", err)
		}
	}
	return nil
}

// findYearOfManufacture находит год выпуска в названии автомобиля
// Входной параметр: carName - название автомобиля
func findYearOfManufacture(carName string) (int, error) {
	rexp := regexp.MustCompile(`\d{4}`)
	matches := rexp.FindAllStringIndex(carName, -1)
	if len(matches) > 0 {
		lastMatch := matches[len(matches)-1]
		year, err := strconv.Atoi(carName[lastMatch[0]:lastMatch[1]])
		if err != nil {
			return -1, fmt.Errorf("error from `Atoi` function, package `strconv`: %#v", err)
		}
		return year, nil
	} else {
		return -1, fmt.Errorf("there'is no year in carName string")
	}
}

// scrapePageOfComplectationLink собирает характеристики автомобиля с страницы, содержащей сведения о комплектации автомобиля
// Входные параметры: car - автомобиль, complectationLink - ссылка на страницу комплектации
// Цикломатическая сложность игнорируется в целях оптимизации
//
//gocyclo:ignore
func scrapePageOfComplectationLink(car *models.Car, complectationLink string) error {
	document, err := getWebPage(complectationLink)
	if err != nil {
		return fmt.Errorf("error from `getWebPage` function, package `gateway`: %#v", err)
	}

	var errCapacity, errAcceleration, errMaxSpeed, errClearance, errLength, errWidth, errHeight,
		errNumberOfSeats, errWheelbase, errFrontTrackWidth, errBackTrackWidth, errMassKg,
		errTrunkVolume, errDragCoeff, errEngineMaxPower, errCityConsum, errHighwayConsum,
		errMixedConsum, errFrontTiresWidth, errFrontTiresAspectRatio, errFrontTiresRimDiameter,
		errBackTiresWidth, errBackTiresAspectRatio, errBackTiresRimDiameter error

	document.Find("td").EachWithBreak(func(index int, element *goquery.Selection) bool {
		switch element.Text() {
		case "Название комплектации":
			car.TrimLevel = strings.TrimSpace(element.Next().Text())

		case "Тип привода":
			car.Specs.Drive = strings.TrimSpace(element.Next().Text())

		case "Тип кузова":
			if car.Specs.Body == "" {
				car.Specs.Body = strings.TrimSpace(element.Next().Text())
			}

		case "Тип трансмиссии":
			car.Specs.Gearbox = strings.TrimSpace(element.Next().Text())

		case "Объем двигателя, куб.см":
			car.Specs.Engine.Capacity, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errCapacity = err
				return false
			}

		case "Время разгона 0-100 км/ч, с":
			car.Specs.Acceleration0To100, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errAcceleration = err
				return false
			}

		case "Максимальная скорость, км/ч":
			car.Specs.MaxSpeed, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if errMaxSpeed != nil {
				errMaxSpeed = err
				return false
			}

		case "Клиренс (высота дорожного просвета), мм":
			car.Specs.GroundClearance, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errClearance = err
				return false
			}

		case "Габариты кузова (Д x Ш x В), мм":
			carSize := regexp.MustCompile(`(\d+)`)
			dimensions := carSize.FindAllStringSubmatch(strings.TrimSpace(element.Next().Text()), -1)
			if len(dimensions) != 0 {
				car.Specs.Length, err = strconv.ParseFloat(dimensions[0][0], 64)
				if err != nil {
					errLength = err
				}

				car.Specs.Width, err = strconv.ParseFloat(dimensions[1][0], 64)
				if err != nil {
					errWidth = err
				}

				car.Specs.Height, err = strconv.ParseFloat(dimensions[2][0], 64)
				if err != nil {
					errHeight = err
				}
			}

		case "Число мест":
			car.Specs.NumberOfSeats, err = strconv.Atoi(strings.TrimSpace(element.Next().Text()))
			if err != nil {
				errNumberOfSeats = err
				return false
			}

		case "Колесная база, мм":
			car.Specs.Wheelbase, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errWheelbase = err
				return false
			}

		case "Ширина передней колеи, мм":
			car.Specs.FrontTrackWidth, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errFrontTrackWidth = err
				return false
			}

		case "Ширина задней колеи, мм":
			car.Specs.BackTrackWidth, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errBackTrackWidth = err
				return false
			}

		case "Масса, кг":
			car.Specs.Mass, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errMassKg = err
				return false
			}

		case "Объем багажника, л":
			trunk := regexp.MustCompile(`(\d+)`)
			car.Specs.TrunkVolume, err = strconv.ParseFloat(trunk.FindStringSubmatch(element.Next().Text())[1], 64)
			if err != nil {
				errTrunkVolume = err
				return false
			}

		case "Коэффициент аэродинамического сопротивления, cW":
			car.Specs.DragCoefficient, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errDragCoeff = err
				return false
			}

		case "Используемое топливо":
			car.Specs.Engine.FuelUsed = strings.TrimSpace(element.Next().Text())

		case "Тип двигателя":
			car.Specs.Engine.EngineType = strings.TrimSpace(element.Next().Text())

		case "Максимальная мощность, л.с. (кВт) при об./мин.":
			car.Specs.Engine.MaxPower, err = strconv.ParseFloat(extractEnginePowerValue(strings.TrimSpace(element.Next().Text())), 64)
			if err != nil {
				errEngineMaxPower = err
				return false
			}

		case "Максимальный крутящий момент, Н*м (кг*м) при об./мин.":
			car.Specs.Engine.MaxTorque = strings.TrimSpace(element.Next().Text())

		case "Расход топлива в городском цикле, л/100 км":
			car.Specs.CityFuelConsumption, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errCityConsum = err
				return false
			}

		case "Расход топлива за городом, л/100 км":
			car.Specs.HighwayFuelConsumption, err = strconv.ParseFloat(strings.TrimSpace(element.Next().Text()), 64)
			if err != nil {
				errHighwayConsum = err
				return false
			}

		case "Расход топлива в смешанном цикле, л/100 км":
			car.Specs.MixedFuelConsumption, err = strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(element.Next().Text()), ",", "."), 64)
			if err != nil {
				errMixedConsum = err
				return false
			}

		case "Гидроусилитель руля":
			if symbol, _ := element.Next().Children().Children().Attr("href"); symbol == checkMark {
				car.Specs.SteeringWheel.PowerSteering = "Гидроусилитель"
			}

		case "Электроусилитель руля":
			if symbol, _ := element.Next().Children().Children().Attr("href"); symbol == checkMark {
				car.Specs.SteeringWheel.PowerSteering = "Электроусилитель"
			}

		case "Передний стабилизатор":
			SetTheCharacteristic(element, &car.Specs.Suspension.FrontStabilizer)

		case "Передняя подвеска":
			car.Specs.Suspension.FrontSuspension = strings.TrimSpace(element.Next().Text())

		case "Задний стабилизатор":
			SetTheCharacteristic(element, &car.Specs.Suspension.BackStabilizer)

		case "Задняя подвеска":
			car.Specs.Suspension.BackSuspension = strings.TrimSpace(element.Next().Text())

		case "Передние колеса":
			tires := regexp.MustCompile(`\d{3}/\d{2}\sR\d{2}`)
			tiresWidth := regexp.MustCompile(`^(\d+)`)
			tiresAspectRatio := regexp.MustCompile(`\/(\d+)`)
			tiresRimDiameter := regexp.MustCompile(`R(\d+)`)
			tiresStr := tires.FindString(element.Next().Text())
			if tiresStr != "" {
				tiresWidth := tiresWidth.FindString(tiresStr)
				if tiresWidth != "" {
					car.Specs.Tires.FrontTiresWidth, err = strconv.Atoi(tiresWidth)
					if err != nil {
						errFrontTiresWidth = err
						return false
					}
				}

				tiresAspectRatio := tiresAspectRatio.FindStringSubmatch(tiresStr)
				if len(tiresAspectRatio) > 1 {
					car.Specs.Tires.FrontTiresAspectRatio, err = strconv.Atoi(tiresAspectRatio[1])
					if err != nil {
						errFrontTiresAspectRatio = err
						return false
					}
				}

				tiresRimDiameter := tiresRimDiameter.FindStringSubmatch(tiresStr)
				if len(tiresRimDiameter) > 1 {
					car.Specs.Tires.FrontTiresRimDiameter, err = strconv.Atoi(tiresRimDiameter[1])
					if err != nil {
						errFrontTiresRimDiameter = err
						return false
					}
				}
			}

		case "Задние колеса":
			tires := regexp.MustCompile(`\d{3}/\d{2}\sR\d{2}`)
			tiresWidth := regexp.MustCompile(`^(\d+)`)
			tiresAspectRatio := regexp.MustCompile(`\/(\d+)`)
			tiresRimDiameter := regexp.MustCompile(`R(\d+)`)
			tiresStr := tires.FindString(element.Next().Text())
			if tiresStr != "" {
				tiresWidth := tiresWidth.FindString(tiresStr)
				if tiresWidth != "" {
					car.Specs.Tires.BackTiresWidth, err = strconv.Atoi(tiresWidth)
					if err != nil {
						errBackTiresWidth = err
						return false
					}
				}

				tiresAspectRatio := tiresAspectRatio.FindStringSubmatch(tiresStr)
				if len(tiresAspectRatio) > 1 {
					car.Specs.Tires.BackTiresAspectRatio, err = strconv.Atoi(tiresAspectRatio[1])
					if err != nil {
						errBackTiresAspectRatio = err
						return false
					}
				}
				tiresRimDiameter := tiresRimDiameter.FindStringSubmatch(tiresStr)
				if len(tiresRimDiameter) > 1 {
					car.Specs.Tires.BackTiresRimDiameter, err = strconv.Atoi(tiresRimDiameter[1])
					if err != nil {
						errBackTiresRimDiameter = err
						return false
					}
				}

			}
		case "Передние тормоза":
			car.Specs.Brakes.FrontBrakes = strings.TrimSpace(element.Next().Text())

		case "Задние тормоза":
			car.Specs.Brakes.BackBrakes = strings.TrimSpace(element.Next().Text())

		case "Стояночный тормоз":
			car.Specs.Brakes.ParkingBrake = strings.TrimSpace(element.Next().Text())

		case "Галогенные фары":
			if symbol, _ := element.Next().Children().Children().Attr("href"); symbol == checkMark {
				car.Features.Lights.Headlights = "Галогенные фары"
			}

		case "Биксеноновые фары":
			if symbol, _ := element.Next().Children().Children().Attr("href"); symbol == checkMark {
				car.Features.Lights.Headlights = "Биксеноновые фары"
			}

		case "Светодиодные фары":
			if symbol, _ := element.Next().Children().Children().Attr("href"); symbol == checkMark {
				car.Features.Lights.Headlights = "Светодиодные фары"
			}

		case "Лазерные фары":
			if symbol, _ := element.Next().Children().Children().Attr("href"); symbol == checkMark {
				car.Features.Lights.Headlights = "Лазерные фары"
			}

		case "Светодиодные ходовые огни":
			SetTheCharacteristic(element, &car.Features.Lights.LEDRunningLights)

		case "Передние противотуманные фары":
			SetTheCharacteristic(element, &car.Features.Lights.FrontFogLights)

		case "Светодиодные противотуманные фары":
			SetTheCharacteristic(element, &car.Features.Lights.FrontFogLights)

		case "Cветодиодные задние фонари":
			SetTheCharacteristic(element, &car.Features.Lights.LEDTailLights)

		case "Задние противотуманные фонари":
			SetTheCharacteristic(element, &car.Features.Lights.BackFogLights)

		case "Датчик света":
			SetTheCharacteristic(element, &car.Features.Lights.LightSensor)

		case "Датчик дождя":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.RainSensor)

		case "Электропривод боковых зеркал":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricDriveOfSideMirrors)

		case "Электроподогрев зеркал":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricHeatingOfSideMirrors)

		case "Электропривод багажника":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricTrunkOpener)

		case "Обогрев заднего стекла":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricHeatingOfRearWindow)

		case "Тканевая обивка салона":
			if s, _ := element.Next().Children().Children().Attr("href"); s == checkMark {
				car.Features.Interior.Upholstery = "Тканевая"
			}

		case "Кожаная обивка салона":
			if s, _ := element.Next().Children().Children().Attr("href"); s == checkMark {
				car.Features.Interior.Upholstery = "Кожаная"
			}

		case "Комбинированная обивка салона":
			if s, _ := element.Next().Children().Children().Attr("href"); s == checkMark {
				car.Features.Interior.Upholstery = "Комбинированная"
			}

		case "Электропривод водительского сиденья":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricDriveOfDriverSeat)

		case "Электропривод передних сидений":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricDriveOfFrontSeats)

		case "Электроподогрев передних сидений":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricHeatingOfFrontSeats)

		case "Электроподогрев задних сидений":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricHeatingOfBackSeats)

		case "Электрические стеклоподъемники передние":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricFrontSideWindowsLifts)

		case "Электрические стеклоподъемники задние":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricBackSideWindowsLifts)

		case "Электроподогрев рулевого колеса":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricHeatingOfSteeringWheel)

		case "Электроподогрев лобового стекла":
			SetTheCharacteristic(element, &car.Features.ElectricOptions.ElectricHeatingOfWindshield)

		case "Подушка безопасности водительская":
			SetTheCharacteristic(element, &car.Features.Airbags.DriverAirbag)

		case "Подушка безопасности переднего пассажира":
			SetTheCharacteristic(element, &car.Features.Airbags.FrontPassengerAirbag)

		case "Подушка безопасности боковая":
			SetTheCharacteristic(element, &car.Features.Airbags.SideAirbags)

		case "Подушки безопасности-шторки":
			SetTheCharacteristic(element, &car.Features.Airbags.CurtainAirbags)

		case "Антиблокировочная система (ABS)":
			SetTheCharacteristic(element, &car.Features.SafetyAndMotionControlSystem.ABS)

		case "Система распределения тормозного усилия (EBD)":
			SetTheCharacteristic(element, &car.Features.SafetyAndMotionControlSystem.EBD)

		case "Вспомогательная система торможения (BAS)":
			SetTheCharacteristic(element, &car.Features.SafetyAndMotionControlSystem.BAS)

		case "Система электронного контроля устойчивости (ESP)":
			SetTheCharacteristic(element, &car.Features.SafetyAndMotionControlSystem.ESP)

		case "Антипробуксовочная система (TCS)":
			SetTheCharacteristic(element, &car.Features.SafetyAndMotionControlSystem.TCS)

		case "Круиз-контроль":
			SetTheCharacteristic(element, &car.Features.SafetyAndMotionControlSystem.CruiseControl)

		case "Передний парктроник":
			SetTheCharacteristic(element, &car.Features.SafetyAndMotionControlSystem.FrontParkingSensor)

		case "Задний парктроник":
			SetTheCharacteristic(element, &car.Features.SafetyAndMotionControlSystem.BackParkingSensor)

		case "Кондиционер":
			SetTheCharacteristic(element, &car.Features.CabinMicroclimate.AirConditioner)

		case "Климат-контроль":
			SetTheCharacteristic(element, &car.Features.CabinMicroclimate.ClimateControl)

		case "Поддержка MP3":
			SetTheCharacteristic(element, &car.Features.MultimediaSystems.MP3Support)

		case "Hands free":
			SetTheCharacteristic(element, &car.Features.MultimediaSystems.HandsFreeSupport)

		case "Камера заднего обзора":
			SetTheCharacteristic(element, &car.Features.SafetyAndMotionControlSystem.RearViewCamera)

		case "Бортовой компьютер":
			SetTheCharacteristic(element, &car.Features.MultimediaSystems.OnBoardComputer)

		case "Сигнализация":
			SetTheCharacteristic(element, &car.Features.CarAlarm)
		}
		return true
	})

	switch {
	case errCapacity != nil:
		return fmt.Errorf("error `errCapacity` from `ParseFloat` function, package `strconv`: %#v", errCapacity)
	case errAcceleration != nil:
		return fmt.Errorf("error `errAcceleration` from `ParseFloat` function, package `strconv`: %#v", errAcceleration)
	case errMaxSpeed != nil:
		return fmt.Errorf("error `errMaxSpeed` from `ParseFloat` function, package `strconv`: %#v", errMaxSpeed)
	case errClearance != nil:
		return fmt.Errorf("error `errClearance` from `ParseFloat` function, package `strconv`: %#v", errClearance)
	case errLength != nil:
		return fmt.Errorf("error `errLength` from `ParseFloat` function, package `strconv`: %#v", errLength)
	case errWidth != nil:
		return fmt.Errorf("error `errWidth` from `ParseFloat` function, package `strconv`: %#v", errWidth)
	case errHeight != nil:
		return fmt.Errorf("error `errHeight` from `ParseFloat` function, package `strconv`: %#v", errHeight)
	case errNumberOfSeats != nil:
		return fmt.Errorf("error `errNumberOfSeats` from `Atoi` function, package `strconv`: %#v", errNumberOfSeats)
	case errWheelbase != nil:
		return fmt.Errorf("error `errWheelbase` from `ParseFloat` function, package `strconv`: %#v", errWheelbase)
	case errFrontTrackWidth != nil:
		return fmt.Errorf("error `errFrontTrackWidth` from `ParseFloat` function, package `strconv`: %#v", errFrontTrackWidth)
	case errBackTrackWidth != nil:
		return fmt.Errorf("error `errBackTrackWidth` from `ParseFloat` function, package `strconv`: %#v", errBackTrackWidth)
	case errMassKg != nil:
		return fmt.Errorf("error `errMassKg` from `ParseFloat` function, package `strconv`: %#v", errMassKg)
	case errTrunkVolume != nil:
		return fmt.Errorf("error `errTrunkVolume` from `ParseFloat` function, package `strconv`: %#v", errTrunkVolume)
	case errDragCoeff != nil:
		return fmt.Errorf("error `errDragCoeff` from `ParseFloat` function, package `strconv`: %#v", errDragCoeff)
	case errEngineMaxPower != nil:
		return fmt.Errorf("error `errEngineMaxPower` from `ParseFloat` function, package `strconv`: %#v", errEngineMaxPower)
	case errCityConsum != nil:
		return fmt.Errorf("error `errCityConsum` from `ParseFloat` function, package `strconv`: %#v", errCityConsum)
	case errHighwayConsum != nil:
		return fmt.Errorf("error `errHighwayConsum` from `ParseFloat` function, package `strconv`: %#v", errHighwayConsum)
	case errMixedConsum != nil:
		return fmt.Errorf("error `errMixedConsum` from `ParseFloat` function, package `strconv`: %#v", errMixedConsum)
	case errFrontTiresWidth != nil:
		return fmt.Errorf("error `errFrontTiresWidth` from `Atoi` function, package `strconv`: %#v", errFrontTiresWidth)
	case errFrontTiresAspectRatio != nil:
		return fmt.Errorf("error `errFrontTiresAspectRatio` from `Atoi` function, package `strconv`: %#v", errFrontTiresAspectRatio)
	case errFrontTiresRimDiameter != nil:
		return fmt.Errorf("error `errFrontTiresRimDiameter` from `Atoi` function, package `strconv`: %#v", errFrontTiresRimDiameter)
	case errBackTiresWidth != nil:
		return fmt.Errorf("error `errBackTiresWidth` from `Atoi` function, package `strconv`: %#v", errBackTiresWidth)
	case errBackTiresAspectRatio != nil:
		return fmt.Errorf("error `errBackTiresAspectRatio` from `Atoi` function, package `strconv`: %#v", errBackTiresAspectRatio)
	case errBackTiresRimDiameter != nil:
		return fmt.Errorf("error `errBackTiresRimDiameter` from `Atoi` function, package `strconv`: %#v", errBackTiresRimDiameter)
	}

	return nil
}

// SetTheCharacteristic устанавливает значения для некоторых характеристик
// Входные параметры: element - представляет HTML документ, characteristic - характеристика автомобиля
func SetTheCharacteristic(element *goquery.Selection, characteristic *models.Availability) {
	element.Parent().Find("use").Each(func(i int, s *goquery.Selection) {
		if val, exists := s.Attr("href"); exists {
			switch val {
			case checkMark:
				*characteristic = models.YesValue
			case option:
				*characteristic = models.OptionValue
			default:
			}
		}
	})

	if element.Next().Find("span").Text() == "—" {
		*characteristic = models.NoValue
	}
}

// extractEnginePowerValue извлекает из строки число, обозначающее количество лошадиных сил
// Входной параметр: maxpower - максимальная мощность двигателя
func extractEnginePowerValue(maxpower string) string {
	digits := make([]rune, 0)
	for _, char := range maxpower {
		if unicode.IsDigit(char) {
			digits = append(digits, char)
		} else {
			break
		}
	}
	return string(digits)
}

// scrapePageOfGenerationLink находит на странице поколения автомобиля, которая содержит ссылки на страницы комплектаций,
// комплектацию, которая подходит текущему автомобилю. С подходящей страницы комплектации собирается информация о характеристиках.
// Входные параметры: car - автомобиль, generationLink - ссылка на страницу поколения этого автомобиля,
// additionalParams - краткий перечень характеристик автомобиля, собранный с страницы автомобиля
func scrapePageOfGenerationLink(car *models.Car, additionalParams map[string]string, generationLink string) error {
	document, err := getWebPage(generationLink)
	if err != nil {
		return fmt.Errorf("error from `getWebPage` function, package `gateway`: %#v", err)
	}

	// specificCharacteristics - краткий перечень характеристик автомобиля, собранный с страницы автомобиля.
	// Будет сопоставляться с краткими перечнями характеристик различных комплектаций,
	// представленных на странице поколения
	specificCharacteristics := make(map[string]string)

	// complectationLink - ссылка на страницу комплектации
	var complectationLink string

	specificCharacteristics["Двигатель"] = additionalParams["Двигатель"]
	specificCharacteristics["Мощность"] = additionalParams["Мощность"]
	specificCharacteristics["Коробка передач"] = additionalParams["Коробка передач"]
	specificCharacteristics["Привод"] = additionalParams["Привод"]

	// удаление из строки вида "249 л.с." подстроки "л.с", поскольку regexp.MatchString() почему-то не находит
	// совпадение строки вида "249 л.с" c строкой вида ""3.5 л, 249 л.с., бензин, АКПП, полный привод (4WD)
	specificCharacteristics["Мощность"] = strings.Replace(specificCharacteristics["Мощность"], "л.с.", "", 1)

	// удаление пробела из строки вида "249 ", потому что функция replace выше не может убрать подстроку " л.с" из "249 л.с."
	specificCharacteristics["Мощность"] = strings.TrimSpace(specificCharacteristics["Мощность"])

	if specificCharacteristics["Двигатель"] == "электро" {
		// замена "электро" на "электричество", поскольку только последний вариант присутствует на странице комплектаций
		specificCharacteristics["Двигатель"] = "электричество"
		// удаление коробки передач, поскольку её нет для электромобилей
		delete(specificCharacteristics, "Коробка передач")
	} else {

		// удаление пробелов из таких, например, строк "бензин, 2.3 л" и "бензин, 2.0 л, гибрид"
		specificCharacteristics["Двигатель"] = strings.ReplaceAll(specificCharacteristics["Двигатель"], " ", "")

		// разделение строки вида "бензин, 2.3 л" на две строки "бензин" и " 2.3 л"
		// может встретиться "бензин, 2.0 л, гибрид", которая будет разбита на "бензин", "2.0 л", "гибрид"
		// если встретиться "электро", то эта строка не будет разделена
		fuel := strings.Split(specificCharacteristics["Двигатель"], ",")

		// удаление "л" из "2.3 л", поскольку команда выше удаляет пробел, и получается "2.3л",
		// которая не равна строке "2.3 л", находящейся на странице комплектаций
		fuel[1] = strings.Replace(fuel[1], "л", "", 1)

		// удаление строки вида "бензин, 2.3 л"
		delete(specificCharacteristics, "Двигатель")

		// добавление этих двух строк "бензин", " 2.3 л" или трех строк "бензин", "2.0 л", "гибрид"
		for _, elem := range fuel {
			specificCharacteristics[elem] = elem
		}

		// если вместо "МКПП" И "АКПП", в specificCharacteristics["Коробка передач"] присутствуют "механика" и "автомат" соответственно,
		// то нужно заменить их на "МКПП" и "АКПП",
		// потому что на странице комплектаций всегда встречается "МКПП" и "АКПП", и нужно чтобы specificCharacteristics["Коробка передач"]
		// совпал с "МКПП" или "АКПП" на странице комплектаций
		if specificCharacteristics["Коробка передач"] == "механика" {
			specificCharacteristics["Коробка передач"] = "МКПП"
		} else if specificCharacteristics["Коробка передач"] == "автомат" {
			specificCharacteristics["Коробка передач"] = "АКПП"
		}
	}

	document.Find("th").EachWithBreak(func(index int, element *goquery.Selection) bool {
		item, exists := element.Attr("colspan")
		if exists {
			// 7 и 6 проверяются, потому что так был составлен HTML страницы поколения
			if item == "7" || item == "6" {
				// счетчик совпадений краткого перечня характеристик, присущего конкретной комплектации с
				// кратким переченем характеристик текущего автомобиля
				match := 0
				for _, value := range specificCharacteristics {
					matchStr := strings.Contains(element.Text(), value)
					if matchStr {
						match++
					}
				}
				// если совпали все характеристики
				if match == len(specificCharacteristics) {
					// rexp - шаблон регулярного выражения для ссылки на комплектацию
					rexp := regexp.MustCompile(`/catalog/.+/.+/\d+/`)
					attempt, _ := element.Parent().Next().Children().Find("a").Attr("href")
					if rexp.MatchString(attempt) {
						complectationLink = attempt
					}
					return false
				}
			}
		}
		return true
	})

	complectationLink = fmt.Sprintf("https://www.drom.ru%s", complectationLink)
	err = scrapePageOfComplectationLink(car, complectationLink)
	if err != nil {
		return fmt.Errorf("error from `scrapePageOfComplectationLink` function, package `gateway`: %#v", err)
	}
	return nil
}
