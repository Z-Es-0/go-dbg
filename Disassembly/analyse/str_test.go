/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-18 03:09:20
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-18 03:38:19
 * @FilePath: \ZesOJ\Disassembly\analyse\str_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AEs
 */
package analyse

import (
	"fmt"
	"testing"
)

func TestStr(t *testing.T) {
	strings, err := FindHardcodedStrings("E:\\ZesOJ\\sever\\test.exe")
	if err != nil {
		panic(err)
	}

	for _, s := range strings {
		fmt.Println(s.Addr, s.Segment, s.Str)
	}
}

// ... 原有测试代码保持不变 ...

func TestSelectStr(t *testing.T) {
	// 准备测试数据
	testData := []StrDate{
		{Addr: 0x1000, Str: "Hello World", Segment: ".rdata"},
		{Addr: 0x2000, Str: "SECRET_KEY_123", Segment: ".data"},
		{Addr: 0x3000, Str: "version:1.0", Segment: ".rdata"},
		{Addr: 0x4000, Str: "http://api.example.com", Segment: ".text"},
	}

	testCases := []struct {
		name        string
		pattern     string
		expected    []string
		expectError bool
	}{
		{
			name:     "基础匹配",
			pattern:  `Hello`,
			expected: []string{"Hello World"},
		},
		{
			name:     "数字匹配",
			pattern:  `\d+`,
			expected: []string{"SECRET_KEY_123", "version:1.0"},
		},
		{
			name:     "URL匹配",
			pattern:  `http://.*`,
			expected: []string{"http://api.example.com"},
		},
		{
			name:     "大小写不敏感",
			pattern:  `(?i)SECRET`,
			expected: []string{"SECRET_KEY_123"},
		},
		{
			name:        "错误正则",
			pattern:     "[invalid",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SelectStr(&testData, tc.pattern)

			if tc.expectError {
				if err == nil {
					t.Error("预期错误但未返回错误")
				}
				return
			}

			if err != nil {
				t.Errorf("意外错误: %v", err)
				return
			}

			if len(result) != len(tc.expected) {
				t.Errorf("预期匹配%d项，实际匹配%d项", len(tc.expected), len(result))
			}

			for i, item := range result {
				if item.Str != tc.expected[i] {
					t.Errorf("第%d项不匹配，预期: %s，实际: %s", i+1, tc.expected[i], item.Str)
				}
			}
		})
	}
}
