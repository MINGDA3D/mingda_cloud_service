package handler

import (

	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/app/service"
	"mingda_cloud_service/internal/pkg/response"
	"mingda_cloud_service/internal/pkg/validator"
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