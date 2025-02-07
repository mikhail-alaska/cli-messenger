package messages

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/mikhail-alaska/cli-messenger/server/internal/http-server/middleware/auth"
	resp "github.com/mikhail-alaska/cli-messenger/server/internal/lib/api/response"
	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/logger/sl"
)

type Request struct {
	Username    string `json:"username" validate:"required"`
	MessageFor1 string `json:"messagefor1" validate:"required"`
	MessageFor2 string `json:"messagefor2" validate:"required"`
}

type Response struct {
	resp.Response
}

type MessageCreater interface {
	CreateMessage(userName1, userName2, msg1, msg2 string) (int64, error)
}

func NewMessage(log *slog.Logger, messageCreater MessageCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.messages.create"

        username := auth.GetUsername(w, r)
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

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		id, err := messageCreater.CreateMessage(username, req.Username, req.MessageFor1, req.MessageFor2)
		if err != nil {
			log.Error("failed to create message", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to create message"))
			return
		}

		log.Info("message sent", slog.Any("id", id))
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
