package middleware

import (
	"context"
	"ecommerce-platform/internal/models"
	"ecommerce-platform/internal/utils"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Session 会话信息
type Session struct {
	UserID    int64
	Username  string
	Role      string
	SellerID  int64
	ExpiresAt time.Time
}

// SessionStore 会话存储
var SessionStore = struct {
	sync.RWMutex
	sessions map[string]*Session
}{
	sessions: make(map[string]*Session),
}

// ContextKey 上下文键
type ContextKey string

const (
	UserKey    ContextKey = "user"
	SessionKey ContextKey = "session"
)

// CreateSession 创建会话
func CreateSession(user *models.User) string {
	token := utils.GenerateToken()
	
	session := &Session{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// 如果是商家，获取商家ID
	if user.Role == "seller" {
		seller, err := models.GetSellerByUserID(user.ID)
		if err == nil {
			session.SellerID = seller.ID
		}
	}

	SessionStore.Lock()
	SessionStore.sessions[token] = session
	SessionStore.Unlock()

	return token
}

// GetSession 获取会话
func GetSession(token string) *Session {
	SessionStore.RLock()
	session, ok := SessionStore.sessions[token]
	SessionStore.RUnlock()

	if !ok || session.ExpiresAt.Before(time.Now()) {
		return nil
	}
	return session
}

// DeleteSession 删除会话
func DeleteSession(token string) {
	SessionStore.Lock()
	delete(SessionStore.sessions, token)
	SessionStore.Unlock()
}

// CORS 跨域中间件
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logger 日志中间件
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		_ = duration // 可以在此处添加日志记录
	})
}

// Auth 认证中间件
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从Header获取Token
		authHeader := r.Header.Get("Authorization")
		var token string
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 也尝试从Cookie获取
		if token == "" {
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			utils.Unauthorized(w, "请先登录")
			return
		}

		session := GetSession(token)
		if session == nil {
			utils.Unauthorized(w, "会话已过期，请重新登录")
			return
		}

		// 将会话信息添加到上下文
		ctx := context.WithValue(r.Context(), SessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// OptionalAuth 可选认证中间件（不强制登录）
func OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从Header获取Token
		authHeader := r.Header.Get("Authorization")
		var token string
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 也尝试从Cookie获取
		if token == "" {
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token != "" {
			session := GetSession(token)
			if session != nil {
				ctx := context.WithValue(r.Context(), SessionKey, session)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	}
}

// RequireRole 角色验证中间件
func RequireRole(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return Auth(func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value(SessionKey).(*Session)
			
			for _, role := range roles {
				if session.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			
			utils.Forbidden(w, "没有权限访问")
		})
	}
}

// RequireAdmin 管理员权限
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return RequireRole("admin")(next)
}

// RequireSeller 商家权限
func RequireSeller(next http.HandlerFunc) http.HandlerFunc {
	return RequireRole("seller", "admin")(next)
}

// RequireCustomer 顾客权限
func RequireCustomer(next http.HandlerFunc) http.HandlerFunc {
	return RequireRole("customer", "seller", "admin")(next)
}

// GetCurrentSession 从上下文获取当前会话
func GetCurrentSession(r *http.Request) *Session {
	session, ok := r.Context().Value(SessionKey).(*Session)
	if !ok {
		return nil
	}
	return session
}
