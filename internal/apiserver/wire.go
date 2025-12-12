// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

//go:build wireinject

package apiserver

import (
	"github.com/TobyIcetea/miniblog/internal/apiserver/biz"
	"github.com/TobyIcetea/miniblog/internal/apiserver/store"
	ginmw "github.com/TobyIcetea/miniblog/internal/pkg/middleware/gin"
	"github.com/TobyIcetea/miniblog/internal/pkg/server"
	"github.com/TobyIcetea/miniblog/internal/pkg/validation"
	"github.com/TobyIcetea/miniblog/pkg/auth"
	"github.com/google/wire"
)

func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		ProvideDB, // 提供数据库实例
		validation.ProviderSet,
		wire.NewSet(
			wire.Struct(new(UserRetriever), "*"),
			wire.Bind(new(ginmw.UserRetriever), new(*UserRetriever)),
		),
		auth.ProviderSet,
	)
	return nil, nil
}
