/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-21 23:10:02
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-25 14:20:48
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
