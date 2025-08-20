package controllers

import (
	"emoji-maker-backend/services"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	NewPassword string `json:"newPassword"`
}

// Register 用户注册接口
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.APIResponse{
			Code:    1,
			Message: "Invalid request body",
		})
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	if match, _ := regexp.MatchString(phoneRegex, req.Phone); !match {
		return c.Status(fiber.StatusBadRequest).JSON(services.APIResponse{
			Code:    1,
			Message: "Invalid phone number format",
		})
	}

	// 验证密码强度
	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(services.APIResponse{
			Code:    1,
			Message: "Password must be at least 8 characters long",
		})
	}

	response, err := h.userService.Register(req.Phone, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.APIResponse{
			Code:    1,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Login 用户登录接口
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.APIResponse{
			Code:    1,
			Message: "Invalid request body",
		})
	}

	if req.Phone == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(services.APIResponse{
			Code:    1,
			Message: "Phone and password are required",
		})
	}

	response, err := h.userService.Login(req.Phone, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(services.APIResponse{
			Code:    1,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// ChangePassword 修改密码接口
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.APIResponse{
			Code:    1,
			Message: "Invalid request body",
		})
	}

	if req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(services.APIResponse{
			Code:    1,
			Message: "New password is required",
		})
	}

	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(services.APIResponse{
			Code:    1,
			Message: "Invalid user ID in token",
		})
	}

	response, err := h.userService.ChangePassword(userID, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.APIResponse{
			Code:    1,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
