// Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/onexstack/miniblog. The professional
// version of this repository is https://github.com/onexstack/onex.

package rid_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/TobyIcetea/miniblog/internal/pkg/rid"
	"github.com/stretchr/testify/assert"
)

// Mock Salt function used for testing
func Salt() string {
	return "staticSalt"
}

func TestResourceID_String(t *testing.T) {
	// 测试 UserID 转换为字符串
	userID := rid.UserID
	assert.Equal(t, "user", userID.String(), "UserID.String() should return \"user\"")

	// 测试 PostID 转换为字符串
	postID := rid.PostID
	assert.Equal(t, "post", postID.String(), "PostID.String() should return \"post\"")
}

func TestResourceID_New(t *testing.T) {
	// 测试生成的 ID 是否带有正确前缀
	userID := rid.UserID
	uniqueID := userID.New(1)

	assert.True(t, len(uniqueID) > 0, "Generated UserID should not be empty")
	assert.Contains(t, uniqueID, "user-", "Generated UserID should start with \"user-\"")

	// 生成另外一个唯一标识符，确保生成的值不同
	anotherID := userID.New(2)
	assert.NotEqual(t, uniqueID, anotherID, "Generated UserID should be unique")
	fmt.Printf("Generated UserID: %s\n", uniqueID)
	fmt.Printf("Generated Another UserID: %s\n", anotherID)
}

func BenchmarkResourceID_New(b *testing.B) {
	// 性能测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := rid.UserID
		_ = userID.New(uint64(i))
	}
}

func FuzzResourceID_New(f *testing.F) {
	// 添加预制测试数据
	f.Add(uint64(1))      // 添加一个种子值 counter 为 1
	f.Add(uint64(123456)) // 添加一个较大的种子值

	f.Fuzz(func(t *testing.T, counter uint64) {
		// 测试 UserID 的 New 方法
		result := rid.UserID.New(counter)

		// 断言结果不为空
		assert.NotEmpty(t, result, "Generated UserID should not be empty")

		// 断言结果必须包含资源标识符前缀
		assert.Contains(t, result, rid.UserID.String()+"-", "Generated UserID should start with \"user-\"")

		// 断言前缀不会与 uniqueStr 部分重叠
		splitParts := strings.SplitN(result, "-", 2)
		assert.Equal(t, rid.UserID.String(), splitParts[0], "The prefix of generated UserID should be \"user\"")

		// 断言生成的 ID 具有固定长度
		if len(splitParts) == 2 {
			assert.Equal(t, 6, len(splitParts[1]), "The unique part of generated UserID should be 6 characters long")
		} else {
			t.Errorf("Generated UserID %s does not contain a unique part", result)
		}
	})
}
