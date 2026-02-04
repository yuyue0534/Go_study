package service

import (
	"crm-erp-system/config"
	"crm-erp-system/database"
	"crm-erp-system/model"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct{}

func (s *UserService) Register(req *model.RegisterRequest) error {
	// 检查用户名是否存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 插入用户
	_, err = database.DB.Exec(
		"INSERT INTO users (username, password, email, phone) VALUES (?, ?, ?, ?)",
		req.Username, string(hashedPassword), req.Email, req.Phone,
	)
	return err
}

func (s *UserService) Login(req *model.LoginRequest) (string, error) {
	var user model.User
	err := database.DB.QueryRow(
		"SELECT id, username, password FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 生成JWT
	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) generateToken(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func (s *UserService) GetUserInfo(userID int64) (*model.User, error) {
	var user model.User
	err := database.DB.QueryRow(
		"SELECT id, username, email, phone, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.CreatedAt)

	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}
