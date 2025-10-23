// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package apiserver

import (
	"time"

	"github.com/TobyIcetea/miniblog/internal/pkg/server"
	"github.com/onexstack/onexstack/pkg/log"
	genericoptions "github.com/onexstack/onexstack/pkg/options"
)

const (
	// GRPCServerMode 定义 gRPC 服务模式
	// 使用 gRPC 框架启动一个 gRPC 服务器
	GRPCServerMode = "grpc"
	// GRPCGatewayServerMode 定义 gRPC + HTTP 服务模式
	// 使用 gRPC 框架启动一个 gRPC 服务器 + HTTP 反向代理服务器
	GRPCGatewayServerMode = "grpc-gateway"
	// GinServerMode 定义 Gin 服务模式
	// 使用 Gin Web 框架启动一个 HTTP 服务器
	GinServerMode = "gin"
)

// Config 配置结构体，用于存储应用相关的配置
// 不用 viper.Get，是因为这种方式能更加清晰的知道应用提供了哪些配置项
type Config struct {
	ServerMode  string
	JWTKey      string
	Expiration  time.Duration
	HTTPOptions *genericoptions.HTTPOptions
	GRPCOptions *genericoptions.GRPCOptions
}

// UnionServer 定义了一个联合服务器。根据 ServerMode 决定要启动的服务器类型
//
// 联合服务器分为以下 2 大类：
//  1. Gin 服务器：由 Gin 框架创建的标准的 REST 服务器。根据是否开启 TLS,
//     来判断启动 HTTP 或者 HTTPS
//  2. GRPC 服务器：由 gRPC 框架创建的标准的 RPC 服务器
//  3. HTTP 反向代理服务器：由 grpc-gateway 框架创建的 HTTP 反向代理服务器
//     根据是否开启 TLS，来判断启动 HTTP 或者 HTTPS
//
// HTTP 反向代理服务器依赖 gRPC 服务器，所以在开启 HTTP 反向代理服务器时，会先启动 gRPC 服务器
type UnionServer struct {
	srv server.Server
}

// ServerConfig 包含服务器的核心依赖和配置
type ServerConfig struct {
	cfg *Config
}

// NewUnionServer 根据配置创建联合服务器
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	// 一些初始化代码

	// 创建服务器配置，这些配置可用来创建服务器
	serverConfig, err := cfg.NewServerConfig()
	if err != nil {
		return nil, err
	}

	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode)

	// 根据服务模式创建对应的服务实例
	// 实际企业开发中，可以根据需要只选择一种服务器模式
	// 这里为了方便展示，通过 cfg.ServerMode 同时支持了 Gin 和 GRPC 2 种服务器模式
	// 默认为 gRPC 服务器模式
	var srv server.Server
	switch cfg.ServerMode {
	case GinServerMode:
		srv, err = serverConfig.NewGinServer(), nil
	default:
		srv, err = serverConfig.NewGRPCServerOr()
	}
	if err != nil {
		return nil, err
	}

	return &UnionServer{srv: srv}, nil
}

// Run 运行应用
func (s *UnionServer) Run() error {
	s.srv.RunOrDie()
	return nil
}

// NewServerConfig 创建一个 *ServerConfig 示例
func (cfg *Config) NewServerConfig() (*ServerConfig, error) {
	return &ServerConfig{cfg: cfg}, nil
}
