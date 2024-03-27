package gateway

import (
	"database/sql"
	"fmt"
	"strings"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecases/repository"
)

type selectionRepository struct {
	// ctx - переменная контекста
	ctx adapters.Context
	// vehiclesDB - клиент для подключения к реляционной БД, хранящей датасет автомобилей
	vehiclesDB *sql.DB
}

func NewSelectionRepository(ctx adapters.Context, vehiclesDB *sql.DB) repository.SelectionRepository {
	return &selectionRepository{ctx, vehiclesDB}
}

// SetPriorities сохраняет приоритеты, расставленные пользователем, в cookie
// Приоритеты - это свойства: "Комфорт", "Экономичность", "Безопасность", "Динамика", "Управляемость",
// по которым будут ранжироваться автомобили
// Входной параметр: priorities - приоритеты, расставленные пользователем
func (slr *selectionRepository) SetPriorities(priorities []string) {
	prioritiesStr := strings.Join(priorities, ",")
	slr.ctx.SetCookie("priorities", prioritiesStr, 3600, "/", "localhost", false, true)
}

// SetPrice сохраняет диапазон цен, заданный пользователем, в cookies
// Входные параметры: minPrice - минимальная цена, maxPrice - максимальная цена
func (slr *selectionRepository) SetPrice(minPrice, maxPrice string) {
	slr.ctx.SetCookie("minPrice", minPrice, 3600, "/", "localhost", false, true)
	slr.ctx.SetCookie("maxPrice", maxPrice, 3600, "/", "localhost", false, true)
}

// SetManufacturers сохраняет названия стран-производителей, выбранных пользователем, в cookie
// Входной параметр: manufacturers - страны, выбранные пользователем
func (slr *selectionRepository) SetManufacturers(manufacturers []string) {
	manufacturersStr := strings.Join(manufacturers, ",")
	slr.ctx.SetCookie("manufacturers", manufacturersStr, 3600, "/", "localhost", false, true)
}

// GetSelectionParams получает параметры, заданные ранее пользователем, из cookies
func (slr *selectionRepository) GetSelectionParams() (*models.Selection, error) {
	slc := new(models.Selection)
	prioritiesStr, err := slr.ctx.Cookie("priorities")
	if err != nil {
		return nil, fmt.Errorf("error from `Cookie` method, package `gin`: %#v", err)
	}
	slc.Priorities = strings.Split(prioritiesStr, ",")

	slc.MinPrice, err = slr.ctx.Cookie("minPrice")
	if err != nil {
		return nil, fmt.Errorf("error from `Cookie` method, package `gin`: %#v", err)
	}

	slc.MaxPrice, err = slr.ctx.Cookie("maxPrice")
	if err != nil {
		return nil, fmt.Errorf("error from `Cookie` method, package `gin`: %#v", err)
	}

	manufacturersStr, err := slr.ctx.Cookie("manufacturers")
	if err != nil {
		return nil, fmt.Errorf("error from `Cookie` method, package `gin`: %#v", err)
	}
	slc.Manufacturers = strings.Split(manufacturersStr, ",")
	return slc, nil
}
