package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/stretchr/testify/require"
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
		UserID:  "40afbc9a-27e3-4b38-97f9-2930b8790a9f",
		Code:    "testcode",
		CodeExp: parseTime("2019-04-17T11:57:00+00:00", t),
	}
	c, err := model.Create(code)
	require.Nil(t, err)
	require.Equal(t, "testcode", c)
}

func TestVerificationModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := VerificationModel{db}
	c, err := model.Get("code_for_user_2")
	require.Nil(t, err)
	require.Equal(t, "dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570", c.UserID)
}

func TestVerificationModel_GetForUser(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := VerificationModel{db}
	c, err := model.GetForUser("dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570")
	require.Nil(t, err)
	require.Equal(t, "dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570", c.UserID)
	require.Equal(t, "code_for_user_2", c.Code)
}

func TestVerificationModel_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := VerificationModel{db}
	c := models.VerificationCode{
		UserID: "dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570",
		Code:   "code_for_user_2",
	}
	err := model.Delete(c)
	require.Nil(t, err)
	_, err = model.Get("code_for_user_2")
	require.Equal(t, models.ErrNoRecord, err)
}
