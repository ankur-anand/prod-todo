// +build unit_tests all_tests

package pkg

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type dummyRepo struct {
	returnFunc  func() UserModel
	returnStore func(UserModel) (uuid.UUID, error)
}

func (d dummyRepo) Find(ctx context.Context, id uuid.UUID) (UserModel, error) {
	panic("implement me")
}

func (d dummyRepo) FindByEmail(ctx context.Context,
	email string) (UserModel, error) {
	return d.returnFunc(), nil
}

func (d dummyRepo) Update(ctx context.Context, user UserModel) error {
	panic("implement me")
}

func (d dummyRepo) Store(ctx context.Context, user UserModel) (uuid.UUID, error) {
	return d.returnStore(user)
}

func TestService_IsValidEmail(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "invalid email @missing",
			email: "ankur.com",
			want:  false,
		},
		{
			name:  "invalid email model",
			email: "ankur@.com",
			want:  false,
		},
		{
			name:  "valid email",
			email: "ankur@example.com",
			want:  true,
		},
		{
			name:  "invalid user name is more than 64 characters",
			email: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@gm.com",
			want:  false,
		},
	}
	as := NewRegAndAuthService(dummyRepo{})
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if as.IsValidEmail(tc.email) != tc.want {
				t.Errorf("email validation failed for %s, want %v, got %v", tc.email, tc.want, as.IsValidEmail(tc.email))
			}
		})

	}
}

func TestService_IsValidPassword(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "invalid len password small",
			password: "ankur",
			want:     false,
		},
		{
			name:     "invalid length password too large",
			password: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want:     false,
		},
		{
			name:     "valid password",
			password: "ankur@example.com",
			want:     true,
		},
	}
	as := NewRegAndAuthService(dummyRepo{})
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if as.IsValidPassword(tc.password) != tc.want {
				t.Errorf("email validation failed for %s, want %v, got %v", tc.password, tc.want, as.IsValidPassword(tc.password))
			}
		})

	}
}

func TestService_IsDuplicateRegistration(t *testing.T) {
	t.Parallel()
	dummyR := dummyRepo{}
	dummyR.returnFunc = func() UserModel {
		return UserModel{
			ID:       uuid.New(),
			Email:    "ankuranand@example.com",
			Password: "garbage",
			Username: "ankuranand",
		}
	}
	as := NewRegAndAuthService(dummyR)
	ok, _ := as.IsDuplicateRegistration(context.Background(), "anKuranand@example.com")
	if !ok {
		t.Errorf("duplicate Registration validation failed for %s", "ankuranand@example.com")
	}

	ok, _ = as.IsDuplicateRegistration(context.Background(), "anKur@example.com")
	if ok {
		t.Errorf("duplicate Registration validation failed for %s", "ankur@example.com")
	}
}

func TestService_IsCredentialValid(t *testing.T) {
	t.Parallel()
	password := "ankuranand"
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	dummyR := dummyRepo{}
	dummyR.returnFunc = func() UserModel {
		return UserModel{
			ID:       uuid.New(),
			Email:    "ankuranand@example.com",
			Password: string(encryptedPass),
			Username: "ankuranand",
		}
	}
	as := NewRegAndAuthService(dummyR)
	ok, user, _ := as.IsCredentialValid(context.Background(), "ankuranand@example.com", password)
	if !ok && user.Email != "ankuranand@example.com" {
		t.Errorf("credentail validation failed")
	}

	ok, user, _ = as.IsCredentialValid(context.Background(), "ankuranand@example.com", "garbage")
	if ok && user != NilUserModel {
		t.Errorf("credentail validation failed")
	}
}

func TestService_StoreUser(t *testing.T) {
	t.Parallel()
	var err error
	password := "ankuranand"

	dummyR := dummyRepo{}
	userReceived := make(chan UserModel)

	dummyR.returnStore = func(model UserModel) (uuid.UUID, error) {
		go func() {
			userReceived <- model
		}()
		return model.ID, nil
	}

	as := NewRegAndAuthService(dummyR)
	usr := UserModel{
		Email:     "AnkurananD@example.com", // email should be normalized
		Password:  password,
		FirstName: "Ankur",
		LastName:  "Anand",
	}
	_, err = as.StoreUser(context.Background(), usr)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case user := <-userReceived:
		close(userReceived)
		ok := user.Email == "ankuranand@example.com" && user.FirstName == "Ankur" && password != user.Password
		if !ok {
			t.Errorf("StoreUser failed to received the expected user model")
		}
	case <-time.After(time.Second * 1):
		t.Errorf("storeUser timedout")
	}
}
