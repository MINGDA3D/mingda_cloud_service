package handler

import (
	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/app/service"
	"mingda_cloud_service/internal/pkg/response"
	"mingda_cloud_service/internal/pkg/validator"
	"strings"
	"errors"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(jwtSecret, aesKey string) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(jwtSecret, aesKey),
	}
}

type RegisterRequest struct {
	SN    string `json:"sn" binding:"required,min=10,max=32"`
	Model string `json:"model" binding:"required"`
}

type AuthRequest struct {
	SN        string `json:"sn" binding:"required"`
	Sign      string `json:"sign" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

// Register 设备注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := validator.BindAndValid(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	device, err := h.authService.RegisterDevice(req.SN, req.Model)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, device)
}

// Authenticate 设备认证
func (h *AuthHandler) Authenticate(c *gin.Context) {
	var req AuthRequest
	if err := validator.BindAndValid(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	device, err := h.authService.AuthenticateDevice(req.SN, req.Sign, req.Timestamp)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 生成访问令牌
	token, err := h.authService.GenerateToken(device)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"device": device,
	})
}

// RefreshToken 刷新访问令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 从请求头获取当前token
	auth := c.GetHeader("Authorization")
	if auth == "" {
		response.Error(c, errors.New(errors.ErrUnauthorized, "missing authorization header"))
		return
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.Error(c, errors.New(errors.ErrUnauthorized, "invalid authorization format"))
		return
	}

	// 刷新token
	newToken, err := h.authService.RefreshToken(parts[1])
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"token": newToken,
	})
} 