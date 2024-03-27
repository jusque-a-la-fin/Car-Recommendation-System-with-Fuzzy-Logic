package repository

type QuestionRepository interface {
	// GetIdsOfUnansweredQuestions получает из БД id вопросов, на которые пользователь ещё не отвечал
	// Входной параметр: userID - идентификатор пользователя
	GetIdsOfUnansweredQuestions(userID string) ([]string, error)

	// GetQuestion получает вопрос для пользователя из БД
	// Входной параметр: questionID - идентификатор вопроса
	GetQuestion(questionID string) (string, []string, error)

	// InsertAnswer записывает ответ пользователя в БД
	// Входные параметры: userID - идентификатор пользователя, questionID - идентификатор вопроса,
	// answer - ответ пользователя
	InsertAnswer(userID, questionID, answer string) error

	// SetQuestionID сохраняет в cookie идентификатор вопроса
	// Входной параметр: questionID - идентификатор вопроса
	SetQuestionID(questionID string)

	// GetQuestionID получает из cookie идентификатор вопроса
	GetQuestionID() (string, error)
}
