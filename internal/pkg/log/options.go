// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package log

import "go.uber.org/zap/zapcore"

// 通过改结构体，可以自定义日志的输出格式、级别以及其他相关配置.
type Options struct {
	// DisableCaller 指定是否禁用 caller 信息
	// 如果设置为 false（默认值），日志中会显式调用日志所在的文件名和行号
	DisableCaller bool

	// DisableStacktrace 指定是否禁用堆栈信息
	// 如果设置为 false（默认值），在日志级别为 panic 或更高时，会打印堆栈跟踪信息
	DisableStacktrace bool

	// Level 指定日志级别
	// 可选值包括：debug、info、warn、error、dpanic、panic、fatal
	// 默认值为 Info
	Level string

	// Format 指定日志的输出格式
	// 可选值包括：console(控制台格式) 和 json(JSON 格式)
	// 默认值为 console
	Format string

	// OutputPaths 指定日志的输出位置
	// 默认值为标准输出(stdout)，也可以指定文件路径或其他输出目标
	OutputPaths []string
}

// 该方法用于初始化日志配置选项，提供默认的日志级别、格式和输出位置.
func NewOptions() *Options {
	return &Options{
		// 默认启用 caller 信息
		DisableCaller: false,
		// 默认启用堆栈信息
		DisableStacktrace: false,
		// 默认日志级别为 info
		Level: zapcore.InfoLevel.String(),
		// 默认日志输出格式为 console
		Format: "console",
		// 默认日志输出位置为标准输出
		OutputPaths: []string{"stdout"},
	}
}
