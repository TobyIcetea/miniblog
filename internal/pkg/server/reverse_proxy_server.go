// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/TobyIcetea/miniblog/internal/pkg/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	genericoptions "github.com/onexstack/onexstack/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// GRPCGatewayServer 代表一个 GRPC 网关服务器
type GRPCGatewayServer struct {
	srv *http.Server
}

// NewGRPCGatewayServer 创建一个新的 GRPC 网关服务器实例
func NewGRPCGatewayServer(
	httpOptions *genericoptions.HTTPOptions,
	grpcOptions *genericoptions.GRPCOptions,
	registerHandler func(mux *runtime.ServeMux, conn *grpc.ClientConn) error,
) (*GRPCGatewayServer, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 10 * time.Second, // 最小连接时间
		}),
	}
	dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(grpcOptions.Addr, dialOptions...)
	if err != nil {
		log.Errorw("Failed to connect GRPC server", "err", err)
		return nil, err
	}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			// 设置序列化 protobuf 数据时，枚举类型的字段以数字格式输出
			// 否则，默认会以字符串格式输出，跟枚举类型定义不一致，带来理解成本
			UseEnumNumbers: true,
		},
	}))

	if err := registerHandler(gwmux, conn); err != nil {
		log.Errorw("Failed to register handler", "err", err)
		return nil, err
	}

	return &GRPCGatewayServer{
		srv: &http.Server{
			Addr:    httpOptions.Addr,
			Handler: gwmux,
		},
	}, nil
}

// RunOrDie 启动 GRPC 网关服务器并在出错时记录致命错误
func (s *GRPCGatewayServer) RunOrDie() {
	log.Infow("Start to listening the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("Failed to start GRPC gateway server", "err", err)
	}
}

// GracefulStop 优雅地关闭 GRPC 网关服务器
func (s *GRPCGatewayServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop grpc gateway server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("Failed to shutdown GRPC gateway server", "err", err)
	}
}
