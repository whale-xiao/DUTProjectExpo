package service

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"showcase-backend/config"
	"showcase-backend/dto"
	"showcase-backend/pkg/errcode"
	pkgjwt "showcase-backend/pkg/jwt"
	"showcase-backend/repository"
)

// AuthService 认证业务：登录校验 + 签发 JWT。
type AuthService struct {
	db       *gorm.DB
	userRepo *repository.UserRepo
	cfg      *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:       db,
		userRepo: repository.NewUserRepo(db),
		cfg:      cfg,
	}
}

// LoginResult 登录结果。
type LoginResult struct {
	Token string      `json:"token"`
	User  dto.UserDTO `json:"user"`
}

// Login 校验用户名密码，成功则签发 24h Token 并更新最近登录时间。
func (s *AuthService) Login(username, password string) (*LoginResult, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return nil, NewError(errcode.InvalidParam, "用户名或密码不能为空")
	}
	u, err := s.userRepo.FindByUsername(username)
	if err != nil {
		// 用户不存在与密码错误返回同一提示，避免枚举账号
		return nil, NewError(errcode.Unauthorized, "用户名或密码不正确")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, NewError(errcode.Unauthorized, "用户名或密码不正确")
	}
	_ = s.userRepo.UpdateLastLogin(u.ID) // 更新失败不影响登录

	token, err := pkgjwt.Generate(s.cfg.JWTSecret, u.ID, u.Username, u.Role, 24*time.Hour)
	if err != nil {
		return nil, NewError(errcode.InternalError, "生成令牌失败")
	}
	return &LoginResult{
		Token: token,
		User: dto.UserDTO{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Role:     u.Role,
		},
	}, nil
}
