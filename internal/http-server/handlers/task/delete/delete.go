package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	resp "taskManager/internal/lib/api/response"
	"taskManager/internal/lib/logger/sl"
)

type Response struct {
	DeletedID int
	resp.Response
}

type TaskDeleter interface {
	DeleteTask(id int) error
}

func New(log *slog.Logger, taskDeleter TaskDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requst_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		taskID, err := strconv.Atoi(id)
		if err != nil {
			log.Error("wrong id type", sl.Err(err))

			render.JSON(w, r, resp.Error("wrong id type"))

			return
		}
		if taskID == 0 {
			log.Info("id is zero")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err = taskDeleter.DeleteTask(taskID)
		if err != nil {
			log.Error("can not delete task", sl.Err(err))

			render.JSON(w, r, resp.Error("can not delete task"))

			return
		}

		log.Info("task deleted", slog.Int("deleted task ID: ", taskID))

		render.JSON(w, r, Response{
			DeletedID: taskID,
			Response:  resp.Ok(),
		})
	}
}
