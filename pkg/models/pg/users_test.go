package pg

import (
	"github.com/lib/pq"
	"testing"
)

func TestUserModel_Create_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := UserModel{db}
	_, err := model.Create("Test", "test_email_1@example.com", "some_password")
	if err != nil {
		t.Error("Failed to create new user")
	}
	user, err := model.GetByEmail("test_email_1@example.com")
	if err != nil {
		t.Error("failed to read created test user")
	}
	if user.Name != "Test" {
		t.Errorf("invalid user name: %s", user.Name)
	}
	if user.Email != "test_email_1@example.com" {
		t.Errorf("invalid user email: %s", user.Email)
	}
	if user.Password != "some_password" {
		t.Errorf("invalid user password: %s", user.Password)
	}
}

func TestUserModel_Create_Email_Exists(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := UserModel{db}
	_, err := model.Create("Test", "testuser@example.com", "some_password")
	if err == nil {
		t.Error("expect unique constraint violation")
	}
	if err, ok := err.(*pq.Error); ok {
		if "unique_violation" != err.Code.Name() {
			t.Error("expect unique constraint violation")
		}
	}
}
