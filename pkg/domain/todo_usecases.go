package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// TodoService provides the use cases implementation to work
// with the entities of the underlying domain during
// Todo create and update
type TodoService struct {
	repo TodoRepository
}

// NewTodoService returns an initialized TodoService
func NewTodoService(rep TodoRepository) (TodoService, error) {
	if rep == nil {
		return TodoService{}, fmt.Errorf("required TodoRepository cannot be nil")
	}
	return TodoService{repo: rep}, nil
}

// FindByID returns a TodoModel associated with the given ids
func (ts TodoService) FindByID(ctx context.Context, id uuid.UUID) (TodoModel, error) {
	return ts.repo.Find(ctx, id)
}

// GetUserTodoList returns an TodoIterator which contains all todo associated with
// the given user id
func (ts TodoService) GetUserTodoList(ctx context.Context, userID uuid.UUID) (TodoIterator, error) {
	return ts.repo.FindAll(ctx, userID)
}

// UpdateOne updates the provided todo model in the underlying storage
func (ts TodoService) UpdateOne(ctx context.Context, todo TodoModel) error {
	if todo.Id == uuid.Nil {
		return fmt.Errorf("cannot update todo with empty id")
	}

	if len(todo.Title) == 0 || len(todo.Content) == 0 {
		return fmt.Errorf("todo title and content should not be empty")
	}
	return ts.repo.Update(ctx, todo)
}

// DeleteOne deletes the provided todo model in the underlying storage
func (ts TodoService) DeleteOne(ctx context.Context, todo TodoModel) error {
	if todo.Id == uuid.Nil {
		return fmt.Errorf("cannot delete todo with empty id")
	}
	return ts.repo.Delete(ctx, todo.Id)
}

// StoreNew saves the provided todo model in the underlying storage
func (ts TodoService) StoreNew(ctx context.Context, todo TodoModel) (uuid.UUID, error) {
	if len(todo.Title) == 0 || len(todo.Content) == 0 {
		return uuid.Nil, fmt.Errorf("todo title and content should not be empty")
	}
	tdo := TodoModel{
		Id:       uuid.New(),
		Title:    todo.Title,
		Content:  todo.Content,
		Finished: todo.Finished,
	}
	return ts.repo.Store(ctx, tdo)
}
