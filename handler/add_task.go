package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Satoshi-Tb/go_todo_app/entity"
	"github.com/Satoshi-Tb/go_todo_app/store"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type AddTask struct {
	DB        *sqlx.DB
	Repo      *store.Repository
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	if err := at.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		Title:  b.Title,
		Status: entity.TaskStatusTodo,
	}

	if err := at.Repo.AddTask(ctx, at.DB, t); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		ID int `json:"id"`
	}{ID: int(t.ID)}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
