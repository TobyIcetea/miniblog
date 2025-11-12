// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/TobyIcetea/miniblog/examples/helper"
	"github.com/TobyIcetea/miniblog/internal/pkg/known"
	apiv1 "github.com/TobyIcetea/miniblog/pkg/api/apiserver/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	addr  = flag.String("addr", "localhost:6666", "The grpc server address to connect to")
	limit = flag.Int64("limit", 10, "limit to list users")
)

func main() {
	flag.Parse()

	// 建立与 gRPC 服务器之间的链接
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	defer conn.Close() // 确保连接在函数结束时关闭

	client := apiv1.NewMiniBlogClient(conn) // 创建 MiniBlog 客户端

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_ = uuid.New().String()

	createUserRequest := helper.ExampleCreateUserRequest()
	createUserResponse, err := client.CreateUser(ctx, createUserRequest)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Printf("[CreateUser     ] Success to create user: %v", createUserResponse.UserID)

	loginResponse, err := client.Login(ctx, &apiv1.LoginRequest{
		Username: createUserRequest.Username,
		Password: createUserRequest.Password,
	})
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	log.Printf("[Login			] Success to login user, username: %s, token: %s", createUserRequest.Username, loginResponse.Token)

	// 创建 metadata，用于传递 Token
	md := metadata.Pairs("Authorization", "Bearer "+loginResponse.Token, known.XUserID, createUserResponse.UserID)
	// 将 metadata 附加到上下文中
	ctx = metadata.NewOutgoingContext(ctx, md)

	defer func() {
		ctx = helper.MustWithAdminToken(ctx, client)
		_, _ = client.DeleteUser(ctx, &apiv1.DeleteUserRequest{UserID: createUserResponse.UserID})
	}()

	// 请求 CreatePost 接口
	createPostRequest := &apiv1.CreatePostRequest{
		Title:   "Hello, World",
		Content: "This is a test blog of miniblog platform.",
	}
	createPostResponse, err := client.CreatePost(ctx, createPostRequest)
	if err != nil {
		log.Printf("[CreatePost		] Failed to create post: %v", err)
		return
	}
	log.Printf("[CreatePost		] Success to create post: %v", createPostResponse.PostID)

	// 请求 UpdatePost 接口
	newTitle := "Hello World Modified"
	_, err = client.UpdatePost(ctx, &apiv1.UpdatePostRequest{
		PostID: createPostResponse.PostID,
		Title:  &newTitle,
	})
	if err != nil {
		log.Printf("[UpdatePost		] Failed to update post: %v", err)
		return
	}
	log.Printf("[UpdatePost		] Success to update post: %v", createPostResponse.PostID)

	// 请求 GetPost 接口
	getPostResponse, err := client.GetPost(ctx, &apiv1.GetPostRequest{PostID: createPostResponse.PostID})
	if err != nil {
		log.Printf("[GetPost			] Failed to get post: %v", err)
		return
	}
	if getPostResponse.Post.PostID != createPostResponse.PostID || getPostResponse.Post.Title != newTitle {
		log.Printf("[GetPost			] Failed to get post, post_id: %s, title: %s", getPostResponse.Post.PostID, getPostResponse.Post.Title)
		return
	}
	log.Printf("[GetPost			] Success to get post, post_id: %s, title: %s", getPostResponse.Post.PostID, getPostResponse.Post.Title)

	// 请求 ListPost 接口
	listResponse, err := client.ListPost(ctx, &apiv1.ListPostRequest{Offset: 0, Limit: *limit})
	if err != nil {
		log.Printf("[ListPost		] Failed to list post: %v", err)
		return
	}
	log.Printf("[ListPost		] Success to list post: %v", listResponse.Posts)

	// 请求 DeletePost 接口
	_, err = client.DeletePost(ctx, &apiv1.DeletePostRequest{PostIDs: []string{createPostResponse.PostID}})
	if err != nil {
		log.Printf("[DeletePost		] Failed to delete post: %v", err)
		return
	}
	log.Printf("[DeletePost		] Success to delete post: %v", createPostResponse.PostID)

	log.Printf("[All		] Success")
}
