package models

type Selection struct {
	// Priorities - приоритеты или нечеткие множества, например,
	// "Экономичность", "Комфорт", "Управляемость", "Динамика", "Безопасность"
	Priorities []string
	// MinPrice - нижний предел цены
	MinPrice string
	// MaxPrice - верхний предел цены
	MaxPrice string
	// Manufacturers - страны-производители
	Manufacturers []string
}

type Makes struct {
	// Make - название марки
	Make string
	// NumberOfCars - количество авто этой марки
	NumberOfCars int
}
