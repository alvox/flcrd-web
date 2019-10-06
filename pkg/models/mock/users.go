package mock

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var mockUser = &models.User{
	ID:       "test_user_id",
	Name:     "test_user_name",
	Email:    "test_user_email@example.com",
	Password: hashPassword("test_password"),
	Created:  time.Now(),
	Token: models.Token{
		AccessToken:     "test_access_token",
		RefreshToken:    "test_refresh_token",
		RefreshTokenExp: time.Now(),
	},
}

type UserModel struct{}

func (m *UserModel) Create(user *models.User) (*string, error) {
	return &mockUser.ID, nil
}

func (m *UserModel) Get(userID string) (*models.User, error) {
	switch userID {
	case "test_user_id":
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) GetProfile(userID string) (*models.User, error) {
	switch userID {
	case "test_user_id":
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) GetByEmail(email string) (*models.User, error) {
	switch email {
	case "test_user_email@example.com":
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) UpdateRefreshToken(user *models.User) error {
	switch user.ID {
	case "test_user_id":
		return nil
	default:
		return models.ErrNoRecord
	}
}

func (m *UserModel) Update(user *models.User) error {
	return nil
}

func (m *UserModel) Delete(userID string) error {
	return nil
}

func hashPassword(pwd string) string {
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}
