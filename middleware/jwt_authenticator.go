package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenFromHeader := r.Header.Get("Authorization")
		if tokenFromHeader == "" {
			http.Error(w, fmt.Errorf("authorization Token Header is missing").Error(), http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(tokenFromHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, fmt.Errorf("invalid Bearer Token format").Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing methods : %v", t.Header["alg"])
			}

			// Load env
			err := godotenv.Load()
			if err != nil {
				return nil, err
			}

			return []byte(os.Getenv("REFRESH_JWT_SECRET")), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		// Success checking
		next.ServeHTTP(w, r)
	}
}
