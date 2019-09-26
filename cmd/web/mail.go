package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type MailSender struct {
	baseUrl string
	apiKey  string
}

func (s *MailSender) SendConfirmation(to, link string) {
	msg, contentType, err := composeMessage(to, link)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", s.baseUrl, msg)
	if err != nil {
		fmt.Printf("Error while preparing email request: %s\n", err.Error())
		return
	}
	req.SetBasicAuth("api", s.apiKey)
	req.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error while sending email: %s\n", err.Error())
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error while reading mailgun response: %s\n", err.Error())
		return
	}
	_, err = io.Copy(ioutil.Discard, resp.Body)
	fmt.Println(string(body))
}

func composeMessage(to, link string) (*bytes.Buffer, string, error) {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	v := map[string]string{
		"from":    "alex@flashcards.rocks",
		"to":      to,
		"subject": "Confirm your email, please",
		"text":    fmt.Sprintf("Thank you for joining!\nPlease, click this link to complete your registration: %s", link),
	}
	for k, v := range v {
		if tmp, err := writer.CreateFormField(k); err == nil {
			_, err = tmp.Write([]byte(v))
			if err != nil {
				fmt.Println(err)
				return nil, "", err
			}
		} else {
			fmt.Println(err)
			return nil, "", err
		}
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}
	return data, writer.FormDataContentType(), nil
}
