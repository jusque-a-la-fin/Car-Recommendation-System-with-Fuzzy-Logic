package gateway

import (
	"database/sql"
	"fmt"
	"vehicles/packages/adapters"
	"vehicles/packages/usecases/repository"
)

type questionRepository struct {
	// ctx - переменная контекста
	ctx adapters.Context
	// questionsDB - клиент для подключения к реляционной БД под управлением PostgreSQL,
	// хранящей информацию, связанную с опросом пользователей
	questionsDB *sql.DB
}

func NewQuestionRepository(ctx adapters.Context, questionsDB *sql.DB) repository.QuestionRepository {
	return &questionRepository{ctx, questionsDB}
}

// GetIdsOfUnansweredQuestions получает из БД id вопросов, на которые пользователь ещё не отвечал
// Входной параметр: userID - идентификатор пользователя
func (qnr *questionRepository) GetIdsOfUnansweredQuestions(userID string) ([]string, error) {
	query := `
        SELECT id
        FROM questions
        WHERE id NOT IN (
          SELECT question_id
          FROM user_responses
          WHERE user_id = (
            SELECT id
            FROM users
            WHERE hashed_ip = $1
          )
        );
    `

	rows, err := qnr.questionsDB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error from `Query` method, package `sql`: %#v", err)
	}
	defer rows.Close()

	var questionIDs []string
	var id string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		questionIDs = append(questionIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error from `Err` method, package `sql`: %#v", err)
	}
	return questionIDs, nil
}

// GetQuestion получает вопрос для пользователя из БД
// Входной параметр: questionID - идентификатор вопроса
func (qnr *questionRepository) GetQuestion(questionID string) (string, []string, error) {
	query := `
        SELECT question
        FROM questions
        WHERE id = $1;
    `

	var questionText string
	err := qnr.questionsDB.QueryRow(query, questionID).Scan(&questionText)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, fmt.Errorf("error: sql.ErrNoRows from `QueryRow` method, package `sql`: %#v", sql.ErrNoRows)
		} else {
			return "", nil, fmt.Errorf("error from `QueryRow` method, package `sql`: %#v", err)
		}
	}

	query = `
        SELECT possible_answer
        FROM possible_answers
        WHERE question_id = $1;
    `

	rows, err := qnr.questionsDB.Query(query, questionID)
	if err != nil {
		return "", nil, fmt.Errorf("error from `Query` method, package `sql`: %#v", err)
	}
	defer rows.Close()

	var possibleAnswers []string
	var possibleAnswer string
	for rows.Next() {
		err := rows.Scan(&possibleAnswer)
		if err != nil {
			return "", nil, fmt.Errorf("error from `Scan` method, package `sql`: %#v", err)
		}
		possibleAnswers = append(possibleAnswers, possibleAnswer)
	}
	if err := rows.Err(); err != nil {
		return "", nil, fmt.Errorf("error from `Err` method, package `sql`: %#v", err)
	}

	return questionText, possibleAnswers, nil
}

// InsertAnswer записывает ответ пользователя в БД
// Входные параметры: userID - идентификатор пользователя, questionID - идентификатор вопроса,
// answer - ответ пользователя
func (qnr *questionRepository) InsertAnswer(userID, questionID, answer string) error {
	var exists bool
	err := qnr.questionsDB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE hashed_ip = $1)", userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error from `QueryRow` method, package `sql`: %#v", err)
	}

	if !exists {

		sqlQuery := `
        WITH new_user AS (
            INSERT INTO users (hashed_ip)
            SELECT CAST($1 AS VARCHAR(160))
            WHERE NOT EXISTS (
                SELECT 1 FROM users WHERE hashed_ip = $1
        )
            RETURNING id
        ),
        upsert_response AS (
            INSERT INTO user_responses (user_id, question_id, answer)
            SELECT new_user.id, $2, $3
            FROM new_user
		    ON CONFLICT ON CONSTRAINT unique_user_question_responses
            DO UPDATE SET answer = EXCLUDED.answer
            RETURNING *
        )
        SELECT * FROM upsert_response`

		_, err := qnr.questionsDB.Exec(sqlQuery, userID, questionID, answer)
		if err != nil {
			return fmt.Errorf("error while inserting answer for non-existing user, error from `Exec` method, package `sql`: %#v", err)
		}

	} else {
		sqlQuery := `
          WITH existing_user AS (
	          SELECT id FROM users WHERE hashed_ip = $1
          ),
          upsert_response AS (
	          INSERT INTO user_responses (user_id, question_id, answer)
	          SELECT existing_user.id, $2, $3
	          FROM existing_user
	          ON CONFLICT ON CONSTRAINT unique_user_question_responses
	          DO UPDATE SET answer = EXCLUDED.answer
	          RETURNING *
           ) 
           SELECT * FROM upsert_response`

		_, err := qnr.questionsDB.Exec(sqlQuery, userID, questionID, answer)
		if err != nil {
			return fmt.Errorf("error while inserting answer for existing user, error from `Exec` method, package `sql`: %#v", err)
		}
	}
	return nil
}
