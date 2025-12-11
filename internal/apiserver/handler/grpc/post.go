// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package grpc

import (
	"context"

	apiv1 "github.com/TobyIcetea/miniblog/pkg/api/apiserver/v1"
)

// CreatePost 创建博客帖子.
func (h *Handler) CreatePost(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error) {
	return h.biz.PostV1().Create(ctx, rq)
}

// UpdatePost 更新博客帖子.
func (h *Handler) UpdatePost(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error) {
	return h.biz.PostV1().Update(ctx, rq)
}

// DeletePost 删除博客帖子.
func (h *Handler) DeletePost(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error) {
	return h.biz.PostV1().Delete(ctx, rq)
}

// GetPost 获取博客帖子.
func (h *Handler) GetPost(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error) {
	return h.biz.PostV1().Get(ctx, rq)
}

// ListPost 列出所有博客帖子.
func (h *Handler) ListPost(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error) {
	return h.biz.PostV1().List(ctx, rq)
}
