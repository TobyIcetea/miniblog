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
	"strings"
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
	addr  = flag.String("addr", "localhost:6666", "The grpc server address to connect to.")
	limit = flag.Int64("limit", 10, "Limit to list users.")
)

func main() {
	flag.Parse()

	// 建立与 gRPC 服务器的链接
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close() // 确保连接在函数关闭时结束

	client := apiv1.NewMiniBlogClient(conn) // 创建 MiniBlog 客户端

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_ = uuid.New().String()

	createUserRequest := helper.ExampleCreateUserRequest()
	createUserRequest.Nickname = nil // 不设置 Nickname 字段
	createUserResponse, err := client.CreateUser(ctx, createUserRequest)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Printf("[CreateUser	] Success to create user: %v", createUserResponse)

	loginResponse, err := client.Login(ctx, &apiv1.LoginRequest{
		Username: createUserRequest.Username,
		Password: createUserRequest.Password,
	})
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	if loginResponse.Token == "" {
		log.Printf("Failed to validate login response: %v", loginResponse)
		return
	}
	log.Printf("[Login			] Success to login with user account: %v", loginResponse)

	// 创建 metadata，用于传递 token
	md := metadata.Pairs("Authorization", "Bearer "+loginResponse.Token, known.XUserID, createUserResponse.UserID)
	// 将 metadata 附加到上下文中
	ctx = metadata.NewOutgoingContext(ctx, md)

	defer func() {
		_, _ = client.DeleteUser(ctx, &apiv1.DeleteUserRequest{UserID: createUserResponse.UserID})
	}()

	getUserResponse, err := client.GetUser(ctx, &apiv1.GetUserRequest{UserID: createUserResponse.UserID})
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return
	}
	if getUserResponse.User.Nickname != "你好世界" {
		log.Printf("error default value for nickname")
		return
	}
	log.Printf("[GetUser		] Success to get user: %v", getUserResponse)

	createUserRequest2 := helper.ExampleCreateUserRequest()
	createUserRequest2.Email = "bad email address" // 不设置 nickname 字段
	_, err = client.CreateUser(ctx, createUserRequest2)
	if !strings.Contains(err.Error(), "invalid email format") {
		log.Printf("error create user with invalid email format: %v", err)
		return
	}
	log.Printf("[GetUser		] Success to create user with invalid email format: %v", err)
}
