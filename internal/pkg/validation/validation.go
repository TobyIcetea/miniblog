// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package validation

import (
	"regexp"

	"github.com/TobyIcetea/miniblog/internal/apiserver/store"
	"github.com/TobyIcetea/miniblog/internal/pkg/errno"
)

// Validator 是验证逻辑的实现结构体
type Validator struct {
	// 有些复杂的验证逻辑，可能需要直接查询数据库
	// 这里只是一个举例，如果验证时，有其他依赖的客户端/服务/资源等,
	// 都可以一并注入进来
	store store.IStore
}

// New 创建一个新的 Validator 实例
func New(store store.IStore) *Validator {
	return &Validator{store: store}
}

// isValidUsername 校验用户名是否合法
func isValidUsername(username string) bool {
	// 用户名必须仅包含字母、数字和下划线，并且长度在 3 到 20 个字符之间
	var (
		lengthRegex = `^.{3,20}$`       // 长度在 3 到 20 个字符之间
		validRegex  = `^[A-Za-z0-9_]+$` // 仅包含字母、数字和下划线
	)

	// 校验长度
	if matched, _ := regexp.MatchString(lengthRegex, username); !matched {
		return false
	}
	// 校验字符合法性
	if matched, _ := regexp.MatchString(validRegex, username); !matched {
		return false
	}

	return true
}

// isValidPassword 判断密码是否符合复杂度要求
func isValidPassword(password string) error {
	// 检查新密码是否为空
	if password == "" {
		return errno.ErrInvalidArgument.WithMessage("password cannot be empty")
	}

	// 检查新密码的长度要求
	if len(password) < 8 {
		return errno.ErrInvalidArgument.WithMessage("password must be at least 8 characters long")
	}

	// 使用正则表达式检查是否至少包含一个字母
	letterPattern := regexp.MustCompile(`[A-Za-z]`)
	if !letterPattern.MatchString(password) {
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one letter")
	}

	// 使用正则表达式检查是否至少包含一个数字
	numberPattern := regexp.MustCompile(`\d`)
	if !numberPattern.MatchString(password) {
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one number")
	}

	return nil
}

// isValidEmail 判断电子邮件是否合法
func isValidEmail(email string) error {
	// 检查电子邮件地址格式
	if email == "" {
		return errno.ErrInvalidArgument.WithMessage("email cannot be empty")
	}

	// 使用正则表达式校验电子邮件格式
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(email) {
		return errno.ErrInvalidArgument.WithMessage("invalid email format")
	}

	return nil
}

// isValidPhone 判断手机号码是否合法
func isValidPhone(phone string) error {
	// 检查手机号码格式
	if phone == "" {
		return errno.ErrInvalidArgument.WithMessage("phone cannot be empty")
	}

	// 使用正则表达式校验手机号码格式（假设是中国手机号，11 位数字）
	phonePattern := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phonePattern.MatchString(phone) {
		return errno.ErrInvalidArgument.WithMessage("invalid phone format")
	}

	return nil
}
