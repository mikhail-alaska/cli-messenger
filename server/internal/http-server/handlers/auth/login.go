package login

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v4"
	resp "github.com/mikhail-alaska/cli-messenger/server/internal/lib/api/response"
	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/logger/sl"
)

// jwtSecret — секрет для подписи токена.
// В реальном приложении его лучше хранить в переменных окружения или в безопасном хранилище.
var jwtSecret = []byte("your_secret_key")

// LoginRequest описывает данные, которые ожидаются в запросе.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse описывает структуру ответа с токеном.
type LoginResponse struct {
	resp.Response
	Token string `json:"token"`
}

type PaswordGetter interface {
	GetPassword(username string) (string, error)
}

// LoginHandler проверяет логин и пароль и возвращает JWT-токен, если авторизация успешна.
func LoginHandler(log *slog.Logger, paswordGetter PaswordGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.login"
		var req LoginRequest

		log = log.With(
			slog.String("op", op),
		)

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

        pass, err := paswordGetter.GetPassword(req.Username)

		// Простейшая проверка учетных данных.
		// Обычно здесь происходит проверка в базе данных.
		if req.Password != pass{

			log.Error("Invalid credentials", sl.Err(err))
			render.JSON(w, r, resp.Error("Invalid credentials"))
			return

		}

		// Создаем токен с необходимыми claims.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(72 * time.Hour).Unix(), // Токен действителен 72 часа.
		})

		// Подписываем токен секретным ключом.
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			log.Error("failed to sign", sl.Err(err))
			render.JSON(w, r, resp.Error("internal server error lol"))
			return
		}

		// Отдаем токен в JSON-формате.
		response := LoginResponse{
            Response: resp.OK(),
			Token: tokenString,
		}

		w.Header().Set("Content-Type", "application/json")

        render.JSON(w,r, response)
	}
}
