package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"testing"
)

func TestVerificationModel_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := VerificationModel{db}
	code := models.VerificationCode{
		UserID:  "testuser_id_1",
		Code:    "testcode",
		CodeExp: parseTime("2019-04-17T11:57:00+00:00", t),
	}
	c, err := model.Create(code)
	if err != nil {
		t.Errorf("failed to create new verification code: %s", err.Error())
	}
	if c != "testcode" {
		t.Errorf("fnvalid verification code; want: %s, got: %s", "testcode", c)
	}
}

func TestVerificationModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := VerificationModel{db}
	c, err := model.Get("code_for_user_2")
	if err != nil {
		t.Errorf("failed to get verification code: %s", err.Error())
	}
	if c.UserID != "testuser_id_2" {
		t.Errorf("invalid user id; want: %s, got: %s", "testuser_id_2", c.UserID)
	}
}

func TestVerificationModel_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := VerificationModel{db}
	c := models.VerificationCode{
		UserID: "testuser_id_2",
		Code:   "code_for_user_2",
	}
	err := model.Delete(c)
	if err != nil {
		t.Errorf("failed to delete code: %s", err)
	}
	_, err = model.Get("code_for_user_2")
	if err != models.ErrNoRecord {
		t.Error("code hasn't been deleted in delete test")
	}
}
