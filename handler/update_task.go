package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Satoshi-Tb/go_todo_app/entity"
	"github.com/go-playground/validator/v10"
)

type UpdateTask struct {
	Service   UpdateTaskService
	Validator *validator.Validate
}

func (ut *UpdateTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		TaskID entity.TaskID     `json:"id" validate:"required"`
		Title  string            `json:"title" validate:"required"`
		Status entity.TaskStatus `json:"status" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	if err := ut.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	cnt, err := ut.Service.UpdateTask(ctx, b.TaskID, b.Title, b.Status)

	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		Cnt int `json:"count"`
	}{Cnt: cnt}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
