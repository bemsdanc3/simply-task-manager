package update

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

type Request struct {
	ContentToUpdate string `json:"content"`
}

type Response struct {
	UpdatedContent string
	resp.Response
}

type TaskUpdater interface {
	UpdateTask(contentToUpdate string, id int) error
}

func New(log *slog.Logger, TaskUpdater TaskUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.update.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded: ", slog.Any("request", req))

		contentToUpdate := req.ContentToUpdate
		if contentToUpdate == "" {
			log.Warn("update field has no value")
		}

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

		err = TaskUpdater.UpdateTask(contentToUpdate, taskID)
		if err != nil {
			log.Error("failed to update task", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to update task"))

			return
		}

		log.Info("task updated", slog.String("updated task: ", contentToUpdate))

		render.JSON(w, r, Response{
			UpdatedContent: contentToUpdate,
			Response:       resp.Ok(),
		})
	}
}
