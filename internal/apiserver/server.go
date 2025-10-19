// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package apiserver

import (
	"net"
	"time"

	handler "github.com/TobyIcetea/miniblog/internal/apiserver/handler/grpc"
	"github.com/TobyIcetea/miniblog/internal/pkg/log"
	v1 "github.com/TobyIcetea/miniblog/pkg/api/apiserver/v1"
	genericoptions "github.com/onexstack/onexstack/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	cfg *Config
	srv *grpc.Server
	lis net.Listener
}

// NewUnionServer 根据配置创建联合服务器
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	lis, err := net.Listen("tcp", cfg.GRPCOptions.Addr)
	if err != nil {
		log.Fatalw("Failed to listen on %s: %v", cfg.GRPCOptions.Addr, err)
		return nil, err
	}

	// 创建 GRPC Server 实例
	grpcsrv := grpc.NewServer()
	v1.RegisterMiniBlogServer(grpcsrv, handler.NewHandler())
	reflection.Register(grpcsrv)

	return &UnionServer{cfg: cfg, srv: grpcsrv, lis: lis}, nil
}

// Run 运行应用
func (s *UnionServer) Run() error {
	// 打印一条日志，用来提示 GRPC 服务已经起来，方便排障
	log.Infow("Start to listening the incoming requests on grpc address", "addr", s.cfg.GRPCOptions.Addr)
	return s.srv.Serve(s.lis)
}
