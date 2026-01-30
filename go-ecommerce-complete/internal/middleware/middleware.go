package middleware

import (
	"encoding/json"
	"net/http"
)

var Sessions = make(map[string]map[string]interface{})

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "未登录", http.StatusUnauthorized)
			return
		}
		if _, ok := Sessions[cookie.Value]; !ok {
			http.Error(w, "会话已过期", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func RoleMiddleware(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, "未登录", http.StatusUnauthorized)
				return
			}
			session, ok := Sessions[cookie.Value]
			if !ok {
				http.Error(w, "会话已过期", http.StatusUnauthorized)
				return
			}
			userRole := session["role"].(string)
			allowed := false
			for _, role := range roles {
				if userRole == role {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, "权限不足", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}

func JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": message,
	})
}
