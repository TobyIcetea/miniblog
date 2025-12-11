// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package options

import (
	"errors"
	"fmt"
	"time"

	"github.com/TobyIcetea/miniblog/internal/apiserver"
	genericoptions "github.com/onexstack/onexstack/pkg/options"
	stringsutil "github.com/onexstack/onexstack/pkg/util/strings"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
)

// 定义支持的服务器模式集合.
var availableServerModes = sets.New(
	apiserver.GinServerMode,
	apiserver.GRPCServerMode,
	apiserver.GRPCGatewayServerMode,
)

// ServerOptions 包含服务器配置选项.
type ServerOptions struct {
	// ServerMode 定义服务器模式：gRPC、Gin HTTP、HTTP Reverse Proxy
	ServerMode string `json:"server-mode" mapstructure:"server-mode"`
	// JWTKey 定义 JWT 密钥
	JWTKey string `json:"jwt-key" mapstructure:"jwt-key"`
	// Expiration 定义 JWT Token 的过期时间
	Expiration time.Duration `json:"expiration" mapstructure:"expiration"`
	// TLSOptions 包含 TLS 配置选项
	TLSOptions *genericoptions.TLSOptions `json:"tls" mapstructure:"tls"`
	// HTTPOptions 包含 HTTP 配置选项
	HTTPOptions *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`
	// GRPCOptions 包含 gRPC 配置选项
	GRPCOptions *genericoptions.GRPCOptions `json:"grpc" mapstructure:"grpc"`
	// MySQLOptions 包含 MySQL 配置选项
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
}

// NewServerOptions 创建带有默认值的 ServerOptions 实例.
func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode:   apiserver.GRPCGatewayServerMode,
		JWTKey:       "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration:   2 * time.Hour,
		TLSOptions:   genericoptions.NewTLSOptions(),
		HTTPOptions:  genericoptions.NewHTTPOptions(),
		GRPCOptions:  genericoptions.NewGRPCOptions(),
		MySQLOptions: genericoptions.NewMySQLOptions(),
	}
	opts.HTTPOptions.Addr = ":5555"
	opts.GRPCOptions.Addr = ":6666"
	opts.MySQLOptions.Username = "miniblog"
	opts.MySQLOptions.Password = "miniblog1234"
	opts.MySQLOptions.Database = "miniblog"
	return opts
}

// 通过使用 pflag 包，可以实现从命令行中解析这些选项的功能.
func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, fmt.Sprintf("Server mode, available options: %v", availableServerModes.UnsortedList()))
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "JWT key for authentication. Must be at least 6 characters long.")
	// 绑定 JWT Token 的过期时间选项到命令行标志
	// 参数名称为 --expiration，默认值为 o.Expiration
	fs.DurationVar(&o.Expiration, "expiration", o.Expiration, "JWT Token expiration time.")
	o.TLSOptions.AddFlags(fs)
	o.HTTPOptions.AddFlags(fs)
	o.GRPCOptions.AddFlags(fs)
	o.MySQLOptions.AddFlags(fs)
}

// Validate 校验 ServerOptions 中的选项是否合法.
func (o *ServerOptions) Validate() error {
	errs := []error{}

	// 校验 ServerMode 是否有效
	if !availableServerModes.Has(o.ServerMode) {
		errs = append(errs, fmt.Errorf("invalid server mode: %s", o.ServerMode))
	}

	// 校验 JWTKey 是否至少 6 个字符长
	if len(o.JWTKey) < 6 {
		errs = append(errs, errors.New("jwt-key must be at least 6 characters long"))
	}

	// 校验子选项
	errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)

	// 如果是 gRPC 或 gRPC-Gateway 模式，校验 gRPC 配置
	if stringsutil.StringIn(o.ServerMode, []string{apiserver.GRPCServerMode, apiserver.GRPCGatewayServerMode}) {
		errs = append(errs, o.GRPCOptions.Validate()...)
	}

	// 合并所有错误并返回
	return utilerrors.NewAggregate(errs)
}

// Config 基于 ServerOptions 构建 apiserver.Config.
func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode:   o.ServerMode,
		JWTKey:       o.JWTKey,
		Expiration:   o.Expiration,
		TLSOptions:   o.TLSOptions,
		HTTPOptions:  o.HTTPOptions,
		GRPCOptions:  o.GRPCOptions,
		MySQLOptions: o.MySQLOptions,
	}, nil
}
