package main

import (
	"bytes"
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
	if status != http.StatusOK {
		t.Errorf("status: want 200; got %d", status)
	}
	if u.Email != "test_user_email@example.com" {
		t.Errorf("email: want test_user_email@example.com; got %s", u.Email)
	}
	if u.Name != "test_user_name" {
		t.Errorf("name: want test_user_name; got %s", u.Name)
	}
	if u.Token.RefreshToken == "" {
		t.Errorf("refreshToken: empty")
	}
	if u.Token.AccessToken == "" {
		t.Errorf("accessToken: empty")
	}
	if u.Password != "" {
		t.Errorf("password: password should be empty; got %s", u.Password)
	}
}

func TestLogin_EmptyRequest(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, resp := ts.post(t, "/v0/users/login", "")

	if status != http.StatusBadRequest {
		t.Errorf("status: want 400; got %d", status)
	}
	wr := `{"code":"004","message":"can't read request body"}`
	if !bytes.Contains([]byte(wr), resp) {
		t.Errorf("response: want %s; got %s", wr, string(resp))
	}
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

	if status != http.StatusUnauthorized {
		t.Errorf("status: want 401; got %d", status)
	}
	wr := `{"code":"006","message":"email or password incorrect"}`
	if !bytes.Contains([]byte(wr), resp) {
		t.Errorf("response: want %s; got %s", wr, string(resp))
	}
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

	if status != http.StatusUnauthorized {
		t.Errorf("status: want 401; got %d", status)
	}
	wr := `{"code":"006","message":"email or password incorrect"}`
	if !bytes.Contains([]byte(wr), resp) {
		t.Errorf("response: want %s; got %s", wr, string(resp))
	}
}
