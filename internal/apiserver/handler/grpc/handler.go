// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package grpc

import (
	"github.com/TobyIcetea/miniblog/internal/apiserver/biz"
	apiv1 "github.com/TobyIcetea/miniblog/pkg/api/apiserver/v1"
)

// Handler 负责处理博客模块的请求.
type Handler struct {
	apiv1.UnimplementedMiniBlogServer

	biz biz.IBiz
}

// NewHandler 创建一个新的 Handler 实例.
func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{biz: biz}
}
