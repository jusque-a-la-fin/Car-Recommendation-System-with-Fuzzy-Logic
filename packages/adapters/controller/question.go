package controller

import (
	"fmt"
	"vehicles/packages/adapters"
	usecase "vehicles/packages/usecases/usecases"
)

type questionController struct {
	ctx             adapters.Context
	questionUseCase usecase.QuestionInput
}

type Question interface {
	GetAnswer() error
}

func NewQuestionController(ctx adapters.Context, qni usecase.QuestionInput) Question {
	return &questionController{ctx, qni}
}

// GetAnswer получает ответ от пользователя, необходимый для "обучения" нечеткого алгоритма, и записывает его в базу данных
// Пользователь дает этот ответ в опросе на странице, где показан список автомобилей
func (qnc *questionController) GetAnswer() error {
	answer := qnc.ctx.PostForm("radio")
	sessionID := qnc.ctx.Query("guest")
	err := qnc.questionUseCase.GetAnswer(sessionID, answer)
	if err != nil {
		return fmt.Errorf("error from `GetAnswer` method, package `usecase`: %#v", err)
	}
	return nil
}
