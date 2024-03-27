package usecase

import (
	"fmt"
	"math/rand"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecases/repository"
)

type QuestionInput interface {
	PickQuestion() (string, []string, error)
	SetQuestionID(questionID string)
	GetQuestionID() (string, error)
	GetAnswer(sessionID, answer string) error
}

type questionUseCase struct {
	questionRepo repository.QuestionRepository
	carsRepo     repository.CarsRepository
	userUseCase  UserInput
	output       SearchOutput
}

func NewQuestionUseCase(r repository.QuestionRepository, cr repository.CarsRepository, u UserInput, o SearchOutput) QuestionInput {
	return &questionUseCase{r, cr, u, o}
}

// PickQuestion выбирает вопрос для пользователя
// Ответ на этот вопрос необходим для "обучения" нечеткого алгоритма
func (qnu *questionUseCase) PickQuestion() (string, []string, error) {
	user := models.User{}
	var err error
	user.ID, err = qnu.userUseCase.GetUserID()
	if err != nil {
		return "", nil, fmt.Errorf("error from `GetUserID` method, package `usecase`: %#v", err)
	}

	questionIDs, err := qnu.questionRepo.GetIdsOfUnansweredQuestions(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("error from `GetIdsOfUnansweredQuestions` method, package `gateway`: %#v", err)
	}

	randIndex := rand.Intn(len(questionIDs))
	qtn := models.Question{}
	qtn.ID = questionIDs[randIndex]

	var possibleAnswers []string
	qtn.Question, possibleAnswers, err = qnu.questionRepo.GetQuestion(qtn.ID)
	if err != nil {
		return "", nil, fmt.Errorf("error from `GetQuestion` method, package `gateway`: %#v", err)
	}

	qnu.questionRepo.SetQuestionID(qtn.ID)

	return qtn.Question, possibleAnswers, nil
}

// GetAnswer принимает ответ пользователя. Этот ответ необходим
// для "обучения" нечеткого алгоритма
func (qnu *questionUseCase) GetAnswer(sessionID, answer string) error {
	user := models.User{}
	var err error
	user.ID, err = qnu.userUseCase.GetUserID()
	if err != nil {
		return fmt.Errorf("error from `GetUserID` method, package `usecase`: %#v", err)
	}

	qtn := models.Question{}
	qtn.ID, err = qnu.questionRepo.GetQuestionID()
	if err != nil {
		return fmt.Errorf("error from `GetQuestionID` method, package `gateway`: %#v", err)
	}
	qtn.Answer = answer

	err = qnu.questionRepo.InsertAnswer(user.ID, qtn.ID, qtn.Answer)
	if err != nil {
		return fmt.Errorf("error from `InsertAnswer` method, package `gateway`: %#v", err)
	}

	cars, err := qnu.carsRepo.GetCarsData(sessionID)
	if err != nil {
		return fmt.Errorf("error from `GetCarsData` method, package `gateway`: %#v", err)
	}

	qnu.output.ShowCars(sessionID, cars)
	return nil
}

// SetQuestionID сохраняет идентификатор вопроса в cookie
func (qnu *questionUseCase) SetQuestionID(questionID string) {
	qtn := models.Question{}
	qtn.ID = questionID
	qnu.questionRepo.SetQuestionID(qtn.ID)
}

// GetQuestionID получает идентификатор вопроса из cookie
func (qnu *questionUseCase) GetQuestionID() (string, error) {
	qtn := models.Question{}
	var err error
	qtn.ID, err = qnu.questionRepo.GetQuestionID()
	if err != nil {
		return "", fmt.Errorf("error from `GetQuestionID` method, package `gateway`: %#v", err)
	}
	return qtn.ID, nil
}
