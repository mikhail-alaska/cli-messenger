package users

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/mikhail-alaska/cli-messenger/server/internal/lib/api/response"
	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/logger/sl"
)

type AllUserGetter interface {
	GetAllUsers() ([]string, error)
}

type UserAllResponse struct{
    Users []string
	resp.Response
}

func AllUsers(log *slog.Logger, allUserGetter AllUserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.getall."

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)


		users, err := allUserGetter.GetAllUsers()
		if err != nil {
			log.Error("failed to het all users", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to het all users"))
			return
		}


		log.Info("users gotten", slog.Any("users num", len(users)))
        render.JSON(w,r, UserAllResponse{
            Response: resp.OK(),
            Users: users,
        })
	}
}
