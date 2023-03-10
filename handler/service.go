package handler

import (
	"context"

	"github.com/Satoshi-Tb/go_todo_app/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService AddTaskService RegisterUserService LoginService DelTaskService UpdateTaskService
type ListTasksService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}

type DelTaskService interface {
	DelTask(ctx context.Context, taskID entity.TaskID) (int, error)
}

type UpdateTaskService interface {
	UpdateTask(ctx context.Context, taskID entity.TaskID, title string, status entity.TaskStatus) (int, error)
}

type RegisterUserService interface {
	RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error)
}

type LoginService interface {
	Login(ctx context.Context, name, pw string) (string, error)
}
