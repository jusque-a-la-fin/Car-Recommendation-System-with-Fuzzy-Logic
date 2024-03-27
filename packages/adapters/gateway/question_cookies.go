package gateway

import "fmt"

// SetQuestionID сохраняет в cookie идентификатор вопроса
// Входной параметр: questionID - идентификатор вопроса
func (qnr *questionRepository) SetQuestionID(questionID string) {
	qnr.ctx.SetCookie("questionID", questionID, 0, "/", "localhost", false, true)
}

// GetQuestionID получает из cookie идентификатор вопроса
func (qnr *questionRepository) GetQuestionID() (string, error) {
	cookieName := "questionID"
	questionID, err := qnr.ctx.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("error from `Cookie` method, package `gin`: %#v", err)
	}

	return questionID, nil
}
