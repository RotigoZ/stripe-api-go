package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"

func AuthMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusBadRequest)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer"{
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}
			tokenString := headerParts[1]

			secretKey := []byte(os.Getenv("JWT_SECRET"))
			token, erro := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secretKey, nil
			})

			if erro != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userID := claims["sub"]
				userRole := claims["role"]

				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				ctx = context.WithValue(ctx, UserRoleKey, userRole)
            
        
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            http.Error(w, "Could not parse token claims", http.StatusInternalServerError)
        }
			
			
	})
}

func AdminMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			role, ok := r.Context().Value(UserRoleKey).(string)

			if !ok || role != "admin" && role != "superadmin" {
				http.Error(w, "Forbidden: You do not have permission to permorf this action", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
	})
}

func SuperAdminMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			role, ok := r.Context().Value(UserRoleKey).(string)

			if !ok || role != "superadmin" {
				http.Error(w, "Forbidden: You do not have permission to permorf this action", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
	})
}