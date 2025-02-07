package users

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/mikhail-alaska/cli-messenger/server/internal/lib/api/response"
	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/logger/sl"
	"github.com/mikhail-alaska/cli-messenger/server/internal/storage"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	OpenKey  int    `json:"openkey" validate:"required"`
}

type Response struct {
	resp.Response
}

type UserCreater interface {
	CreateUser(userName string, openKey int) (int64, error)
}


func NewUser(log *slog.Logger, userCreater UserCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.create"

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

		id, err := userCreater.CreateUser(req.Username,req.OpenKey)
		if errors.Is(err, storage.ErrUserNameExists) {
            log.Info("username already exists", slog.String("username", req.Username))
			render.JSON(w, r, resp.Error("username already exists"))
			return
		}
		if err != nil {
			log.Error("failed to save user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save user"))
			return
		}


		log.Info("user added", slog.Any("id", id))
        render.JSON(w,r, Response{
            Response: resp.OK(),
        })
	}
}
