// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type NewDocument struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
