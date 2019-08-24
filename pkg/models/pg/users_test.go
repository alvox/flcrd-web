package pg

import (
	"testing"
)

func TestUserModel_Create_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := UserModel{db}
	id, err := model.Create("Test", "test_email_1@example.com", "some_password")
	if err != nil {
		t.Error("Failed to create new user")
	}
	user, err := model.Get(*id)
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
