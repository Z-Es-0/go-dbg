/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-03 21:03:27
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-03 21:40:30
 * @FilePath: \ZesOJ\Disassembly\gdb\winhead_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package gdb

import (
	"testing"
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
