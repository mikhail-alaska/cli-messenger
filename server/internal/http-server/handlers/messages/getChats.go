package messages

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/mikhail-alaska/cli-messenger/server/internal/http-server/middleware/auth"
	resp "github.com/mikhail-alaska/cli-messenger/server/internal/lib/api/response"
	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/logger/sl"
)


type ResponseAllChats struct {
	resp.Response
	Chats []string
}

type ChatGetter interface {
	GetChatsByUsername(username string) ([]string, error)
}

func GetChats(log *slog.Logger, chatGetter ChatGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.messages.chatget"

        username := auth.GetUsername(w, r)
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)


		out, err := chatGetter.GetChatsByUsername(username)
		if err != nil {
			log.Error("failed to get chats", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get chats"))
			return
		}

		log.Info("chats gotten", slog.Any("chats len", len(out)))
		render.JSON(w, r, ResponseAllChats{
			Response: resp.OK(),
            Chats: out,
		})
	}
}
