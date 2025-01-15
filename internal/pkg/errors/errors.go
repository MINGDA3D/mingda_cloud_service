package errors

// ErrorCode 错误码类型
type ErrorCode int

const (
	// 系统级错误码 (1000-1999)
	ErrUnknown         ErrorCode = 1000 // 未知错误
	ErrInvalidParams   ErrorCode = 1001 // 参数错误
	ErrUnauthorized    ErrorCode = 1002 // 未授权
	ErrInvalidToken    ErrorCode = 1003 // 无效的Token
	ErrTokenExpired    ErrorCode = 1004 // Token已过期
	ErrInvalidSign     ErrorCode = 1005 // 签名错误
	ErrTooManyReq      ErrorCode = 1006 // 请求过于频繁
	ErrTimeout         ErrorCode = 1007 // 请求超时
	ErrServiceBusy     ErrorCode = 1008 // 服务繁忙

	// 业务错误码 (2000-2999)
	// 设备相关 (2000-2099)
	ErrDeviceNotFound    ErrorCode = 2001 // 设备不存在
	ErrDeviceDisabled    ErrorCode = 2002 // 设备已禁用
	ErrInvalidSN         ErrorCode = 2003 // 无效的SN码
	ErrDeviceTypeInvalid ErrorCode = 2004 // 无效的设备类型
	ErrDeviceOffline     ErrorCode = 2005 // 设备离线
	ErrDeviceBusy        ErrorCode = 2006 // 设备忙
	ErrDeviceTimeout     ErrorCode = 2007 // 设备响应超时
	
	// 软件版本相关 (2100-2199)
	ErrInvalidVersion    ErrorCode = 2101 // 无效的版本号
	ErrVersionNotMatch   ErrorCode = 2102 // 版本不匹配
	ErrVersionTooOld     ErrorCode = 2103 // 版本过旧
	ErrVersionTooNew     ErrorCode = 2104 // 版本过新

	// 基础设施错误码 (3000-3999)
	// 数据库相关 (3000-3099)
	ErrDatabase         ErrorCode = 3001 // 数据库操作失败
	ErrDatabaseTimeout  ErrorCode = 3002 // 数据库超时
	ErrDatabaseConnect  ErrorCode = 3003 // 数据库连接失败
	ErrDatabaseDup      ErrorCode = 3004 // 数据重复
	
	// 缓存相关 (3100-3199)
	ErrRedis           ErrorCode = 3101 // Redis操作失败
	ErrRedisTimeout    ErrorCode = 3102 // Redis超时
	ErrRedisConnect    ErrorCode = 3103 // Redis连接失败
	ErrRedisNotFound   ErrorCode = 3104 // Redis键不存在
	
	// 消息队列相关 (3200-3299)
	ErrRabbitMQ        ErrorCode = 3201 // RabbitMQ操作失败
	ErrRabbitMQTimeout ErrorCode = 3202 // RabbitMQ超时
	ErrRabbitMQConnect ErrorCode = 3203 // RabbitMQ连接失败
)

// Error 自定义错误
type Error struct {
	Code    ErrorCode `json:"code"`    // 错误码
	Message string    `json:"message"` // 错误信息
}

// Error 实现error接口
func (e *Error) Error() string {
	return e.Message
}

// New 创建新的错误
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewWithError 使用已有错误创建新的错误
func NewWithError(code ErrorCode, err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Code:    code,
		Message: err.Error(),
	}
}

// IsErrorCode 判断错误是否为指定错误码
func IsErrorCode(err error, code ErrorCode) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
} 