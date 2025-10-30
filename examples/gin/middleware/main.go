// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义一个通用中间件：打印请求路径
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Request Path: %s", c.Request.URL.Path)
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// 使用全局中间件，所有路由都会经过该中间件
	r.Use(LogMiddleware())

	// 定义普通路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Home",
		})
	})

	// 定义一个路由组，并且组添加中间件
	apiGroup := r.Group("/api", LogMiddleware())
	apiGroup.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, API!",
		})
	})
	apiGroup.GET("/world", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "World, API!",
		})
	})

	// 为单个路由添加中间件
	r.GET("/secure", LogMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Secure Resource",
		})
	})

	r.Run(":8080")
}
