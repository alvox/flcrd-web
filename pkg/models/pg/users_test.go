package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/lib/pq"
	"reflect"
	"testing"
	"time"
)

func TestUserModel_Create_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := UserModel{db}

	exp, e := time.Parse(
		time.RFC3339,
		"2019-12-01T22:08:41+00:00")
	if e != nil {
		t.Errorf("unexpected error while preparing test data: %s", e.Error())
	}

	u := &models.User{
		Name:     "Test",
		Email:    "test_email_1@example.com",
		Password: "some_password",
		Token: models.Token{
			AccessToken:     "authtoken",
			RefreshToken:    "refreshtoken",
			RefreshTokenExp: exp,
		},
	}

	_, err := model.Create(u)
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

	u := &models.User{
		Name:     "Test",
		Email:    "testuser1@example.com",
		Password: "some_password",
		Token: models.Token{
			AccessToken:  "authtoken",
			RefreshToken: "refreshtoken",
		},
	}

	_, err := model.Create(u)
	if err == nil {
		t.Error("expect unique constraint violation")
	}
	if err, ok := err.(*pq.Error); ok {
		if "unique_violation" != err.Code.Name() {
			t.Error("expect unique constraint violation")
		}
	}
}

func TestUserModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	tests := []struct {
		name     string
		userId   string
		wantUser *models.User
		wantErr  error
	}{
		{
			name:   "User 1",
			userId: "testuser_id_1",
			wantUser: &models.User{
				ID:       "testuser_id_1",
				Name:     "Testuser1",
				Email:    "testuser1@example.com",
				Password: "12345",
				Created:  time.Date(2019, 1, 1, 9, 0, 0, 0, time.UTC),
				Token: models.Token{
					AccessToken:     "",
					RefreshToken:    "refreshtoken",
					RefreshTokenExp: time.Date(2019, 2, 2, 10, 0, 0, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:     "Non-existent user",
			userId:   "non-existing-testuser",
			wantUser: nil,
			wantErr:  models.ErrNoRecord,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := UserModel{db}
			user, err := model.Get(tt.userId)
			if err != tt.wantErr {
				t.Errorf("want %v; got %s", tt.wantErr, err)
			}
			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %v; got %v", tt.wantUser, user)
			}
		})
	}
}

func TestUserModel_UpdateRefreshToken(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}

	db, teardown := newTestDB(t)
	defer teardown()
	model := UserModel{db}

	exp, e := time.Parse(
		time.RFC3339,
		"2019-04-17T11:57:00+00:00")
	if e != nil {
		t.Errorf("unexpected error while preparing test data: %s", e.Error())
	}

	u := &models.User{
		ID: "testuser_id_2",
		Token: models.Token{
			RefreshToken:    "newnew",
			RefreshTokenExp: exp,
		},
	}
	err := model.UpdateRefreshToken(u)
	if err != nil {
		t.Errorf("failed to update refresh token: %s", err.Error())
	}
	user, err := model.Get("testuser_id_2")
	if err != nil {
		t.Errorf("failed to read test user: %s", err.Error())
	}
	if user.Token.RefreshToken != "newnew" {
		t.Errorf("invalid refresh token: %s", user.Token.RefreshToken)
	}
	if user.Token.RefreshTokenExp != time.Date(2019, 4, 17, 11, 57, 0, 0, time.UTC) {
		t.Errorf("invalid refresh token expiration time: %s", user.Token.RefreshTokenExp)
	}

}
