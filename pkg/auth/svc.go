package auth

import (
	"context"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// normalize normalizes email address.
func normalize(email string) string {
	// Trim whitespaces.
	email = strings.TrimSpace(email)

	// Trim extra dot in hostname.
	email = strings.TrimRight(email, ".")

	// Lowercase.
	email = strings.ToLower(email)

	return email
}

// Service for auth
type Service struct {
	repo Repository
}

// NewService returns a new AUTH Service initialized with
// a concrete repo implementation
func NewService(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

// IsValidEmail checks if an email is valid or not
func (as Service) IsValidEmail(email string) bool {
	// email addresses have a practical limit of 254 bytes
	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}

	// {64}@{255}
	at := strings.LastIndex(email, "@")
	user := email[:at]
	isLenOk := len(user) > 64
	return !isLenOk
}

// IsValidPassword validate if the password is valid or not
func (as Service) IsValidPassword(password string) bool {
	// min 8 character only and less than 254
	if len(strings.TrimSpace(password)) > 254 || len(strings.TrimSpace(password)) < 8 {
		return false
	}
	return true
}

// IsCredentialValid checks if the Credential is ok
func (as Service) IsCredentialValid(ctx context.Context, email string,
	password string) (bool, error) {
	user, err := as.repo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if user.Email != email {
		return false, nil
	}
	// if password is invalid no need to compare invoke compare hash
	if !as.IsValidPassword(password) {
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

// IsDuplicateRegistration checks if the user is already registered
func (as Service) IsDuplicateRegistration(ctx context.Context, email string) (bool,
	error) {
	user, err := as.repo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if user.Email != email {
		return false, nil
	}
	return true, nil
}

// StoreUser stores the user inside the storage
func (as Service) StoreUser(ctx context.Context, email, password,
	username string) (int, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	return as.repo.Store(ctx, UserModel{
		Email:    email,
		Password: string(encryptedPass),
		Username: username,
	})
}
