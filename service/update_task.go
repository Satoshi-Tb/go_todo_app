package service

import (
	"context"
	"fmt"

	"github.com/Satoshi-Tb/go_todo_app/auth"
	"github.com/Satoshi-Tb/go_todo_app/entity"
	"github.com/Satoshi-Tb/go_todo_app/store"
)

type UpdateTask struct {
	DB   store.Execer
	Repo TaskUpdater
}

func (ut *UpdateTask) UpdateTask(ctx context.Context, taskID entity.TaskID, title string, status entity.TaskStatus) (int, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return -1, fmt.Errorf("user_id not found")
	}

	t := &entity.Task{
		ID:     taskID,
		UserID: id,
		Title:  title,
		Status: status,
	}

	cnt, err := ut.Repo.UpdateTask(ctx, ut.DB, t)
	if err != nil {
		return -1, fmt.Errorf("failed to update: %w", err)
	}
	return cnt, nil
}
