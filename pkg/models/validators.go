package models

import (
	"regexp"
	"unicode/utf8"
)

var r = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//todo: tests
func (u User) Validate(validateName bool) *ValidationErrors {
	errs := validateEmail(u.Email)

	if validateName && u.Name == "" {
		errs.Add("name", "field is required")
	}
	if validateName && u.Name != "" && utf8.RuneCountInString(u.Name) > 50 {
		errs.Add("name", "max length is 50 characters")
	}
	return errs
}

func (u User) ValidateForUpdate() *ValidationErrors {
	errs := validateEmail(u.Email)

	if u.Name == "" {
		errs.Add("name", "field is required")
	}
	if u.Name != "" && utf8.RuneCountInString(u.Name) > 50 {
		errs.Add("name", "max length is 50 characters")
	}
	return errs
}

func validateEmail(email string) *ValidationErrors {
	errs := NewValidationErrors()

	if email == "" {
		errs.Add("email", "field is required")
	}
	if email != "" && utf8.RuneCountInString(email) > 120 {
		errs.Add("email", "max length is 120 characters")
	}
	if email != "" && !r.MatchString(email) {
		errs.Add("email", "invalid email address")
	}
	return errs
}

//todo: tests
func (d Deck) Validate() *ValidationErrors {
	errs := NewValidationErrors()

	if d.Name == "" {
		errs.Add("name", "field is required")
	}
	if utf8.RuneCountInString(d.Name) > 50 {
		errs.Add("name", "max length is 50 characters")
	}
	if utf8.RuneCountInString(d.Description) > 250 {
		errs.Add("description", "max length is 250 characters")
	}
	return errs
}

func (d Deck) ValidateWithID(id string) *ValidationErrors {
	errs := d.Validate()

	if d.ID == "" {
		errs.Add("id", "field is required")
	}
	if d.ID != id {
		errs.Add("id", "deck id doesn't match with path param")
	}
	return errs
}

//todo: tests
func (f Flashcard) Validate() *ValidationErrors {
	errs := NewValidationErrors()

	if f.FrontType != "TEXT" && f.FrontType != "IMAGE" {
		errs.Add("front_type", "should be either 'TEXT' or 'IMAGE'")
	}
	if f.RearType != "TEXT" && f.RearType != "IMAGE" {
		errs.Add("rear_type", "should be either 'TEXT' or 'IMAGE'")
	}
	if f.FrontType == "TEXT" && f.Front == "" {
		errs.Add("front", "field is required")
	}
	if utf8.RuneCountInString(f.Front) > 250 {
		errs.Add("front", "max length is 250 characters")
	}
	if f.RearType == "TEXT" && f.Rear == "" {
		errs.Add("rear", "field is required")
	}
	if utf8.RuneCountInString(f.Rear) > 250 {
		errs.Add("rear", "max length is 250 characters")
	}
	return errs
}

func (t Token) Validate() *ValidationErrors {
	errs := NewValidationErrors()

	if t.AccessToken == "" {
		errs.Add("access_token", "field is required")
	}
	if t.RefreshToken == "" {
		errs.Add("refresh_token", "field is required")
	}
	return errs
}

func (a AuthRequest) Validate(validateName bool) *ValidationErrors {
	errs := validateEmail(a.Email)

	if validateName && a.Name == "" {
		errs.Add("name", "field is required")
	}
	if validateName && a.Name != "" && utf8.RuneCountInString(a.Name) > 50 {
		errs.Add("name", "max length is 50 characters")
	}
	if a.Password == "" {
		errs.Add("password", "field is required")
	}
	if a.Password != "" && utf8.RuneCountInString(a.Password) < 5 || utf8.RuneCountInString(a.Password) > 50 {
		errs.Add("password", "min length is 5 characters, max length is 50 characters")
	}
	return errs
}

type ValidationErrors struct {
	Errors []*ValidationError `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{}
}

func (e *ValidationErrors) Add(field, msg string) {
	e.Errors = append(e.Errors, &ValidationError{
		Field:   field,
		Message: msg,
	})
}

func (e ValidationErrors) Present() bool {
	return len(e.Errors) != 0
}
