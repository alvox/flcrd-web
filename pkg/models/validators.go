package models

import (
	"net/url"
	"regexp"
)

var r = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//todo: tests
func (u User) Validate(validateName bool) url.Values {
	errs := url.Values{}

	if validateName && u.Name == "" {
		errs.Add("name", "field is required")
	}
	if validateName && len(u.Name) > 50 {
		errs.Add("name", "max length is 50 characters")
	}
	if u.Email == "" {
		errs.Add("email", "field is required")
	}
	if len(u.Email) > 120 {
		errs.Add("email", "max length is 120 characters")
	}
	if !r.MatchString(u.Email) {
		errs.Add("email", "invalid format")
	}
	if u.Password == "" {
		errs.Add("password", "field is required")
	}
	if len(u.Password) < 5 || len(u.Password) > 30 {
		errs.Add("password", "min length is 5 characters, max length is 30 characters")
	}
	return errs
}

//todo: tests
func (d Deck) Validate() url.Values {
	errs := url.Values{}

	if d.Name == "" {
		errs.Add("name", "field is required")
	}
	if len(d.Name) > 50 {
		errs.Add("name", "max length is 50 characters")
	}
	if len(d.Description) > 250 {
		errs.Add("description", "max length is 250 characters")
	}
	return errs
}

//todo: tests
func (f Flashcard) Validate() url.Values {
	errs := url.Values{}

	if f.DeckID == "" {
		errs.Add("deck_id", "field is required")
	}
	if f.Front == "" {
		errs.Add("front", "field is required")
	}
	if len(f.Front) > 250 {
		errs.Add("front", "max length is 250 characters")
	}
	if f.Rear == "" {
		errs.Add("rear", "field is required")
	}
	if len(f.Rear) > 250 {
		errs.Add("rear", "max length is 250 characters")
	}
	return errs
}
