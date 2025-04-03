/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-21 23:10:02
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-04 00:58:53
 * @FilePath: \ZesOJ\Disassembly\gdb\debugger_windows_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

import (
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
		rip := uintptr(0x000007FFC7C291070) // Example address, replace with actual RIP value

		buffer, err = ReadProcessMemory(process, rip, 8)
		if err != nil {
			t.Errorf("ReadProcessMemory() error = %v", err)
			return
		}

		t.Logf("Memory at RIP: %x", buffer)

		// Close handles to avoid resource leaks
		syscall.CloseHandle(process)
		syscall.CloseHandle(thread)
	}
}

// ... 已有代码 ...

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

// ... 其他测试函数 ...
// ... 其他测试函数 ...

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
