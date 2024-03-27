package models

// Search - параметры запроса, по которым будут искаться автомобили
type Search struct {
	// Mark - марка
	Mark string `json:"mark"`
	// Model -модель
	Model string `json:"model"`
	// Gearbox - тип трансмиссии
	Gearbox string `json:"gearbox"`
	// LowPriceLimit - нижний предел цены
	LowPriceLimit string `json:"low_price_limit"`
	// HighPriceLimit - верхний предел цены
	HighPriceLimit string `json:"high_price_limit"`
	// Drive - тип привода
	Drive string `json:"drive"`
	// EarliestYear - самый ранний год выпуска
	EarliestYear string `json:"earliest_year"`
	// LatestYear - самый поздний год выпуска
	LatestYear string `json:"lastest_year"`
	// Fuel - тип топлива
	Fuel string `json:"fuel"`
	// IsNewCar - признак, определяющий, нужен ли пользователю новый автомобиль
	IsNewCar string `json:"new"`
}
