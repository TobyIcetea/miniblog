// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package apiserver

import (
	"context"

	"github.com/TobyIcetea/miniblog/internal/pkg/server"
)

// ginServer 定义一个使用 Gin 框架开发的 HTTP 服务器
type ginServer struct{}

// 确保 *ginServer 实现了 server.Server 接口
var _ server.Server = (*ginServer)(nil)

// NewGinServer 初始化一个新的 Gin 服务器实例
func (cfg *ServerConfig) NewGinServer() server.Server {
	return &ginServer{}
}

// RunOrDie 启动 Gin 服务器，出错则程序崩溃退出
func (s *ginServer) RunOrDie() {
	select {}
}

// GracefulStop 优雅停止服务器
func (s *ginServer) GracefulStop(ctx context.Context) {}
