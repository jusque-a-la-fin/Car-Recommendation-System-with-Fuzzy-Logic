package repository

type UserRepository interface {
	// SetUserIDInCookie сохраняет в cookie идентификатор пользователя
	// Входной параметр: userID - идентификатор пользователя
	SetUserIDInCookie(userID string)

	// GetUserIDFromCookie получает из cookie идентификатор пользователя
	GetUserIDFromCookie() (string, error)
}
