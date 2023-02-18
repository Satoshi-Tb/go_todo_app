package store

import (
	"context"

	"github.com/Satoshi-Tb/go_todo_app/entity"
)

func (r *Repository) ListTasks(ctx context.Context, db Queryer, id entity.UserID) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT id, user_id, title, status, created, modified FROM task WHERE user_id = ?;`
	if err := db.SelectContext(ctx, &tasks, sql, id); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(ctx context.Context, db Execer, t *entity.Task) error {
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()

	sql := `INSERT INTO task(user_id, title, status, created, modified) VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(
		ctx, sql, t.UserID, t.Title, t.Status, t.Created, t.Modified,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = entity.TaskID(id)
	return nil
}

func (r *Repository) DelTask(ctx context.Context, db Execer, t *entity.Task) (int, error) {
	sql := `DELETE FROM task WHERE id = ? AND user_id = ?`
	result, err := db.ExecContext(ctx, sql, t.ID, t.UserID)
	if err != nil {
		return -1, err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return int(cnt), nil
}

func (r *Repository) UpdateTask(ctx context.Context, db Execer, t *entity.Task) (int, error) {
	t.Modified = r.Clocker.Now()
	sql := `UPDATE task SET title = ?, status = ?, modified = ? WHERE id = ? AND user_id = ?`
	result, err := db.ExecContext(ctx, sql, t.Title, t.Status, t.Modified, t.ID, t.UserID)
	if err != nil {
		return -1, err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return int(cnt), nil
}
