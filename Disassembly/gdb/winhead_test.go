/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-03 21:03:27
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-03 23:49:21
 * @FilePath: \ZesOJ\Disassembly\gdb\winhead_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package gdb

import (
	"reflect"
	"testing"
	"unsafe"
)

// func TestContextSize(t *testing.T) {
// 	goSize := unsafe.Sizeof(CONTEXT{})
// 	cSize := C.get_context_size()

// 	if goSize != uintptr(cSize) {
// 		t.Errorf("结构体大小不匹配: Go %d vs C %d", goSize, cSize)
// 	}
// }

// func TestFieldOffsets(t *testing.T) {
// 	if C.verify_offsets() != 1 {
// 		t.Fatal("关键字段偏移不匹配")
// 	}
// }

func TestUnionAccess(t *testing.T) {
	ctx := CONTEXT{}

	// 测试联合体访问
	ctx.unionArea[0] = 0xAA
	ctx.unionArea[1] = 0xBB

	if ctx.FltSave().ControlWord != 0xBBAA {
		t.Error("联合体访问失败 (FltSave)")
	}

	if ctx.D()[0] != 0xBBAA {
		t.Error("联合体访问失败 (D)")
	}
}

func TestOffsetOfFieldByName(t *testing.T) {
	// 创建CONTEXT实例用于计算期望值
	var context CONTEXT

	tests := []struct {
		name     string
		field    string
		expected uintptr
	}{
		{
			name:     "测试Rax寄存器偏移",
			field:    "Rax",
			expected: unsafe.Offsetof(context.Rax),
		},
		{
			name:     "测试Rip寄存器偏移",
			field:    "Rip",
			expected: unsafe.Offsetof(context.Rip),
		},
		{
			name:     "测试SegCs段寄存器偏移",
			field:    "SegCs",
			expected: unsafe.Offsetof(context.SegCs),
		},
		{
			name:     "测试P1Home",
			field:    "P1Home",
			expected: 0,
		},
		// {
		// 	name:     "测试联合体字段",
		// 	field:    "unionArea",
		// 	expected: unsafe.Offsetof(context.unionArea),
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OffsetOfFieldByName(tt.field)

			if err != nil {
				t.Fatalf("获取字段偏移量时出错: %v", err)
			}

			// 有效情况验证
			if got != tt.expected {
				t.Errorf("字段 %s 偏移量不匹配\n期望: 0x%X\n实际: 0x%X",
					tt.field, tt.expected, got)
			}

			// 验证反射结果与实际内存布局一致
			field, _ := reflect.TypeOf(context).FieldByName(tt.field)
			if got != field.Offset {
				t.Errorf("反射偏移量不一致\n反射: 0x%X\n计算: 0x%X",
					field.Offset, got)
			}

			t.Logf("字段 %s 偏移量: 0x%X", tt.field, got)

		})
	}
}
