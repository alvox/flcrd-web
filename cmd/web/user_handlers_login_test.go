package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestLogin_Positive(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	req := `{
            "email": "test_user_email@example.com", 
            "password": "test_password"
            }`
	status, _, resp := ts.post(t, "/v0/users/login", req)
	u, valid := parseUser(string(resp))
	if !valid {
		t.Error()
	}
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "test_user_email@example.com", u.Email)
	require.Equal(t, "test_user_name", u.Name)
	require.NotEmpty(t, u.Token.RefreshToken)
	require.NotEmpty(t, u.Token.AccessToken)
}

func TestLogin_EmptyRequest(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, resp := ts.post(t, "/v0/users/login", "")
	wr := `{"code":"004","message":"can't read request"}`
	require.Equal(t, http.StatusBadRequest, status)
	require.Equal(t, []byte(wr), resp)
}

func TestLogin_UserNotFound(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	req := `{
            "email": "non_existing_user@example.com", 
            "password": "123456"
            }`
	status, _, resp := ts.post(t, "/v0/users/login", req)
	wr := `{"code":"006","message":"email or password incorrect"}`
	require.Equal(t, http.StatusUnauthorized, status)
	require.Equal(t, []byte(wr), resp)
}

func TestLogin_InvalidPassword(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	req := `{
            "email": "test_user_email@example.com", 
            "password": "123456"
            }`
	status, _, resp := ts.post(t, "/v0/users/login", req)
	wr := `{"code":"006","message":"email or password incorrect"}`
	require.Equal(t, http.StatusUnauthorized, status)
	require.Equal(t, []byte(wr), resp)
}
