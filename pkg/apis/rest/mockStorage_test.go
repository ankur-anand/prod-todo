// +build unit_tests all_tests

package rest

import (
	"context"

	"github.com/ankur-anand/prod-todo/pkg/domain"
	"github.com/google/uuid"
)

type _mockUserRepoStorage struct {
	returnFunc  func() domain.UserModel
	returnStore func(domain.UserModel) (uuid.UUID, error)
}

func (m *_mockUserRepoStorage) Find(ctx context.Context, id uuid.UUID) (domain.UserModel, error) {
	return m.returnFunc(), nil
}

func (m *_mockUserRepoStorage) FindByEmail(ctx context.Context, email string) (domain.UserModel, error) {
	return m.returnFunc(), nil
}

func (m *_mockUserRepoStorage) FindAll(ctx context.Context) (domain.UserIterator, error) {
	panic("implement me")
}

func (m *_mockUserRepoStorage) Update(ctx context.Context, user domain.UserModel) error {
	panic("implement me")
}

func (m *_mockUserRepoStorage) Store(ctx context.Context, user domain.UserModel) (uuid.UUID, error) {
	return m.returnStore(user)
}
