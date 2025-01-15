package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	if e, ok := err.(*errors.Error); ok {
		c.JSON(http.StatusOK, Response{
			Code:    int(e.Code),
			Message: e.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Code:    int(errors.ErrUnknown),
		Message: err.Error(),
	})
}

// ErrorWithMsg 带消息的错误响应
func ErrorWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.ErrUnknown),
		Message: msg,
	})
} 