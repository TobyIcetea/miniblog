// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package http

import (
	"github.com/TobyIcetea/miniblog/internal/apiserver/biz"
	"github.com/TobyIcetea/miniblog/internal/pkg/validation"
)

// Handler 处理博客模块的请求.
type Handler struct {
	biz biz.IBiz
	val *validation.Validator
}

// NewHandler 创建新的 Handler 实例.
func NewHandler(biz biz.IBiz, val *validation.Validator) *Handler {
	return &Handler{
		biz: biz,
		val: val,
	}
}
