package models

const UndefinedStr string = "Неизвестно"
const Absent string = "Отсутствует"

// Car - автомобиль
type Car struct {
	ID int
	// FullName - название
	FullName string
	// Description - описание
	Description string
	// Generation - название поколения
	Generation string `db:"generation"`
	// Trimlevel - название комплектации
	TrimLevel string `db:"trim_level"`
	// Specifications - технические характеристики
	Specs Specifications
	// Features - опции
	Features Features
	// Offering - сведения для покупателя
	Offering Offering
}

// Specifications - технические характеристики
type Specifications struct {
	// Body - тип кузова
	Body string `db:"body"`
	// Length - длина, мм
	Length float64 `db:"length"`
	// Width - ширина, мм
	Width float64 `db:"width"`
	// Height - высота, мм
	Height float64 `db:"height"`
	// GroundClearance - клиренс, мм
	GroundClearance float64 `db:"ground_clearance"`
	// DragCoefficient - коэффициент аэродинамического сопротивления, cW
	DragCoefficient float64 `db:"drag_coefficient"`
	// FrontTrackWidth - ширина передней колеи, мм
	FrontTrackWidth float64 `db:"front_track_width"`
	// BackTrackWidth - ширина задней колеи, мм"
	BackTrackWidth float64 `db:"back_track_width"`
	// Wheelbase - колесная база, мм
	Wheelbase float64 `db:"wheelbase"`
	// Acceleration0To100 - время разгона 0-100 км/ч, с
	Acceleration0To100 float64 `db:"acceleration_0_to_100"`
	// MaxSpeed - максимальная скорость, км/ч
	MaxSpeed float64 `db:"max_speed"`
	// CityFuelConsumption - расход топлива в городском цикле, л/100 км
	CityFuelConsumption float64 `db:"city_fuel_consumption"`
	// HighwayFuelConsumption - расход топлива за городом, л/100 км
	HighwayFuelConsumption float64 `db:"highway_fuel_consumption"`
	// MixedFuelConsumption - расход топлива в смешанном цикле, л/100 км
	MixedFuelConsumption float64 `db:"mixed_fuel_consumption"`
	// NumberOfSeats - число мест
	NumberOfSeats int `db:"number_of_seats"`
	// TrunkVolume - объем багажника, литры
	TrunkVolume float64 `db:"trunk_volume"`
	// Mass - масса, кг
	Mass float64 `db:"mass"`
	// Gearbox - тип трансмиссии
	Gearbox string `db:"gearbox"`
	// Drive - тип привода
	Drive string `db:"drive"`
	// CrashTestEstimate - баллы за краш-тест
	CrashTestEstimate float64 `db:"crash_test_estimate"`
	// Engine - двигатель
	Engine Engine
	// SteeringWheel - рулевое колесо
	SteeringWheel SteeringWheel
	// Suspension - подвеска
	Suspension Suspension
	// Brakes - тормоза
	Brakes Brakes
	// Tires - шины
	Tires Tires
}

// Engine - двигатель
type Engine struct {
	// FuelUsed - используемое топливо
	FuelUsed string `db:"fuel_used"`
	// EngineType - тип двигателя
	EngineType string `db:"engine_type"`
	// Capacity - объем двигателя, куб.см
	Capacity float64 `db:"capacity"`
	// MaxPower - максимальная мощность, л.с.
	MaxPower float64 `db:"max_power"`
	// MaxTorque - максимальный крутящий момент, Н*м (кг*м) при об./мин.
	MaxTorque string `db:"max_torque"`
}

type SteeringWheelPosition string

const (
	LeftPos      SteeringWheelPosition = "Левый руль"
	RightPos     SteeringWheelPosition = "Правый руль"
	UndefinedPos SteeringWheelPosition = "Неизвестно"
)

type PowerSteering string

const (
	ElectricPS         PowerSteering = "Электроусилитель руля"
	ElectrohydraulicPS PowerSteering = "Электрогидроусилитель руля"
	HydraulicPS        PowerSteering = "Гидроусилитель руля"
	NoPS               PowerSteering = "Нет"
	UndefinedPS        PowerSteering = "Неизвестно"
)

// SteeringWheel - рулевое колесо
type SteeringWheel struct {
	// SteeringWheelPosition - положение руля(слева, справа и по центру)
	SteeringWheelPosition SteeringWheelPosition `db:"position"`
	// PowerSteering - тип усилителя рулевого управления
	PowerSteering PowerSteering `db:"power_steering"`
}

type Availability string

const (
	YesValue       Availability = "Есть"
	NoValue        Availability = "Нет"
	OptionValue    Availability = "Опция производителя"
	UndefinedValue Availability = "Неизвестно"
)

// Suspension - подвеска
type Suspension struct {
	// FrontStabilizer - наличие переднего стабилизатора
	FrontStabilizer Availability `db:"front_stabilizer"`
	// BackStabilizer - наличие заднего стабилизатора
	BackStabilizer Availability `db:"back_stabilizer"`
	// FrontSuspension - название типа передней подвески
	FrontSuspension string `db:"front_suspension"`
	// BackSuspension - название типа задней подвески
	BackSuspension string `db:"back_suspension"`
}

// Brakes - тормоза
type Brakes struct {
	// FrontBrakes - тип передних тормозов
	FrontBrakes string `db:"front_brakes"`
	// BackBrakes - тип задних тормозов
	BackBrakes string `db:"back_brakes"`
	// ParkingBrake - тип стояночного тормоза
	ParkingBrake string `db:"parking_brake"`
}

// Tires - шины
type Tires struct {
	// FrontTiresWidth - ширина передних шин, мм
	FrontTiresWidth int `db:"front_tires_width"`
	// BackTiresWidth - ширина задних шин, мм
	BackTiresWidth int `db:"back_tires_width"`
	// FrontTiresAspectRatio - процентное соотношение высоты профиля передних шин к их ширине
	FrontTiresAspectRatio int `db:"front_tires_aspect_ratio"`
	// BackTiresAspectRatio - процентное соотношение высоты профиля задних шин к их ширине
	BackTiresAspectRatio int `db:"back_tires_aspect_ratio"`
	// FrontTiresRimDiameter - диаметр обода передних шин, мм
	FrontTiresRimDiameter int `db:"front_tires_rim_diameter"`
	// BackTiresRimDiameter - диаметр обода задних шин, мм
	BackTiresRimDiameter int `db:"back_tires_rim_diameter"`
}

// Features - опции
type Features struct {
	// SafetyAndMotionControlSystems - электронные системы безопасности и контроля движения
	SafetyAndMotionControlSystem SafetyAndMotionControlSystems
	// Lights - фонари и фары
	Lights Lights
	// Interior- интерьер
	Interior Interior
	// CabinMicroclimate - микроклимат салона
	CabinMicroclimate CabinMicroclimate
	// ElectricOptions - пакет электрических опций
	ElectricOptions SetOfElectricOptions
	// Airbags - подушки безопасности
	Airbags SetOfAirbags
	// MultimediaSystem - системы мультимедиа
	MultimediaSystems MultimediaSystems
	// CarAlarm - наличие сигнализации
	CarAlarm Availability `db:"car_alarm"`
	// Color - цвет
	Color string `db:"color"`
}

// SafetyAndMotionControlSystems - электронные системы безопасности и контроля движения
type SafetyAndMotionControlSystems struct {
	// ABS - наличие антиблокировочной системы (ABS)
	ABS Availability `db:"abs_system"`
	// ESP - наличие системы электронного контроля устойчивости (ESP)
	ESP Availability `db:"esp_system"`
	// EBD - наличие системы распределения тормозного усилия (EBD)
	EBD Availability `db:"ebd_system"`
	// BAS - наличие вспомогательной системы торможения (BAS)
	BAS Availability `db:"bas_system"`
	// TCS - наличие антипробуксовочной системы (TCS)
	TCS Availability `db:"tcs_system"`
	// FrontParkingSensor - наличие переднего парктроника
	FrontParkingSensor Availability `db:"front_parking_sensor"`
	// BackParkingSensor - наличие заднего парктроника
	BackParkingSensor Availability `db:"back_parking_sensor"`
	// RearViewCamera - наличие камеры заднего обзора
	RearViewCamera Availability `db:"rear_view_camera"`
	// CruiseControl - наличие круиз-контроля
	CruiseControl Availability `db:"cruise_control"`
}

// Lights - фонари и фары
type Lights struct {
	// Headlights - тип передних фар
	Headlights string `db:"headlights"`
	// LEDRunningLights - наличие светодиодных ходовых огней
	LEDRunningLights Availability `db:"led_running_lights"`
	// LEDTailLights - наличие светодиодных задних фонарей
	LEDTailLights Availability `db:"led_tail_lights"`
	// LightSensor - наличие датчика света
	LightSensor Availability `db:"light_sensor"`
	// FrontFogLights - наличие передних противотуманных фар
	FrontFogLights Availability `db:"front_fog_lights"`
	// BackFogLights - наличие задних противотуманных фонарей
	BackFogLights Availability `db:"back_fog_lights"`
}

// Interior - интерьер
type Interior struct {
	// Upholstery - тип обивки салона
	Upholstery string `db:"upholstery"`
}

// CabinMicroclimate - микроклимат салона
type CabinMicroclimate struct {
	// AirConditioner - наличие кондиционера
	AirConditioner Availability `db:"air_conditioner"`
	// ClimateControl - наличие климат-контроля
	ClimateControl Availability `db:"climate_control"`
}

// SetOfElectricOptions - пакет электрических опций
type SetOfElectricOptions struct {
	// ElectricFrontSideWindowsLifts - наличие электрических передних стеклоподъемников
	ElectricFrontSideWindowsLifts Availability `db:"electric_front_side_windows_lifts"`
	// ElectricBackSideWindowsLifts - наличие электрических задних стеклоподъемников
	ElectricBackSideWindowsLifts Availability `db:"electric_back_side_windows_lifts"`
	// ElectricHeatingOfFrontSeats - наличие электроподогрева передних сидений
	ElectricHeatingOfFrontSeats Availability `db:"electric_heating_of_front_seats"`
	// ElectricHeatingOfBackSeats - наличие электроподогрева задних сидений
	ElectricHeatingOfBackSeats Availability `db:"electric_heating_of_back_seats"`
	// ElectricHeatingOfSteeringWheel - наличие электроподогрева рулевого колеса
	ElectricHeatingOfSteeringWheel Availability `db:"electric_heating_of_steering_wheel"`
	// ElectricHeatingOfWindshield - наличие электроподогрева лобового стекла
	ElectricHeatingOfWindshield Availability `db:"electric_heating_of_windshield"`
	// ElectricHeatingOfRearWindow -- наличие обогрева заднего стекла
	ElectricHeatingOfRearWindow Availability `db:"electric_heating_of_rear_window"`
	// ElectricHeatingOfSideMirrors - наличие электроподогрева боковых зеркал
	ElectricHeatingOfSideMirrors Availability `db:"electric_heating_of_side_mirrors"`
	// ElectricDriveOfDriverSeat - наличие электропривода водительского сидения
	ElectricDriveOfDriverSeat Availability `db:"electric_drive_of_driver_seat"`
	// ElectricDriveOfFrontSeats - наличие электропривода передних сидений
	ElectricDriveOfFrontSeats Availability `db:"electric_drive_of_front_seats"`
	// ElectricDriveOfSideMirrors - наличие электропривода боковых зеркал
	ElectricDriveOfSideMirrors Availability `db:"electric_drive_of_side_mirrors"`
	// ElectricTrunkOpener - наличие электропривода багажника
	ElectricTrunkOpener Availability `db:"electric_trunk_opener"`
	// RainSensor - наличие датчика дождя
	RainSensor Availability `db:"rain_sensor"`
}

// SetOfAirbags - подушки безопасности
type SetOfAirbags struct {
	// DriverAirbag - наличие водительской подушки безопасности
	DriverAirbag Availability `db:"driver_airbag"`
	// FrontPassengerAirbag - наличие подушки безопасности переднего пассажира
	FrontPassengerAirbag Availability `db:"front_passenger_airbag"`
	// SideAirbags - наличие боковых подушек безопасности
	SideAirbags Availability `db:"side_airbags"`
	// CurtainAirbags - наличие подушек безопасности-шторок
	CurtainAirbags Availability `db:"curtain_airbags"`
}

// MultimediaSystems - системы мультимедиа
type MultimediaSystems struct {
	// OnBoardComputer - наличие бортового компьютера
	OnBoardComputer Availability `db:"on_board_computer"`
	// MP3Support - наличие поддержки MP3
	MP3Support Availability `db:"mp3_support"`
	// HandsFreeSupport - наличие поддержки Hands free
	HandsFreeSupport Availability `db:"hands_free_support"`
}

// Offering - сведения для покупателя
type Offering struct {
	// Price - цена, руб
	Price string `db:"price"`
	// Year - год выпуска
	Year int `db:"year"`
	// Kilometerage - пробег, км
	Kilometerage string `db:"kilometerage"`
	// PhotoURLs - фотографии
	PhotoURLs []string `db:"photo_urls"`
}

func NewCar() Car {
	car := Car{}
	car.FullName = UndefinedStr
	car.Description = Absent
	car.Generation = UndefinedStr
	car.TrimLevel = UndefinedStr
	car.Specs.Body = UndefinedStr
	car.Specs.DragCoefficient = 0.29
	car.Specs.Gearbox = UndefinedStr
	car.Specs.Drive = UndefinedStr
	car.Specs.Engine.FuelUsed = UndefinedStr
	car.Specs.Engine.EngineType = UndefinedStr
	car.Specs.Engine.MaxTorque = UndefinedStr
	car.Specs.SteeringWheel.SteeringWheelPosition = UndefinedPos
	car.Specs.SteeringWheel.PowerSteering = UndefinedPS
	car.Specs.Suspension.FrontSuspension = UndefinedStr
	car.Specs.Suspension.FrontStabilizer = UndefinedValue
	car.Specs.Suspension.BackSuspension = UndefinedStr
	car.Specs.Suspension.BackStabilizer = UndefinedValue
	car.Specs.Brakes.FrontBrakes = UndefinedStr
	car.Specs.Brakes.BackBrakes = UndefinedStr
	car.Specs.Brakes.ParkingBrake = UndefinedStr
	car.Features.SafetyAndMotionControlSystem.ABS = UndefinedValue
	car.Features.SafetyAndMotionControlSystem.ESP = UndefinedValue
	car.Features.SafetyAndMotionControlSystem.EBD = UndefinedValue
	car.Features.SafetyAndMotionControlSystem.BAS = UndefinedValue
	car.Features.SafetyAndMotionControlSystem.TCS = UndefinedValue
	car.Features.SafetyAndMotionControlSystem.FrontParkingSensor = UndefinedValue
	car.Features.SafetyAndMotionControlSystem.BackParkingSensor = UndefinedValue
	car.Features.SafetyAndMotionControlSystem.RearViewCamera = UndefinedValue
	car.Features.SafetyAndMotionControlSystem.CruiseControl = UndefinedValue
	car.Features.Lights.Headlights = UndefinedStr
	car.Features.Lights.LEDRunningLights = UndefinedValue
	car.Features.Lights.LEDTailLights = UndefinedValue
	car.Features.Lights.LightSensor = UndefinedValue
	car.Features.Lights.FrontFogLights = UndefinedValue
	car.Features.Lights.BackFogLights = UndefinedValue
	car.Features.Interior.Upholstery = UndefinedStr
	car.Features.CabinMicroclimate.AirConditioner = UndefinedValue
	car.Features.CabinMicroclimate.ClimateControl = UndefinedValue
	car.Features.ElectricOptions.ElectricFrontSideWindowsLifts = UndefinedValue
	car.Features.ElectricOptions.ElectricBackSideWindowsLifts = UndefinedValue
	car.Features.ElectricOptions.ElectricHeatingOfFrontSeats = UndefinedValue
	car.Features.ElectricOptions.ElectricHeatingOfBackSeats = UndefinedValue
	car.Features.ElectricOptions.ElectricHeatingOfSteeringWheel = UndefinedValue
	car.Features.ElectricOptions.ElectricHeatingOfWindshield = UndefinedValue
	car.Features.ElectricOptions.ElectricHeatingOfRearWindow = UndefinedValue
	car.Features.ElectricOptions.ElectricHeatingOfSideMirrors = UndefinedValue
	car.Features.ElectricOptions.ElectricDriveOfDriverSeat = UndefinedValue
	car.Features.ElectricOptions.ElectricDriveOfFrontSeats = UndefinedValue
	car.Features.ElectricOptions.ElectricDriveOfSideMirrors = UndefinedValue
	car.Features.ElectricOptions.ElectricTrunkOpener = UndefinedValue
	car.Features.ElectricOptions.RainSensor = UndefinedValue
	car.Features.Airbags.DriverAirbag = UndefinedValue
	car.Features.Airbags.FrontPassengerAirbag = UndefinedValue
	car.Features.Airbags.SideAirbags = UndefinedValue
	car.Features.Airbags.CurtainAirbags = UndefinedValue
	car.Features.MultimediaSystems.OnBoardComputer = UndefinedValue
	car.Features.MultimediaSystems.MP3Support = UndefinedValue
	car.Features.MultimediaSystems.HandsFreeSupport = UndefinedValue
	car.Features.CarAlarm = UndefinedValue
	car.Features.Color = UndefinedStr
	car.Offering.Price = UndefinedStr
	car.Offering.Kilometerage = UndefinedStr
	return car
}
