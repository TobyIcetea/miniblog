// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package main

import (
	"errors"
	"fmt"
)

type LoginRequest struct {
	Username string
	Password string
}

func validate(req LoginRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

func main() {
	req := LoginRequest{
		Username: "user",
		Password: "12345",
	}
	err := validate(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Validation passed!")
}
