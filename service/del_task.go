package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Satoshi-Tb/go_todo_app/auth"
	"github.com/Satoshi-Tb/go_todo_app/entity"
	"github.com/Satoshi-Tb/go_todo_app/store"
)

type DelTask struct {
	DB   store.Execer
	Repo TaskDeleter
}

func (d *DelTask) DelTask(ctx context.Context, taskID entity.TaskID) (int, error) {
	log.Print("DelTask service start")
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return -1, fmt.Errorf("user_id not found")
	}

	t := &entity.Task{
		ID:     taskID,
		UserID: id,
	}

	cnt, err := d.Repo.DelTask(ctx, d.DB, t)
	if err != nil {
		return -1, fmt.Errorf("failed to delete: %w", err)
	}
	return cnt, nil
}
