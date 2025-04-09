/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-21 23:10:02
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-09 21:12:26
 * @FilePath: \ZesOJ\Disassembly\gdb\debugger_windows_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

import (
	"fmt"
	"syscall"
	"testing"
)

// TestCreateAndBlockProcess 测试 CreateAndBlockProcess 函数。
func TestCreateAndBlockProcess(t *testing.T) {
	tests := []struct {
		exePath string
		cmdLine string
		wantErr bool
	}{
		{"E:\\ZesOJ\\sever\\test.exe", "", false},
		// {"C:\\InvalidPath\\invalid.exe", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.exePath, func(t *testing.T) {
			process, thread, err := CreateAndBlockProcess(tt.exePath, tt.cmdLine)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAndBlockProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if process == 0 || thread == 0 {
					t.Errorf("CreateAndBlockProcess() returned invalid handles: process = %v, thread = %v", process, thread)
				}

				// // Wait for 10 seconds
				// time.Sleep(10 * time.Second)
				// Close handles to avoid resource leaks
				syscall.CloseHandle(process)
				syscall.CloseHandle(thread)
			}
		})
	}
}

func TestReadProcessMemory(t *testing.T) {
	type Test struct {
		exePath string
		cmdLine string
		wantErr bool
	}

	var test Test
	test.exePath = "E:\\ZesOJ\\sever\\test.exe"
	test.cmdLine = ""
	test.wantErr = false

	process, thread, err := CreateAndBlockProcess(test.exePath, test.cmdLine)

	if (err != nil) != test.wantErr {
		t.Errorf("CreateAndBlockProcess() error = %v, wantErr %v", err, test.wantErr)
		return
	}
	if !test.wantErr {
		if process == 0 || thread == 0 {
			t.Errorf("CreateAndBlockProcess() returned invalid handles: process = %v, thread = %v", process, thread)
		}

		t.Log("OK with CreateAndBlockProcess")

		var buffer []byte
		context, err := GetThreadContext(thread)
		if err != nil {
			t.Errorf("GetThreadContext() error = %v", err)
			return
		}
		t.Logf("RIP: %x", context.Rip)
		rip := uintptr(context.Rip) // Example address, replace with actual RIP value

		buffer, err = ReadProcessMemory(process, rip, 8)
		if err != nil {
			t.Errorf("ReadProcessMemory() error = %v", err)
			return
		}

		// 添加反汇编验证
		if len(buffer) > 0 {
			asm, err := Disassemble(buffer, uint64(rip), 64)
			if err == nil {
				t.Logf("内存反汇编结果: %s", asm.Armcode)
			} else {
				t.Errorf("反汇编失败，原始字节: % X, 错误信息: %v", buffer, err)
			}
		}

		t.Logf("Memory at RIP: %x", buffer)
		// Close handles to avoid resource leaks
		syscall.CloseHandle(process)
		syscall.CloseHandle(thread)
	}
}

func TestGetThreadContext(t *testing.T) {
	// 测试用例参数化
	tests := []struct {
		name         string
		exePath      string
		wantValidRip bool
		wantErr      bool
	}{
		{"正常进程上下文", "E:\\ZesOJ\\sever\\test.exe", true, false},
		// {"无效线程句柄", "E:\\ZesOJ\\sever\\test.exe", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试进程
			process, thread, err := CreateAndBlockProcess(tt.exePath, "")
			if err != nil && !tt.wantErr {
				t.Fatalf("CreateAndBlockProcess() 错误 = %v", err)
			}
			defer syscall.CloseHandle(process)
			defer syscall.CloseHandle(thread)

			// 测试无效句柄情况
			var invalidThread syscall.Handle = 0
			if tt.name == "无效线程句柄" {
				thread = invalidThread
			}

			// 获取线程上下文
			context, err := GetThreadContext(thread)

			// 错误情况验证
			if (err != nil) != tt.wantErr {
				t.Fatalf("期望错误 = %v, 实际错误 = %v", tt.wantErr, err)
			}

			// 有效情况验证
			if !tt.wantErr {
				// 验证上下文非空
				if context == nil {
					t.Fatal("获取到空的上下文对象")
				}

				t.Logf("线程上下文验证:")

				// 验证关键寄存器值
				registerTests := []struct {
					regName  string
					regValue uint64
				}{
					{"RIP", context.Rip},
					{"RSP", context.Rsp},
					{"RAX", context.Rax},
					{"RBX", context.Rbx},
					{"RCX", context.Rcx},
					{"RDX", context.Rdx},
					{"RSI", context.Rsi},
					{"RDI", context.Rdi},
					{"R8", context.R8},
					{"R9", context.R9},
					{"R10", context.R10},
					{"R11", context.R11},
					{"R12", context.R12},
					{"R13", context.R13},
					{"R14", context.R14},
					{"R15", context.R15},
				}

				for _, rt := range registerTests {
					// if rt.regValue == 0 {
					// 	t.Errorf("%s 寄存器值为0（可能无效）", rt.regName)
					// } else {
					t.Logf("%-5s = 0x%016x", rt.regName, rt.regValue)
					// }
				}

				// // 验证上下文标志位
				// if (context.ContextFlags & 0x10007) != 0x10007 {
				// 	t.Errorf("上下文标志位异常，期望 0x10007 实际 0x%x", context.ContextFlags)
				// }
			}
		})
	}
}

func TestWriteProcessMemory(t *testing.T) {
	tests := []struct {
		name       string
		exePath    string
		data       []byte
		injectAddr uintptr

		setup func() syscall.Handle // 特殊句柄构造
	}{
		// {
		// 	name:    "正常内存写入",
		// 	exePath: "E:\\ZesOJ\\sever\\test.exe",
		// 	data:    []byte{0x90, 0x90, 0x90}, // 3个NOP指令

		// },
		// {
		// 	name:    "无效线程句柄",
		// 	exePath: "E:\\ZesOJ\\sever\\test.exe",
		// 	data:    []byte{0xCC},
		// 	// setup:   func() syscall.Handle { return 0 },
		// },
		{
			name:    "空数据写入",
			exePath: "E:\\ZesOJ\\sever\\test.exe",
			data:    []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试进程
			process, thread, err := CreateAndBlockProcess(tt.exePath, "")

			if err != nil {
				t.Fatalf("CreateAndBlockProcess() 错误 = %v", err)
			}
			defer syscall.CloseHandle(process)
			defer syscall.CloseHandle(thread)
			context, err := GetThreadContext(thread)
			if err != nil {
				t.Fatalf("GetThreadContext() 错误 = %v", err)
			}

			t.Logf("线程上下文验证: rip = 0x%016x", context.Rip)

			ip := (uintptr)(context.Rip)

			buffer, err := ReadProcessMemory(process, ip, 8)

			if err != nil {
				t.Fatalf("ReadProcessMemory() 错误 = %v", err)
			}

			t.Logf("内存 at RIP: %x", buffer)

			ts, err := WriteProcessMemory(process, ip, tt.data)
			if err != nil {
				t.Fatalf("WriteProcessMemory() 错误 = %v", err)
			}
			t.Logf("WriteProcessMemory() 成功，写入字节数 = %d", ts)

			buffer, err = ReadProcessMemory(process, ip, 8)

			if err != nil {
				t.Fatalf("ReadProcessMemory() 错误 = %v", err)
			}

			t.Logf("内存 at RIP: %x", buffer)

		})
	}

}

func TestReviseThreadContext(t *testing.T) {
	// 创建测试进程
	exePath := "E:\\ZesOJ\\sever\\test.exe"
	process, thread, err := CreateAndBlockProcess(exePath, "")
	if err != nil {
		t.Fatalf("创建进程失败: %v", err)
	}
	defer syscall.CloseHandle(process)
	defer syscall.CloseHandle(thread)

	// 获取原始上下文
	originalCtx, err := GetThreadContext(thread)
	if err != nil {
		t.Fatalf("获取原始上下文失败: %v", err)
	}

	tests := []struct {
		name      string
		register  string
		testValue uint64
	}{
		{"修改RAX寄存器", "Rax", 0x1234ABCD},
		{"修改RIP寄存器", "Rip", 0x00007FFBC0012345},
		{"修改RSP寄存器", "Rsp", 0x000000C001235500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制原始上下文进行修改
			modifiedCtx := *originalCtx

			// 执行寄存器修改
			if err := ReviseThreadContext(thread, &modifiedCtx, tt.register, tt.testValue); err != nil {
				t.Fatalf("修改寄存器失败: %v", err)
			}

			// 获取修改后的上下文
			updatedCtx, err := GetThreadContext(thread)
			if err != nil {
				t.Fatalf("获取更新后上下文失败: %v", err)
			}

			// 打印对比结果
			t.Logf("\n[%s 修改前后对比]\n"+
				"修改前 %s = 0x%016X\n"+
				"设置值 %s = 0x%016X\n"+
				"实际值 %s = 0x%016X\n"+
				"修改是否生效: %t",
				tt.register,
				tt.register, getRegisterValue(originalCtx, tt.register),
				tt.register, tt.testValue,
				tt.register, getRegisterValue(updatedCtx, tt.register),
				getRegisterValue(updatedCtx, tt.register) == tt.testValue)
		})
	}
}

// 辅助函数获取寄存器值
func getRegisterValue(ctx *CONTEXT, regName string) uint64 {
	switch regName {
	case "Rax":
		return ctx.Rax
	case "Rip":
		return ctx.Rip
	case "Rsp":
		return ctx.Rsp
	// 可根据需要添加其他寄存器
	default:
		return 0
	}
}

func TestDisassembleRange(t *testing.T) {
	// 创建测试进程
	exePath := "E:\\ZesOJ\\sever\\test.exe"
	process, thread, err := CreateAndBlockProcess(exePath, "")
	if err != nil {
		t.Fatalf("创建进程失败: %v", err)
	}
	defer syscall.CloseHandle(process)
	defer syscall.CloseHandle(thread)

	// 获取原始上下文
	originalCtx, err := GetThreadContext(thread)

	if err != nil {
		t.Fatalf("获取原始上下文失败: %v", err)
	}

	// 获取原始内存内容
	ip := (uintptr)(originalCtx.Rip)
	buffer, err := ReadProcessMemory(process, ip, 800)
	if err != nil {
		t.Fatalf("读取内存失败: %v", err)
	}

	// 反汇编内存
	instructions := DisassembleRange(buffer, (uint64)(ip), 64)

	// 打印反汇编结
	t.Logf("反汇编结果:")
	for _, inst := range *instructions {
		addressStr := fmt.Sprintf("%016X", inst.Address)
		mashionStr := ""
		for i, b := range inst.HexCodes {
			if i > 0 {
				mashionStr += " "
			}
			mashionStr += fmt.Sprintf("%02X", b)
		}
		fmt.Printf("0x%-19s  %-30s  %-30s \n", addressStr, mashionStr, inst.Armcode /*, inst.Length*/)
	}

}

func TestGetProcessID(t *testing.T) {
	tests := []struct {
		name    string
		exePath string
		wantErr bool
	}{
		{"Valid Process", "E:\\ZesOJ\\sever\\test.exe", false},
		// {"Invalid Process Handle", "", true}, // Uncomment to test invalid handle
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test process
			process, thread, err := CreateAndBlockProcess(tt.exePath, "")
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateAndBlockProcess() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer syscall.CloseHandle(process)
			defer syscall.CloseHandle(thread)

			if !tt.wantErr {
				// Get the process ID
				pid, err := GetProcessID(process)
				if err != nil {
					t.Errorf("GetProcessID() error = %v", err)
				} else {
					t.Logf("Process ID: %d", pid)
				}

				// Validate the PID
				if pid == 0 {
					t.Errorf("GetProcessID() returned invalid PID: %d", pid)
				}
			}
		})
	}
}

// func TestBreakpointHandling(t *testing.T) {
// 	// 创建测试进程
// 	exePath := "E:\\ZesOJ\\sever\\test.exe"
// 	process, thread, err := CreateAndBlockProcess(exePath, "")
// 	if err != nil {
// 		t.Fatalf("创建进程失败: %v", err)
// 	}
// 	defer syscall.CloseHandle(process)
// 	defer syscall.CloseHandle(thread)

// 	// 初始化调试机
// 	dbg := &DbgMachine{
// 		process:     process,
// 		thread:      thread,
// 		breakpoints: make(map[uintptr]*Dbgbreak),
// 		textdata:    make(map[uintptr]*Directive),
// 	}

// 	entryPoint := uint64(0x00007FFD745AAF17)

// 	// 读取入口点内存
// 	data, err := ReadProcessMemory(process, uintptr(entryPoint), 16)
// 	if err != nil {
// 		t.Fatalf("读取内存失败: %v", err)
// 	}

// 	// 反汇编入口点指令
// 	directive, err := Disassemble(data, entryPoint, 64)
// 	if err != nil {
// 		t.Fatalf("反汇编失败: %v", err)
// 	}

// 	// 设置断点
// 	t.Logf("设置断点在 0x%X", entryPoint)
// 	if err := dbg.SetBreakpoint(&directive); err != nil {
// 		t.Fatalf("设置断点失败: %v", err)
// 	}

// 	// 调试事件循环
// 	var debugEvent *DEBUG_EVENT
// 	for {
// 		// 等待调试事件
// 		debugEvent, err = WaitForDebug()
// 		if err != nil {
// 			t.Fatalf("等待调试事件失败: %v", err)
// 		}

// 		switch debugEvent.DebugEventCode {
// 		case EXCEPTION_DEBUG_EVENT:
// 			if debugEvent.Exception.ExceptionRecord.ExceptionCode == EXCEPTION_BREAKPOINT {
// 				t.Log("\n=== 断点触发 ===")
// 				ctx, err := GetThreadContext(thread)
// 				if err != nil {
// 					t.Fatalf("获取上下文失败: %v", err)
// 				}

// 				t.Logf("RIP = 0x%016X", ctx.Rip)

// 				// 获取断点信息
// 				bpAddr := uintptr(debugEvent.Exception.ExceptionRecord.ExceptionAddress)
// 				bp, exists := dbg.breakpoints[bpAddr]
// 				if !exists {
// 					t.Fatal("未找到注册的断点信息")
// 				}

// 				// 打印断点详细信息
// 				t.Logf("断点地址: 0x%X", bp.address)
// 				t.Logf("原始指令: % X", bp.rawcode.HexCodes)
// 				t.Logf("指令长度: %d", bp.rawcode.Length)
// 				t.Logf("汇编指令: %s", bp.rawcode.Armcode)

// 				// 验证调试状态
// 				ctx, err = GetThreadContext(thread)
// 				if err != nil {
// 					t.Fatalf("获取上下文失败: %v", err)
// 				}

// 				t.Logf("RIP = 0x%016X", ctx.Rip)
// 				t.Logf("断点地址匹配: %v", ctx.Rip == uint64(bp.address))
// 				// 恢复原始指令

// 				t.Log("恢复原始指令...")
// 				if _, err := WriteProcessMemory(process, bp.address, bp.rawcode.HexCodes); err != nil {
// 					t.Fatalf("恢复指令失败: %v", err)
// 				}
// 				// 在断点处理逻辑中添加以下代码
// 				ctx.Rip = uint64(bp.address) // 将RIP设置为原断点地址

// 				// 更新线程上下文
// 				if err := ReviseThreadContext(thread, ctx, "Rip", ctx.Rip); err != nil {
// 					t.Fatalf("修复RIP失败: %v", err)
// 				}

// 				ctx, err = GetThreadContext(thread)
// 				if err != nil {
// 					t.Fatalf("获取上下文失败: %v", err)
// 				}

// 				t.Logf("RIP = 0x%016X", ctx.Rip)
// 				t.Logf("断点地址匹配: %v", ctx.Rip == uint64(bp.address))

// 				// 继续执行后续代码
// 				if err := ContinueDebugEvent(
// 					debugEvent.ProcessId,
// 					debugEvent.ThreadId,
// 					DBG_CONTINUE,
// 				); err != nil {
// 					t.Fatalf("继续执行失败: %v", err)
// 				}

// 				return // 结束测试
// 			}
// 		case EXIT_PROCESS_DEBUG_EVENT:
// 			t.Fatal("进程意外退出")
// 		}

// 		// 继续执行
// 		if err := ContinueDebugEvent(
// 			debugEvent.ProcessId,
// 			debugEvent.ThreadId,
// 			DBG_CONTINUE,
// 		); err != nil {
// 			t.Fatalf("继续执行失败: %v", err)
// 		}
// 	}
// }

// TestGetThreadID 测试 GetThreadID 函数。
func TestGetThreadID(t *testing.T) {
	// 创建测试进程
	exePath := "E:\\ZesOJ\\sever\\test.exe"
	_, tt, _ := CreateAndBlockProcess(exePath, "")

	// 测试用例参数化
	tests := []struct {
		name   string
		thread syscall.Handle
		// wantErr bool
	}{
		// 这里假设我们有一个有效的线程句柄，实际测试时需要替换为真实有效的句柄
		// {"Valid Thread", validThreadHandle, false},
		{name: "test Thread", thread: tt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tid, err := GetThreadID(tt.thread)
			t.Logf("Thread ID: %d", tid)
			if err != nil {
				t.Errorf("GetThreadID() error = %v", err)
			}

		})
	}
}

// // TestGetProcessID 测试 GetProcessID 函数
// func TestGetProcessID(t *testing.T) {
// 	// 创建一个测试进程
// 	exePath := "E:\\ZesOJ\\sever\\test.exe"
// 	process, _, err := CreateAndBlockProcess(exePath, "")
// 	if err != nil {
// 		t.Fatalf("CreateAndBlockProcess() error = %v", err)
// 	}
// 	defer syscall.CloseHandle(process)

// 	tests := []struct {
// 		name    string
// 		process syscall.Handle
// 		wantErr bool
// 	}{
// 		{"Valid Process", process, false},
// 		{"Invalid Process Handle", 0, true},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			pid, err := GetProcessID(tt.process)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetProcessID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !tt.wantErr {
// 				if pid == 0 {
// 					t.Errorf("GetProcessID() returned invalid PID: %d", pid)
// 				} else {
// 					t.Logf("Process ID: %d", pid)
// 				}
// 			}
// 		})
// 	}
// }
