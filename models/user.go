package models

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var (
	PwPepper              = os.Getenv("PASSWORD_PEPPER")
	ErrInvalidCredentials = errors.New("Invalid username or password")
	RememberTokenLength   = 32
)

type UserService struct {
	db   *sql.DB
	hmac hash.Hash
}

func NewUserService(db *sql.DB) *UserService {
	hmacSecret := os.Getenv("HMAC_SECRET_KEY")
	return &UserService{
		db:   db,
		hmac: hmac.New(sha256.New, []byte(hmacSecret)),
	}
}

func (u *UserService) Authenticate(email, password string) (*User, error) {
	user, err := u.GetBy("email", email)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password+PwPepper))
	switch err {
	case nil:
		return user, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidCredentials
	default:
		return nil, err
	}
}

func (u *UserService) GenerateRememberToken() (string, error) {
	b := make([]byte, RememberTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (u *UserService) hashRememberToken(token string) (string, error) {
	u.hmac.Reset()
	_, err := u.hmac.Write([]byte(token))
	if err != nil {
		return "", err
	}
	b := u.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b), nil
}

func (u *UserService) GetBy(field string, val interface{}) (*User, error) {
	sqlStatement := fmt.Sprintf(`
		SELECT id, name, email, password_hash, remember_hash
		FROM users
		WHERE %s = $1;
	`, field)
	var user User
	row := u.db.QueryRow(sqlStatement, val)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RememberHash)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserService) ByRemember(token string) (*User, error) {
	rememberHash, err := u.hashRememberToken(token)
	if err != nil {
		return nil, err
	}
	return u.GetBy("remember_hash", rememberHash)
}

func (u *UserService) Create(input NewUser) (*User, error) {
	remember, err := u.GenerateRememberToken()
	if err != nil {
		return nil, err
	}
	user := &User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Remember: remember,
	}
	err = user.HashPassword()
	if err != nil {
		return nil, err
	}
	rememberHash, err := u.hashRememberToken(user.Remember)
	if err != nil {
		return nil, err
	}
	user.RememberHash = rememberHash
	sqlStatement := `
		INSERT INTO users (name, email, password_hash, remember_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	id := ""
	err = u.db.QueryRow(sqlStatement, user.Name, user.Email, user.PasswordHash, user.RememberHash).Scan(&id)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (u *UserService) Update(user *User) error {
	if user.Remember != "" {
		rememberHash, err := u.hashRememberToken(user.Remember)
		if err != nil {
			return err
		}
		user.RememberHash = rememberHash
	}
	sqlStatement := `
		UPDATE users
		SET name = $2,
				email = $3,
				password_hash = $4,
				remember_hash = $5
		WHERE id = $1
	`
	_, err := u.db.Exec(sqlStatement, user.ID, user.Name, user.Email, user.PasswordHash, user.RememberHash)
	return err
}

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string
	PasswordHash string
	Remember     string
	RememberHash string
}

func (u *User) HashPassword() error {
	pwBytes := []byte(u.Password + PwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		pwBytes,
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedBytes)
	u.Password = ""
	return nil
}
