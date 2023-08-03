package validatex

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/uc1024/f90/core/validatex/idvalidator"
	"github.com/uc1024/f90/core/validatex/uniform"
)

// * 是否为手机
var rx_mobile = regexp.MustCompile("^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,1,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$")

func ValidateIsMobile(fl validator.FieldLevel) bool {
	return rx_mobile.MatchString(fl.Field().String())
}

// * 统一信用代码
func ValidateIsUniformCode(fl validator.FieldLevel) bool {
	return uniform.CalibrationUniform321002015(fl.Field().String())
}

// * 身份证号
func ValidateIsIDCard(fl validator.FieldLevel) bool {
	return idvalidator.IsValidCitizenNo(fl.Field().String())
}

var iso8601 = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?([+-]\d{2}:?\d{2}|Z)$`)

// * 是否为ISO8601时间格式
func ValidateIsISO8601(fl validator.FieldLevel) bool {
	return iso8601.MatchString(fl.Field().String())
}

// fl >= where
// * 字符串数字大于等于
func ValidateNumericGte(fl validator.FieldLevel) bool {
	n, err := decimal.NewFromString(fl.Field().String())
	if err != nil {
		return false
	}
	p, err := decimal.NewFromString(strings.TrimSpace(fl.Param()))
	if err != nil {
		return false
	}
	return n.GreaterThanOrEqual(p)
}

// fl <= where
// * 字符串数字小于等于
func ValidateNumericLte(fl validator.FieldLevel) bool {
	n, err := decimal.NewFromString(fl.Field().String())
	if err != nil {
		return false
	}
	p, err := decimal.NewFromString(strings.TrimSpace(fl.Param()))
	if err != nil {
		return false
	}
	return n.LessThanOrEqual(p)
}
