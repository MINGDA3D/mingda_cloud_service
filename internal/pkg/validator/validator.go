package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mingda_cloud_service/internal/pkg/errors"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateDeviceSN 验证设备SN码
func ValidateDeviceSN(sn string) error {
	// 1. 检查长度
	if len(sn) != 15 {
		return fmt.Errorf("SN码长度必须为15位")
	}

	// 2. 检查格式
	pattern := `^[A-Z]\d[A-Z][0-9]{2}[0-9]{2}[A-Z][0-9]{2}[0-9]{5}$`
	matched, _ := regexp.MatchString(pattern, sn)
	if !matched {
		return fmt.Errorf("SN码格式不正确")
	}

	// 3. 解析各部分
	modelCode := sn[0:3]    // 机型代码
	year := sn[3:5]         // 年份
	day := sn[5:7]          // 日期
	version := sn[7:8]      // 版本号
	month := sn[8:10]       // 月份
	sequence := sn[10:15]   // 流水号

	// 4. 验证机型代码
	if !isValidModelCode(modelCode) {
		return fmt.Errorf("无效的机型代码")
	}

	// 5. 验证日期
	yearNum, _ := strconv.Atoi("20" + year)
	monthNum, _ := strconv.Atoi(month)
	dayNum, _ := strconv.Atoi(day)
	
	if !isValidDate(yearNum, monthNum, dayNum) {
		return fmt.Errorf("无效的生产日期")
	}

	// 6. 验证版本号
	if !isValidVersion(version) {
		return fmt.Errorf("无效的版本号")
	}

	// 7. 验证流水号
	seqNum, _ := strconv.Atoi(sequence)
	if seqNum < 1 || seqNum > 99999 {
		return fmt.Errorf("无效的生产流水号")
	}

	return nil
}

// isValidModelCode 验证机型代码
func isValidModelCode(code string) bool {
	// 目前支持的机型代码列表
	validCodes := map[string]bool{
		"M1P": true, // MD-1000 PRO
		"M6P": true, // MD-600 PRO
		"M1D": true, // MD-1000D
		"M4D": true, // MD-400D
		"M6D": true, // MD-600D
		// 可以添加更多机型
	}
	return validCodes[code]
}

// isValidDate 验证日期是否有效
func isValidDate(year, month, day int) bool {
	// 检查年份范围（假设从2020年开始）
	if year < 2020 || year > time.Now().Year() {
		return false
	}

	// 检查月份范围
	if month < 1 || month > 12 {
		return false
	}

	// 获取指定年月的最大天数
	t := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
	maxDays := t.Day()

	// 检查日期范围
	if day < 1 || day > maxDays {
		return false
	}

	return true
}

// isValidVersion 验证版本号
func isValidVersion(version string) bool {
	// 版本号必须是A-Z的大写字母
	return version >= "A" && version <= "Z"
}

// ValidateStruct 验证结构体
func ValidateStruct(obj interface{}) error {
	if obj == nil {
		return errors.New(errors.ErrInvalidParams, "参数验证错误")
	}

	if err := validate.Struct(obj); err != nil {
		return errors.New(errors.ErrInvalidParams, err.Error())
	}

	return nil
}

// ValidateJSON 验证JSON参数
func ValidateJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.New(errors.ErrInvalidParams, "参数绑定错误")
	}
	return ValidateStruct(obj)
}

// BindAndValid 绑定并验证请求参数
func BindAndValid(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		return errors.New(errors.ErrInvalidParams, "参数绑定错误")
	}

	return ValidateStruct(obj)
} 