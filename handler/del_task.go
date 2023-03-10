package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Satoshi-Tb/go_todo_app/entity"
	"github.com/go-playground/validator/v10"
)

type DelTask struct {
	Service   DelTaskService
	Validator *validator.Validate
}

func (dt *DelTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		TaskID entity.TaskID `json:"id" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := dt.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	cnt, err := dt.Service.DelTask(ctx, b.TaskID)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		ID int `json:"id"`
	}{ID: cnt}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
