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
	u := &models.User{
		Email:      "test_email_1@example.com",
		ExternalID: "some_external_id",
	}
	_, err := model.Create(u)
	if err != nil {
		t.Error("Failed to create new user")
	}
	user, err := model.GetByEmail("test_email_1@example.com")
	if err != nil {
		t.Error("failed to read created test user")
	}
	if user.Email != "test_email_1@example.com" {
		t.Errorf("invalid user email: %s", user.Email)
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
		Email:      "testuser1@example.com",
		ExternalID: "some_extenal_id",
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
				ID:         "testuser_id_1",
				Email:      "testuser1@example.com",
				ExternalID: "12345",
				Created:    time.Date(2019, 1, 1, 9, 0, 0, 0, time.UTC),
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

func TestUserModel_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := UserModel{db}
	user, err := model.Get("testuser_id_1")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if user == nil {
		t.Errorf("can't find test user testuser_id_1")
	}
	err = model.Delete("testuser_id_1")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	_, err = model.Get("testuser_id_1")
	if err != models.ErrNoRecord {
		t.Errorf("user haven't been deleted; want: %s, got: %s", models.ErrNoRecord, err)
	}
	deckModel := DeckModel{db}
	decks, err := deckModel.GetForUser("testuser_id_1")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(decks) != 0 {
		t.Errorf("decks count invalid; want: 0, got: %d", len(decks))
	}
}
