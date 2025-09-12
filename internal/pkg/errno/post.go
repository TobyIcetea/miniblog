// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/TobyIcetea/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package errno

import (
	"net/http"

	"github.com/TobyIcetea/miniblog/pkg/errorsx"
)

// ErrPostNotFound 表示未找到指定的博客
var ErrPostNotFound = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.PostNotFound", Message: "Post not found."}
