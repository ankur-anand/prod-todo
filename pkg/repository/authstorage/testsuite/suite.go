package testsuite

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/ankur-anand/prod-app/pkg/repository/authstorage"

	"github.com/google/uuid"

	"github.com/ankur-anand/prod-app/pkg/auth"
)

// SuiteBase defines a re-usable set of authstorage-Repository related tests that can
// be executed against any type that implements authstorage.Repository.
type SuiteBase struct {
	r auth.Repository
}

// SetRepo configures the test-suite to run all tests against particular repo.
func (s *SuiteBase) SetRepo(r auth.Repository) {
	s.r = r
}

// TestFindAndStore verifies the find with ID logic.
// and Store Operation
func (s *SuiteBase) TestFindAndStore(t *testing.T) {
	_, err := s.r.Find(context.Background(), uuid.New())
	if err == nil {
		t.Errorf("expected error for random uuid user find")
	}

	if err != nil && errors.Is(err, authstorage.ErrUserNotFound) {
		t.Errorf("expected error type value [`no user found`] got `%s`", err.Error())
	}
	id := uuid.New()
	// store a new user
	user := auth.UserModel{
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
func (s *SuiteBase) TestFindByEmailAndStore(t *testing.T) {
	_, err := s.r.FindByEmail(context.Background(), "ank@an.com")
	if err == nil {
		t.Errorf("expected error for unknown email find")
	}
	if err != nil && errors.Is(err, authstorage.ErrUserNotFound) {
		t.Errorf("expected error type value [`no user found`] got `%s`", err.Error())
	}
	id := uuid.New()
	// store a new user
	user := auth.UserModel{
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
func (s *SuiteBase) TestDuplicateEmailStorePqSQL(t *testing.T) {
	id := uuid.New()
	email := id.String() + "example.com"
	// store a new user
	user := auth.UserModel{
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
