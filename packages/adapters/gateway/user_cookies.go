package gateway

import (
	"fmt"
	"vehicles/packages/adapters"
	"vehicles/packages/usecases/repository"
)

type userRepository struct {
	// ctx - переменная контекста
	ctx adapters.Context
}

func NewUserRepository(ctx adapters.Context) repository.UserRepository {
	return &userRepository{ctx}
}

// SetUserIDInCookie сохраняет в cookie идентификатор пользователя
// Входной параметр: userID - идентификатор пользователя
func (usr userRepository) SetUserIDInCookie(userID string) {
	usr.ctx.SetCookie("userID", userID, 0, "/", "localhost", false, true)
}

// GetUserIDFromCookie получает из cookie идентификатор пользователя
func (usr userRepository) GetUserIDFromCookie() (string, error) {
	cookieName := "userID"
	userID, err := usr.ctx.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("error from `Cookie` method, package `context`: %#v", err)
	}
	return userID, nil
}
