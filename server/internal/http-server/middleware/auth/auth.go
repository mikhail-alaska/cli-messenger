package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your_secret_key") // Задай свой секретный ключ

// JWTAuth — простейшее middleware для проверки JWT токена.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем заголовок авторизации
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Ожидаем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Парсим и валидируем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что метод подписи — HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Извлекаем username из claims
			_, ok := claims["username"].(string)
			if !ok {
				http.Error(w, "Username not found in token", http.StatusUnauthorized)
				return
			}
			// Здесь можно добавить username в контекст запроса или продолжить обработку
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}


func GetUsername(w http.ResponseWriter, r *http.Request) string {
    var username string
		authHeader := r.Header.Get("Authorization")

		parts := strings.Split(authHeader, " ")

		tokenString := parts[1]

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			username, _ = claims["username"].(string)
		} 	
    return username
}
