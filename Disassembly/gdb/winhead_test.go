/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-03 21:03:27
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-09 22:27:01
 * @FilePath: \ZesOJ\Disassembly\gdb\winhead_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package gdb

import (
	"fmt"

	"reflect"
	"syscall"
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

func TestGetUnion(t *testing.T) {

}

// 新增断点工作流测试
// 新增断点工作流测试
func TestBreakpointWorkflow(t *testing.T) {
	exePath := "E:\\ZesOJ\\sever\\test.exe"
	process, thread, err := CreateAndBlockProcess(exePath, "")
	if err != nil {
		t.Fatalf("进程创建失败: %v", err)
	}
	defer syscall.CloseHandle(process)
	defer syscall.CloseHandle(thread)

	debugmashin := &DbgMachine{
		process:     process,
		thread:      thread,
		breakpoints: make(map[uintptr]*Dbgbreak),
		textdata:    make(map[uintptr]*Directive),
	}

	// 设置断点
	err = debugmashin.Maketextdata()
	if err != nil {
		t.Fatalf("设置断点失败: %v", err)
	}

	threadID, err := GetThreadID(thread)
	if err != nil {
		t.Fatalf("获取线程ID失败: %v", err)
	}
	processID, err := GetProcessID(process)
	if err != nil {
		t.Fatalf("获取进程ID失败: %v", err)
	}

	debugEvent := &DEBUG_EVENT{
		DebugEventCode: 0,
		ThreadId:       threadID,
		ProcessId:      processID,
	}
	if err = debugmashin.SetBreakpoint(0x00007FFFD102AF1E); err != nil {
		t.Fatalf("设置断点失败: %v", err)
	}

	// 调试事件循环
	for {
		debugEvent, err = WaitForDebug(debugEvent)
		if err != nil {
			t.Fatalf("等待调试事件失败: %v", err)
		}

		switch debugEvent.DebugEventCode {

		case EXCEPTION_DEBUG_EVENT:
			{

				switch (GetUnion[EXCEPTION_DEBUG_INFO](debugEvent)).ExceptionRecord.ExceptionCode {

				case EXCEPTION_ACCESS_VIOLATION:
					fmt.Println("内存访问冲突")
				case EXCEPTION_BREAKPOINT:
					fmt.Println("断点触发")
					DoEXCEPTION_BREAKPOINT(debugmashin, debugEvent)

					debugmashin.DeleteBreakpoint(0x00007FFFD102AF1E)

				case EXCEPTION_SINGLE_STEP:
					fmt.Println("单步执行异常") // 单步执行异常
					DoEXCEPTION_BREAKPOINT(debugmashin, debugEvent)
				case EXCEPTION_GUARD_PAGE:
					fmt.Println("保护页异常") // 保护页异常

				case EXCEPTION_DATATYPE_MISALIGNMENT:
					fmt.Println("数据类型不匹配异常") // 数据类型不匹配异常

				case EXCEPTION_NONCONTINUABLE_EXCEPTION:
					fmt.Println("不可继续执行异常") // 不可继续执行异常

				default:
					fmt.Println((GetUnion[EXCEPTION_DEBUG_INFO](debugEvent)).ExceptionRecord.ExceptionCode)

					fmt.Println("其他异常")

				}
			}

		case CREATE_THREAD_DEBUG_EVENT:
			fmt.Println("线程创建")
			context, _ := GetThreadContext(debugmashin.thread)
			PrintContext(context)
		case LOAD_DLL_DEBUG_EVENT:
			fmt.Println("DLL加载")
			context, _ := GetThreadContext(debugmashin.thread)
			PrintContext(context)

		case EXIT_PROCESS_DEBUG_EVENT:
			t.Log("目标进程正常退出")
			return

		}

		ContinueDebugEvent(
			debugEvent.ProcessId,
			debugEvent.ThreadId,
			DBG_CONTINUE,
		)

	}
}

func DoEXCEPTION_BREAKPOINT(debugmashin *DbgMachine, DebugEv *DEBUG_EVENT) {
	// 处理断点事件
	fmt.Println("等会写")
}

func PrintContext(ctx *CONTEXT) {
	fmt.Printf("Rip: 0x%X\n", ctx.Rip)
	fmt.Printf("Rax: 0x%X\n", ctx.Rax)
	fmt.Printf("Rcx: 0x%X\n", ctx.Rcx)
	fmt.Printf("Rdx: 0x%X\n", ctx.Rdx)
	fmt.Printf("Rbx: 0x%X\n", ctx.Rbx)
	fmt.Printf("Rsp: 0x%X\n", ctx.Rsp)
	fmt.Printf("Rbp: 0x%X\n", ctx.Rbp)
	fmt.Printf("Rsi: 0x%X\n", ctx.Rsi)
	fmt.Printf("Rdi: 0x%X\n", ctx.Rdi)
}
