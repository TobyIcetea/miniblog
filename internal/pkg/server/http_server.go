// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/TobyIcetea/miniblog/internal/pkg/log"
	genericoptions "github.com/onexstack/onexstack/pkg/options"
)

// HTTPServer 代表一个 HTTP 服务器.
type HTTPServer struct {
	srv *http.Server
}

// NewHTTPServer 创建一个新的 HTTP 服务器实例.
func NewHTTPServer(httpOptions *genericoptions.HTTPOptions, tlsOptions *genericoptions.TLSOptions, handler http.Handler) *HTTPServer {
	var tlsConfig *tls.Config
	if tlsOptions != nil && tlsOptions.UseTLS {
		tlsConfig = tlsOptions.MustTLSConfig()
	}

	return &HTTPServer{
		srv: &http.Server{
			Addr:      httpOptions.Addr,
			Handler:   handler,
			TLSConfig: tlsConfig,
		},
	}
}

// RunOrDir 启动 HTTP 服务器并在出错时记录致命错误.
func (s *HTTPServer) RunOrDie() {
	log.Infow("Start to listening the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)

	// 默认启动 HTTP 服务器
	serveFn := func() error { return s.srv.ListenAndServe() }
	if s.srv.TLSConfig != nil {
		serveFn = func() error { return s.srv.ListenAndServeTLS("", "") }
	}

	if err := serveFn(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("Failed to start HTTP(s) server", "err", err)
	}
}

// GracefulStop 优雅地关闭 HTTP 服务器.
func (s *HTTPServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("Failed to stop HTTP server", "err", err)
	}
}
