package users

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/mikhail-alaska/cli-messenger/server/internal/lib/api/response"
	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/logger/sl"
)

type OpenKeyGetter interface {
	GetOpenKey(username string) (int, error)
}

type OpenKeyRequest struct {
	Username string `json:"username" validate:"required"`
}

type OpenKeyResponse struct {
	OpenKey int
	resp.Response
}

func OpenKeyByUserName(log *slog.Logger, openKeyGetter OpenKeyGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.openKeyGetter"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req OpenKeyRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		opke, err := openKeyGetter.GetOpenKey(req.Username)
		if err != nil {
			log.Error("failed to get open key", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get open key"))
			return
		}

		log.Info("open key gotten", slog.Any("openkey", opke))
		render.JSON(w, r, OpenKeyResponse{
			Response: resp.OK(),
			OpenKey:  opke,
		})
	}
}
