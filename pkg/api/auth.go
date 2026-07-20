package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "todo-app-secret-key"

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func validateToken(tokenString, password string) bool {
	if tokenString == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	hashValue, ok := claims["hash"]
	if !ok {
		return false
	}

	hash, ok := hashValue.(string)
	if !ok {
		return false
	}

	return hash == passwordHash(password)
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Если пароль не задан — аутентификация отключена.
		if todoPassword == "" {
			next(w, r)
			return
		}

		var token string

		cookie, err := r.Cookie("token")
		if err == nil {
			token = cookie.Value
		}

		if !validateToken(token, todoPassword) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

type signinRequest struct {
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var request signinRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if request.Password != todoPassword {
		writeError(w, http.StatusUnauthorized, "неверный пароль")
		return
	}

	claims := jwt.MapClaims{
		"hash": passwordHash(todoPassword),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create token")
		return
	}

	writeJSON(w, signinResponse{
		Token: signedToken,
	})
}
