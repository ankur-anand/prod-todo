// +build unit_tests all_tests

package resthandler

import (
	"context"

	"github.com/ankur-anand/prod-todo/pkg"

	"github.com/google/uuid"
)

type _mockUserRepoStorage struct {
	returnFunc  func() pkg.UserModel
	returnStore func(pkg.UserModel) (uuid.UUID, error)
}

func (m *_mockUserRepoStorage) Find(ctx context.Context, id uuid.UUID) (pkg.UserModel, error) {
	return m.returnFunc(), nil
}

func (m *_mockUserRepoStorage) FindByEmail(ctx context.Context, email string) (pkg.UserModel, error) {
	return m.returnFunc(), nil
}

func (m *_mockUserRepoStorage) FindAll(ctx context.Context) (pkg.UserIterator, error) {
	panic("implement me")
}

func (m *_mockUserRepoStorage) Update(ctx context.Context, user pkg.UserModel) error {
	panic("implement me")
}

func (m *_mockUserRepoStorage) Store(ctx context.Context, user pkg.UserModel) (uuid.UUID, error) {
	return m.returnStore(user)
}
