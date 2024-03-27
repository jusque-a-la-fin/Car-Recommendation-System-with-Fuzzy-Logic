package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecases/repository"
)

type UserInput interface {
	SetUserID(userID string) error
	GetUserID() (string, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) UserInput {
	return &userUseCase{r}
}

// SetUserID сохраняет в cookie захешированный ip-адрес пользователя
// Входной параметр: userIP - ip-адрес пользователя
func (uru *userUseCase) SetUserID(userIP string) error {
	user := models.User{}
	var err error
	user.ID, err = generateHashOfUserID(userIP)
	if err != nil {
		return fmt.Errorf("error from `generateHashOfUserID` function, package `usecase`: %#v", err)
	}

	uru.userRepo.SetUserIDInCookie(user.ID)
	return nil
}

// generateHashOfUserID генерирует хэш ip-адреса пользователя
// Входной параметр: userIP - ip-адрес пользователя
func generateHashOfUserID(userIP string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(userIP))
	if err != nil {
		return "", fmt.Errorf("error from `Write` method, package `io`: %#v", err)
	}
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed), nil
}

// GetUserID получает захешированный ip-адрес пользователя из cookie
func (uru *userUseCase) GetUserID() (string, error) {
	user := models.User{}
	var err error
	user.ID, err = uru.userRepo.GetUserIDFromCookie()
	if err != nil {
		return "", fmt.Errorf("error from `GetUserIDFromCookie` method, package `gateway`: %#v", err)
	}
	return user.ID, nil
}
