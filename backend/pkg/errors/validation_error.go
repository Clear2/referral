package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationMessage converts binding and validation errors into messages safe
// to show to API consumers. Detailed validator internals should stay in logs.
func ValidationMessage(err error) string {
	if err == nil {
		return "请求参数不正确"
	}

	if errors.Is(err, io.EOF) {
		return "请求内容不能为空"
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return "请求内容不是有效的 JSON"
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return fmt.Sprintf("%s的数据类型不正确", fieldLabel(typeErr.Field))
	}

	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return "请求参数不正确"
	}

	messages := make([]string, 0, len(validationErrs))
	for _, fieldErr := range validationErrs {
		field := fieldLabel(fieldErr.Field())
		switch fieldErr.Tag() {
		case "required":
			messages = append(messages, field+"不能为空")
		case "email":
			messages = append(messages, "请输入有效的邮箱地址")
		case "min":
			messages = append(messages, fmt.Sprintf("%s至少需要%s位", field, fieldErr.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s不能超过%s位", field, fieldErr.Param()))
		case "url":
			messages = append(messages, field+"必须是有效的网址")
		case "oneof":
			messages = append(messages, field+"的取值不受支持")
		default:
			messages = append(messages, field+"格式不正确")
		}
	}

	return strings.Join(messages, "；")
}

func fieldLabel(field string) string {
	switch strings.ToLower(field) {
	case "name":
		return "用户名"
	case "email":
		return "邮箱"
	case "password":
		return "密码"
	case "redirecturl", "redirect_url", "redirect":
		return "跳转地址"
	case "client":
		return "登录方式"
	case "lang":
		return "语言"
	case "state":
		return "登录状态"
	case "code":
		return "授权码"
	case "id":
		return "ID"
	case "title":
		return "标题"
	case "content":
		return "内容"
	case "category":
		return "分类"
	case "owner":
		return "所有者"
	case "expiretime":
		return "过期时间"
	default:
		return field
	}
}
