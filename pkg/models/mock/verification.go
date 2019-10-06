package mock

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"time"
)

var mockVerificationCode = &models.VerificationCode{
	UserID:  "test_user_id",
	Code:    "test_code",
	CodeExp: time.Now(),
}

type VerificationModel struct{}

func (m *VerificationModel) Create(code models.VerificationCode) (string, error) {
	return "test_code", nil
}

func (m *VerificationModel) Get(code string) (*models.VerificationCode, error) {
	return mockVerificationCode, nil
}

func (m *VerificationModel) GetForUser(userID string) (*models.VerificationCode, error) {
	return mockVerificationCode, nil
}

func (m *VerificationModel) Delete(code models.VerificationCode) error {
	return nil
}
