package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRegister_Positive(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	req := `{
            "name": "test", 
            "email": "test@test.com", 
            "password": "123456"
            }`
	status, _, resp := ts.post(t, "/v0/users/register", req)
	u, valid := parseUser(string(resp))
	if !valid {
		t.Error()
	}
	require.Equal(t, http.StatusCreated, status)
	require.Equal(t, "test_user_email@example.com", u.Email)
	require.Equal(t, "test_user_name", u.Name)
	require.NotEmpty(t, u.Token.RefreshToken)
	require.NotEmpty(t, u.Token.AccessToken)
}

func TestRegister_EmptyRequest(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, resp := ts.post(t, "/v0/users/register", "")
	wr := `{"code":"004","message":"can't read request"}`
	require.Equal(t, http.StatusBadRequest, status)
	require.Equal(t, []byte(wr), resp)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	req := `{
            "name": "test", 
            "email": "test_user_email@example.com", 
            "password": "123456"
            }`
	status, _, resp := ts.post(t, "/v0/users/register", req)
	wr := `{"code":"007","message":"user with this email already registered"}`
	require.Equal(t, http.StatusBadRequest, status)
	require.Equal(t, []byte(wr), resp)
}
