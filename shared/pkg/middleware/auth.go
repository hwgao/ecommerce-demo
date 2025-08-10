package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
    "ecommerce/shared/pkg/response"
)

func JWTAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            response.Error(w, http.StatusUnauthorized, "Authorization header required")
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("your-secret-key"), nil // Use env variable in production
        })

        if err != nil || !token.Valid {
            response.Error(w, http.StatusUnauthorized, "Invalid token")
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            response.Error(w, http.StatusUnauthorized, "Invalid token claims")
            return
        }

        userID, err := uuid.Parse(claims["user_id"].(string))
        if err != nil {
            response.Error(w, http.StatusUnauthorized, "Invalid user ID")
            return
        }

        ctx := context.WithValue(r.Context(), "user_id", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
