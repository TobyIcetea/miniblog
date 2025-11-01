// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package conversion

import (
	"github.com/TobyIcetea/miniblog/internal/apiserver/model"
	apiv1 "github.com/TobyIcetea/miniblog/pkg/api/apiserver/v1"
	"github.com/onexstack/onexstack/pkg/core"
)

// PostModelToPostV1 将模型层的 PostM（博客模型对象）转换为 Protobuf 层的 Post（v1 博客对象）
func PostModelToPostV1(postModel *model.PostM) *apiv1.Post {
	var protoPost apiv1.Post
	_ = core.CopyWithConverters(&protoPost, postModel)
	return &protoPost
}

// PostV1ToPostModel 将 Protobuf 层的 Post（v1 博客对象）转换为模型层的 PostM（博客模型对象）
func PostV1ToPostModel(protoPost *apiv1.Post) *model.PostM {
	var postModel model.PostM
	_ = core.CopyWithConverters(&postModel, protoPost)
	return &postModel
}
