// +build unit_tests all_tests

package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/ankur-anand/prod-todo/pkg/domain"
	"github.com/google/uuid"
)

type dummyTodorep struct {
	returnFunc  func() domain.TodoModel
	returnStore func(model domain.TodoModel) (uuid.UUID, error)
	errReturn   func(model domain.TodoModel) error
}

func (d dummyTodorep) Find(ctx context.Context, id uuid.UUID) (domain.TodoModel, error) {
	panic("implement me")
}

func (d dummyTodorep) FindAll(ctx context.Context, userID uuid.UUID) (domain.TodoIterator, error) {
	panic("implement me")
}

func (d dummyTodorep) Update(ctx context.Context, todo domain.TodoModel) error {
	return d.errReturn(todo)
}

func (d dummyTodorep) Store(ctx context.Context, todo domain.TodoModel) (uuid.UUID, error) {
	return d.returnStore(todo)
}

func (d dummyTodorep) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

func TestTodoService_UpdateOne(t *testing.T) {
	t.Parallel()
	d := &dummyTodorep{}

	dr, err := domain.NewTodoService(d)
	if err != nil {
		t.Fatal(err)
	}
	d.errReturn = func(model domain.TodoModel) error {
		t.Fatal("should not have been called")
		return nil
	}
	err = dr.UpdateOne(context.Background(), domain.NilTodoModel)
	if err == nil {
		t.Errorf("expected error to be not nil for NilTodoModel")
	}

	todoReceived := make(chan domain.TodoModel)
	d.errReturn = func(model domain.TodoModel) error {
		go func() {
			todoReceived <- model
		}()
		return nil
	}

	tdo := domain.TodoModel{
		Id:       uuid.New(),
		UserID:   uuid.New(),
		Title:    "finish this test",
		Content:  "Hi There",
		Finished: false,
	}
	err = dr.UpdateOne(context.Background(), tdo)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case todo := <-todoReceived:
		close(todoReceived)
		ok := todo.Title == "finish this test" && todo.Content == "Hi There" && todo.Id != uuid.Nil
		if !ok {
			t.Errorf("StoreTodo failed to received the expected Todo model")
		}
	case <-time.After(time.Second * 1):
		t.Errorf("storeTodo timedout")
	}
}

func TestTodoService_StoreNew(t *testing.T) {
	t.Parallel()
	d := &dummyTodorep{}

	dr, err := domain.NewTodoService(d)
	if err != nil {
		t.Fatal(err)
	}
	d.returnStore = func(model domain.TodoModel) (uuid.UUID, error) {
		t.Fatal("should not have been called")
		return uuid.Nil, nil
	}
	_, err = dr.StoreNew(context.Background(), domain.NilTodoModel)
	if err == nil {
		t.Errorf("expected error to be not nil for NilTodoModel")
	}

	todoReceived := make(chan domain.TodoModel)
	d.returnStore = func(model domain.TodoModel) (uuid.UUID, error) {
		go func() {
			todoReceived <- model
		}()
		return model.Id, nil
	}

	tdo := domain.TodoModel{
		UserID:   uuid.New(),
		Title:    "finish this test",
		Content:  "Hi There",
		Finished: false,
	}
	_, err = dr.StoreNew(context.Background(), tdo)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case todo := <-todoReceived:
		close(todoReceived)
		ok := todo.Title == "finish this test" && todo.Content == "Hi There" && todo.Id != uuid.Nil
		if !ok {
			t.Errorf("StoreTodo failed to received the expected Todo model")
		}
	case <-time.After(time.Second * 1):
		t.Errorf("storeTodo timedout")
	}
}
