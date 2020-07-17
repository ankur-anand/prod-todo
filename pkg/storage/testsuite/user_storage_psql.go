package testsuite

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/ankur-anand/prod-todo/pkg/storage/serror"

	"github.com/ankur-anand/prod-todo/pkg"
	"github.com/google/uuid"
)

// UserSuiteBase defines a re-usable set of user storage related tests that can
// be executed against any type that implements pkg.UserStorage.
type UserSuiteBase struct {
	r pkg.UserStorage
}

// SetRepo configures the test-suite to run all tests against particular repo.
func (s *UserSuiteBase) SetRepo(r pkg.UserStorage) {
	s.r = r
}

// TestFindAndStore verifies the find with ID logic.
// and Store Operation
func (s *UserSuiteBase) TestFindAndStore(t *testing.T) {
	_, err := s.r.Find(context.Background(), uuid.New())
	if err == nil {
		t.Errorf("expected error for random uuid user find")
	}

	if err != nil && errors.Is(err, serror.ErrUserNotFound) {
		t.Errorf("expected error type value [`no user found`] got `%s`", err.Error())
	}
	id := uuid.New()
	// store a new user
	user := pkg.UserModel{
		ID:        id,
		Email:     "ankuranand@example.com",
		Password:  "somegibrish&^5$(075",
		FirstName: "Ankur",
		LastName:  "Anand",
		Username:  "ankur-anand",
	}
	rID, err := s.r.Store(context.Background(), user)
	if err != nil {
		t.Errorf("exected a nil error for store got %v", err)
	}

	if rID != id {
		t.Errorf("expected uuid [%v] got [%v]", id, rID)
	}

	findUser, err := s.r.Find(context.Background(), id)
	if err != nil {
		t.Errorf("exected a nil error for store got %v", err)
	}

	if !reflect.DeepEqual(findUser, user) {
		t.Errorf("expected find user [%+v] to have a equal to inserted user [%+v]", findUser, user)
	}
}

// TestFindByEmailAndStore verifies the find with Email logic.
// and Store Operation
func (s *UserSuiteBase) TestFindByEmailAndStore(t *testing.T) {
	_, err := s.r.FindByEmail(context.Background(), "ank@an.com")
	if err == nil {
		t.Errorf("expected error for unknown email find")
	}
	if err != nil && errors.Is(err, serror.ErrUserNotFound) {
		t.Errorf("expected error type value [`no user found`] got `%s`", err.Error())
	}
	id := uuid.New()
	// store a new user
	user := pkg.UserModel{
		ID:        id,
		Email:     "ankuranand1@example.com",
		Password:  "somegibrish&^5$(075",
		FirstName: "Ankur",
		LastName:  "Anand",
		Username:  "ankur-anand",
	}
	rID, err := s.r.Store(context.Background(), user)
	if err != nil {
		t.Errorf("exected a nil error for store got %v", err)
	}

	if rID != id {
		t.Errorf("expected uuid [%v] got [%v]", id, rID)
	}

	findUser, err := s.r.FindByEmail(context.Background(), "ankuranand1@example.com")
	if err != nil {
		t.Errorf("exected a nil error for store got %v", err)
	}

	if !reflect.DeepEqual(findUser, user) {
		t.Errorf("expected find user [%+v] to have a equal to inserted user [%+v]", findUser, user)
	}
}

// TestDuplicateEmailStorePqSQL verifies that storing duplicates results in
// conflict, especially with SQL type query, with current implementation
// with the postgreSQL
func (s *UserSuiteBase) TestDuplicateEmailStorePqSQL(t *testing.T) {
	id := uuid.New()
	email := id.String() + "example.com"
	// store a new user
	user := pkg.UserModel{
		ID:        id,
		Email:     email,
		Password:  "somegibrish&^5$(075",
		FirstName: "Ankur",
		LastName:  "Anand",
		Username:  "ankur-anand",
	}
	rID, err := s.r.Store(context.Background(), user)
	if err != nil {
		t.Errorf("exected a nil error for store got %v", err)
	}

	if rID != id {
		t.Errorf("expected uuid [%v] got [%v]", id, rID)
	}

	_, err = s.r.Store(context.Background(), user)
	if err == nil {
		t.Errorf("exected a non nil error for duplicate store got")
	}
	if err != nil {
		if !strings.Contains(err.Error(), "ERROR: duplicate key value violates unique constraint") {
			t.Errorf("expected error of type duplicate key value violation got %s", err.Error())
		}
	}
}

// TestUpdateUserPqSQL verifies that storing duplicates results in
// conflict, especially with SQL type query, with current implementation
// with the postgreSQL
func (s *UserSuiteBase) TestUpdateUserPqSQL(t *testing.T) {
	id := uuid.New()
	email := id.String() + "example.com"
	// store a new user
	user := pkg.UserModel{
		ID:        id,
		Email:     email,
		Password:  "somegibrish&^5$(075",
		FirstName: "Ankur",
		LastName:  "Anand",
		Username:  "ankur-anand",
	}
	rID, err := s.r.Store(context.Background(), user)
	if err != nil {
		t.Errorf("exected a nil error for store got %v", err)
	}

	if rID != id {
		t.Errorf("expected uuid [%v] got [%v]", id, rID)
	}
	// update email
	user.Email = "updated@email.com"
	err = s.r.Update(context.Background(), user)
	if err != nil {
		t.Errorf("exected a nil error for update user got %v", err)
	}
	userUpdated, err := s.r.FindByEmail(context.Background(), "updated@email.com")
	if err != nil {
		t.Errorf("exected a nil error for update user got %v", err)
	}

	if userUpdated.Email != "updated@email.com" {
		t.Errorf("expected a email %s got %s", "updated@email.com", userUpdated.Email)
	}

	if !reflect.DeepEqual(userUpdated, user) {
		t.Errorf("expected find user [%+v] to have a equal to inserted user [%+v]", userUpdated, user)
	}
}
