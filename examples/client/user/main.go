// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/TobyIcetea/miniblog/examples/helper"
	"github.com/TobyIcetea/miniblog/internal/pkg/known"
	apiv1 "github.com/TobyIcetea/miniblog/pkg/api/apiserver/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"k8s.io/utils/ptr"
)

var (
	addr  = flag.String("addr", "localhost:6666", "The grpc server address to connect to")
	limit = flag.Int64("limit", 10, "Limit to list users")
)

func main() {
	flag.Parse()

	// 建立与 gRPC 服务器之间的连接
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	defer conn.Close() // 确保连接在函数关闭时关闭

	client := apiv1.NewMiniBlogClient(conn) // 创建 MiniBlog 客户端

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_ = uuid.New().String()

	createUserRequest := helper.ExampleCreateUserRequest()
	createUserResponse, err := client.CreateUser(ctx, createUserRequest)
	if err != nil {
		log.Fatalf("Failed to create user: %v, username: %s", err, createUserRequest.Username)
	}
	log.Printf("[CreateUser			] Success to create user, username: %s, user_id: %s", createUserRequest.Username, createUserResponse.UserID)

	loginResponse, err := client.Login(ctx, &apiv1.LoginRequest{
		Username: createUserRequest.Username,
		Password: createUserRequest.Password,
	})
	if err != nil {
		log.Fatalf("Failed to login user: %v, username: %s", err, createUserRequest.Username)
	}
	if loginResponse.Token == "" {
		log.Printf("[Login			] Failed to login user, username: %s", createUserRequest.Username)
		return
	}
	log.Printf("[Login			] Success to login user, username: %s, token: %s", createUserRequest.Username, loginResponse.Token)

	// 创建 metadata，用于传递 Token
	md := metadata.Pairs("Authroization", "Bearer "+loginResponse.Token, known.XUserID, createUserResponse.UserID)
	// 将 metadata 附加到上下文中
	ctx = metadata.NewOutgoingContext(ctx, md)

	defer func() {
		_, _ = client.DeleteUser(ctx, &apiv1.DeleteUserRequest{UserID: createUserResponse.UserID})
	}()

	refreshTokenResponse, err := client.RefreshToken(ctx, &apiv1.RefreshTokenRequest{})
	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
		return
	}
	if refreshTokenResponse.Token == "" {
		log.Printf("Token cannot be empty")
		return
	}
	log.Printf("[RefreshToken			] Success to refresh token, token: %s", refreshTokenResponse.Token)

	// 请求 UpdateUser 接口
	_, err = client.UpdateUser(ctx, &apiv1.UpdateUserRequest{
		UserID:   createUserResponse.UserID,
		Nickname: ptr.To("TobyMint"),
	})
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return
	}
	log.Printf("[UpdateUser     ] Success to update user: %v", createUserResponse.UserID)

	// 请求 ChangePassword 接口
	newPassword := "miniblog123456"
	_, err = client.ChangePassword(ctx, &apiv1.ChangePasswordRequest{
		UserID:      createUserResponse.UserID,
		OldPassword: createUserRequest.Password,
		NewPassword: newPassword,
	})
	if err != nil {
		log.Printf("Failed to change password: %v", err)
		return
	}
	log.Printf("[ChangePassword ] Success to change password: %v", createUserResponse.UserID)

	loginResponse, err = client.Login(ctx, &apiv1.LoginRequest{
		Username: createUserRequest.Username,
		Password: newPassword,
	})
	if err != nil {
		log.Printf("Failed to login with new password: %v", err)
		return
	}
	log.Printf("[Login ] Success to login with new password, username: %s, token: %s", createUserRequest.Username, loginResponse.Token)

	// 创建 metadata，用于传递 Token
	md = metadata.Pairs("Authorization", "Bearer "+loginResponse.Token, known.XUserID, createUserResponse.UserID)
	// 将 metadata 附加到上下文中
	ctx = metadata.NewOutgoingContext(ctx, md)

	getUserResponse, err := client.GetUser(ctx, &apiv1.GetUserRequest{UserID: createUserResponse.UserID})
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return
	}
	if getUserResponse.User.UserID != createUserResponse.UserID || getUserResponse.User.Username != createUserRequest.Username {
		log.Printf("Failed to get user, user_id: %s, username: %s", getUserResponse.User.UserID, getUserResponse.User.Username)
		return
	}
	log.Printf("[GetUser	] Success to get user, user_id: %s, username: %s", getUserResponse.User.UserID, getUserResponse.User.Username)

	// 这里的逻辑我觉得实在是有点问题，regular user 应该可以访问倒一些东西，比如说可以访问到自己
	// 但是这里的逻辑是，只要调用了 ListUser 方法，就直接报错
	listResponse, err := client.ListUser(ctx, &apiv1.ListUserRequest{Offset: 0, Limit: *limit})
	if err != nil {
		log.Printf("[ListUser		] Failed to list user: %v", err)
	} else {
		if len(listResponse.Users) == 1 && listResponse.Users[0].UserID == createUserResponse.UserID {
			log.Printf("[ListUser       ] Normal users can only list themselves")
		} else {
			log.Printf("[ListUser       ] Failed to validate permission: regular users can access the user list")
			return
		}
	}

	ctx = helper.MustWithAdminToken(ctx, client)

	// 请求 ListUser 接口
	listResponse, err = client.ListUser(ctx, &apiv1.ListUserRequest{Offset: 0, Limit: *limit})
	if err != nil {
		log.Printf("Failed to list user: %v", err)
		return
	}
	log.Printf("[ListUser	] Success to list user, count: %d", len(listResponse.Users))
	found := false
	for _, user := range listResponse.Users {
		if user.UserID == createUserResponse.UserID && user.Username == createUserRequest.Username {
			found = true
			break
		}
	}
	if found {
		log.Printf("[ListUser	] Success to find user, user_id: %s, username: %s", createUserResponse.UserID, createUserRequest.Username)
	}

	// 请求 DeleteUser 接口
	_, err = client.DeleteUser(ctx, &apiv1.DeleteUserRequest{UserID: createUserResponse.UserID})
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return
	}
	log.Printf("[DeleteUser	] Success to delete user, user_id: %s", createUserResponse.UserID)

	log.Printf("[All Test Passed] Success to test all user api")
}

// 随机生成一个符合中国大陆手机格式的号码
func GeneratePhoneNumber() string {
	// 手机号码规则：以 1 开头，第二位为 3-9，接下来 9 位随机数字
	prefixes := []int{3, 4, 5, 6, 7, 8, 9}

	// 随机选择第二位
	secondDigit := prefixes[rand.IntN(len(prefixes))]

	// 随机生成后 9 位数字
	phone := fmt.Sprintf("1%d", secondDigit)
	for i := 0; i < 9; i++ {
		phone += fmt.Sprintf("%d", rand.IntN(10)) // 随机生成剩余的 9 位数字
	}

	return phone
}
