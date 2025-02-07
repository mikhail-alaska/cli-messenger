package messages

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/mikhail-alaska/cli-messenger/server/internal/http-server/middleware/auth"
	resp "github.com/mikhail-alaska/cli-messenger/server/internal/lib/api/response"
	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/logger/sl"
	"github.com/mikhail-alaska/cli-messenger/server/internal/storage"
)

type RequestMessages struct {
	Username    string `json:"username" validate:"required"`
}

type ResponseMessages struct {
	resp.Response
    Msgs []string
}

type MessageGetter interface {
	GetMessages(username1, username2 string) ([]storage.StorageMessages, error)
}

func GetMessages(log *slog.Logger, messageGetter MessageGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.messages.getall"

		username := auth.GetUsername(w, r)
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req RequestMessages

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
            fmt.Println("body", r.Body)
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

		id, err := messageGetter.GetMessages(username,req.Username)
		if err != nil {
			log.Error("failed to create message", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to create message"))
			return
		}

		log.Info("message sent", slog.Any("id", id))
        var msgs []string
        for _, j := range id{
            msgs= append(msgs, j.Message)
        }
		render.JSON(w, r, ResponseMessages{
			Response: resp.OK(),
            Msgs: msgs,
		})
	}
}
