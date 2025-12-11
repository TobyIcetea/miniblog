// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package gin

import (
	"github.com/TobyIcetea/miniblog/internal/pkg/contextx"
	"github.com/TobyIcetea/miniblog/internal/pkg/known"
	"github.com/gin-gonic/gin"
	"github.com/onexstack/onexstack/pkg/log"
)

// 用于从 gin.Context 的 Header 中提取用户 ID，模拟所有请求认证通过.
func AuthnBypassMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中提取用户 ID，假设请求头名称为 "X-User-ID"
		userID := "user-000001" // 默认用户 ID
		if val := c.GetHeader(known.XUserID); val != "" {
			userID = val
		}

		log.Debugw("Simulated authentication bypass", "userID", userID)

		// 将用户 ID 和用户名注入到上下文中
		ctx := contextx.WithUserID(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)

		// 继续后续的操作
		c.Next()
	}
}
