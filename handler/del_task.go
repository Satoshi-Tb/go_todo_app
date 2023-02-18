package handler

import (
	"encoding/json"
	"log"
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
	log.Print("DelTask handler start")
	var b struct {
		TaskID entity.TaskID `json:"id" validate:"required"`
	}
	log.Print("DelTask decode json")
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	log.Print("DelTask validate json")
	log.Printf("json value %+v", b)
	if err := dt.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	log.Print("DelTask DelTask service")
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
