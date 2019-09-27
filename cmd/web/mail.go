package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type MailSender struct {
	baseUrl string
	apiKey  string
}

type SendMessageResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

func (s *MailSender) SendConfirmation(to, name, code string) (*SendMessageResponse, error) {
	msg, contentType, err := composeMessage(to, name, code)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.baseUrl, msg)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("api", s.apiKey)
	req.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	r, err := parseMailResponse(resp)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(ioutil.Discard, resp.Body)
	return r, nil
}

func composeMessage(to, name, code string) (*bytes.Buffer, string, error) {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	fields := map[string]string{
		"from":                "alex@flashcards.rocks",
		"to":                  to,
		"subject":             "Confirm your email, please",
		"template":            "verify_email",
		"v:user_name":         name,
		"v:verification_code": code,
	}
	for k, v := range fields {
		if tmp, err := writer.CreateFormField(k); err == nil {
			_, err = tmp.Write([]byte(v))
			if err != nil {
				return nil, "", err
			}
		} else {
			return nil, "", err
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, "", err
	}
	return data, writer.FormDataContentType(), nil
}

func parseMailResponse(resp *http.Response) (*SendMessageResponse, error) {
	var r SendMessageResponse
	err := json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
