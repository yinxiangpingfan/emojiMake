package services

import (
	"emoji-maker-backend/models"
	"emoji-maker-backend/repositories"
	"errors"

	"github.com/gofiber/fiber/v2"
)

// APIResponse 统一API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// UserService 用户服务接口
type UserService interface {
	Register(phone, password string) (*APIResponse, error)
	Login(phone, password string) (*APIResponse, error)
	ChangePassword(userID int64, newPassword string) (*APIResponse, error)
}

// userServiceImpl 用户服务实现
type userServiceImpl struct {
	userRepo   repositories.UserRepository
	jwtService JWTService
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repositories.UserRepository, jwtService JWTService) UserService {
	return &userServiceImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register 用户注册
func (s *userServiceImpl) Register(phone, password string) (*APIResponse, error) {
	// 检查用户是否已存在
	existingUser, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this phone number already exists")
	}

	// 创建新用户并加密密码
	user := &models.User{
		Phone: phone,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// 保存用户
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &APIResponse{
		Code:    0,
		Message: "Registration successful",
	}, nil
}

// Login 用户登录
func (s *userServiceImpl) Login(phone, password string) (*APIResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid phone or password")
	}

	// 校验密码
	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid phone or password")
	}

	// 生成JWT
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &APIResponse{
		Code:    0,
		Message: "Login successful",
		Data:    fiber.Map{"token": token},
	}, nil
}

// ChangePassword 修改密码
func (s *userServiceImpl) ChangePassword(userID int64, newPassword string) (*APIResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 设置新密码
	if err := user.SetPassword(newPassword); err != nil {
		return nil, err
	}

	// 更新用户信息
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &APIResponse{
		Code:    0,
		Message: "Password changed successfully",
	}, nil
}
