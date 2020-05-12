package domain

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
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

// AuthService provides the use cases implementation to work
// with the entities of the underlying domain during
// SignIn and SignUP
type AuthService struct {
	repo UserRepository
}

// NewService returns a new AuthService initialized with
// a concrete repo implementation
func NewService(repo UserRepository) AuthService {
	return AuthService{
		repo: repo,
	}
}

// IsValidEmail checks if an email is valid or not
func (as AuthService) IsValidEmail(email string) bool {
	email = normalize(email)
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
func (as AuthService) IsValidPassword(password string) bool {
	// min 8 character only and less than 254
	if len(strings.TrimSpace(password)) > 254 || len(strings.TrimSpace(password)) < 8 {
		return false
	}
	return true
}

// IsCredentialValid checks if the Credential is ok
func (as AuthService) IsCredentialValid(ctx context.Context, email string,
	password string) (bool, error) {
	email = normalize(email)
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
func (as AuthService) IsDuplicateRegistration(ctx context.Context, email string) (bool,
	error) {
	email = normalize(email)
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
func (as AuthService) StoreUser(ctx context.Context, model UserModel) (uuid.UUID, error) {
	email := normalize(model.Email)
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(model.Password),
		bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}
	return as.repo.Store(ctx, UserModel{
		ID:        uuid.New(),
		Email:     email,
		Password:  string(encryptedPass),
		FirstName: model.FirstName,
		LastName:  model.LastName,
	})
}
