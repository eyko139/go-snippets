package util

import "github.com/eyko139/go-snippets/internal/validator"

type SnippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

type UserSignupForm struct {
    Name string `form:"name"`
    Email string `form:"email"`
    Password string `form:"password"`
    validator.Validator `form:"-"`
}
