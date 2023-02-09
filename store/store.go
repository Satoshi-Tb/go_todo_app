package store

import (
	"errors"

	"github.com/Satoshi-Tb/go_todo_app/entity"
)

// 仮実装。mapを使う
type TaskStore struct {
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

var Tasks = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}}
var ErrNotFound = errors.New("not found")

func (ts *TaskStore) Add(t *entity.Task) (int, error) {
	ts.LastID++ // 最小値1からの連番
	t.ID = ts.LastID
	ts.Tasks[t.ID] = t
	return int(t.ID), nil
}

func (ts *TaskStore) All() entity.Tasks {
	tasks := make([]*entity.Task, len(ts.Tasks))
	for k, v := range ts.Tasks {
		tasks[k-1] = v
	}
	return tasks
}
