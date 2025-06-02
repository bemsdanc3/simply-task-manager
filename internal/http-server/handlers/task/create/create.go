package create

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "taskManager/internal/lib/api/response"
	"taskManager/internal/lib/logger/sl"
)

type Request struct {
	Content string `json:"content"`
}

type Response struct {
	resp.Response
	Content string `json:"content"`
}

type TaskCreator interface {
	CreateTask(contentOfTask string) (int64, error)
}

func New(log *slog.Logger, taskCreator TaskCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.create.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded: ", slog.Any("request", req))

		content := req.Content
		if content == "" {
			content = "Вы ничего не ввели."
		}

		id, err := taskCreator.CreateTask(content)
		if err != nil {
			log.Error("failed to create task", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to create task"))

			return
		}

		log.Info("task created", slog.Int64("id", id))

		responseOk(w, r, content)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, content string) {
	render.JSON(w, r, Response{
		Response: resp.Ok(),
		Content:  content,
	})
}
