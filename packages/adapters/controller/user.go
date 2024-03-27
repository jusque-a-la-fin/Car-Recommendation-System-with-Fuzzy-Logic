package controller

import (
	"fmt"
	"vehicles/packages/adapters"
	usecase "vehicles/packages/usecases/usecases"

	"github.com/gin-gonic/gin"
)

type userController struct {
	ctx         adapters.Context
	userUseCase usecase.UserInput
}

type User interface {
	SaveUserID() error
}

func NewUserController(ctx adapters.Context, uri usecase.UserInput) User {
	return &userController{ctx, uri}
}

// SaveUserID сохраняет идентификатор пользователя
func (usc *userController) SaveUserID() error {
	userIP := usc.ctx.(*gin.Context).Request.RemoteAddr
	err := usc.userUseCase.SetUserID(userIP)
	if err != nil {
		return fmt.Errorf("error from `SetUserID` method, package `usecase`: %#v", err)
	}
	return nil
}
