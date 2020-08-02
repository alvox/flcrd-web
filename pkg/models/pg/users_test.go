package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
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
	u := &models.User{
		Name:   "Test",
		Email:  "test_email_1@example.com",
		Status: "PENDING",
	}
	c := &models.Credentials{
		Password: "testpass",
		Token: models.Token{
			AccessToken:     "authtoken",
			RefreshToken:    "refreshtokenz",
			RefreshTokenExp: parseTime("2019-12-01T22:08:41+00:00", t),
		},
	}
	_, err := model.Create(u, c)
	assert.Nil(t, err, "failed to create new user")
	user, err := model.GetByEmail("test_email_1@example.com")
	assert.Nil(t, err)
	if assert.NotNil(t, user) {
		assert.Equal(t, "Test", user.Name, "invalid user name")
		assert.Equal(t, "test_email_1@example.com", user.Email, "invalid user email")
		assert.Equal(t, "PENDING", user.Status, "invalid user status")
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
		Name:   "Test",
		Email:  "testuser1@example.com",
		Status: "PENDING",
	}
	c := &models.Credentials{}
	_, err := model.Create(u, c)
	assert.NotNil(t, err, "expect unique constraint violation")
	if err, ok := err.(*pgconn.PgError); ok {
		assert.Equal(t, "user_email_idx", err.ConstraintName, "expect unique constraint violation")
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
			userId: "40afbc9a-27e3-4b38-97f9-2930b8790a9f",
			wantUser: &models.User{
				ID:      "40afbc9a-27e3-4b38-97f9-2930b8790a9f",
				Name:    "Testuser1",
				Email:   "testuser1@example.com",
				Status:  "ACTIVE",
				Created: time.Date(2019, 1, 1, 9, 0, 0, 0, time.UTC),
				Token: models.Token{
					AccessToken:     "",
					RefreshToken:    "refreshtoken1",
					RefreshTokenExp: time.Date(2019, 2, 2, 10, 0, 0, 0, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:     "Non-existent user",
			userId:   "4735bdb2-45a5-42d9-a6d4-6db29787b5f1",
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
	c := &models.Credentials{
		UserID: "dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570",
		Token: models.Token{
			RefreshToken:    "newnew",
			RefreshTokenExp: parseTime("2019-04-17T11:57:00+00:00", t),
		},
	}
	err := model.UpdateRefreshToken(c)
	assert.Nil(t, err, "failed to update refresh token")

	c, err = model.GetCredentials("dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570")
	assert.Nil(t, err, "failed to read test user")
	assert.Equal(t, "newnew", c.Token.RefreshToken, "invalid refresh token")
	assert.Equal(t, time.Date(2019, 4, 17, 11, 57, 0, 0, time.UTC), c.Token.RefreshTokenExp, "invalid exp time")
}

func TestUserModel_UpdateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := UserModel{db}
	u := &models.User{
		ID:     "dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570",
		Name:   "updated name",
		Email:  "updated@email.com",
		Status: "ACTIVE",
	}
	err := model.Update(u)
	assert.Nil(t, err, "failed to update user status")
	user, err := model.Get("dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570")
	assert.Nil(t, err, "failed to read test user")
	if assert.NotNil(t, user) {
		assert.Equal(t, "ACTIVE", user.Status, "invalid user status")
		assert.Equal(t, "updated name", user.Name, "invalid user name")
		assert.Equal(t, "updated@email.com", user.Email, "invalid user email")
	}
}

func TestUserModel_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := UserModel{db}
	user, err := model.Get("40afbc9a-27e3-4b38-97f9-2930b8790a9f")
	assert.Nil(t, err)
	assert.NotNil(t, user, "can't find test user")
	err = model.Delete("40afbc9a-27e3-4b38-97f9-2930b8790a9f")
	assert.Nil(t, err)
	_, err = model.Get("40afbc9a-27e3-4b38-97f9-2930b8790a9f")
	assert.Equal(t, models.ErrNoRecord, err, "user still in the database")
	_, err = model.GetCredentials("40afbc9a-27e3-4b38-97f9-2930b8790a9f")
	assert.Equal(t, models.ErrNoRecord, err, "credentials still in the database")

	deckModel := DeckModel{db}
	decks, err := deckModel.GetForUser("40afbc9a-27e3-4b38-97f9-2930b8790a9f")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(decks))
}
