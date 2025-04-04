/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-21 21:49:22
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-04 13:23:18
 * @FilePath: \ZesOJ\Disassembly\gdb\debugger_windows.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/arch/x86/x86asm"
)

// 定义必要的常量
const (
	// 当使用 `CreateProcess` 函数创建一个新进程时，如果在 `dwCreationFlags` 参数中指定了 `DEBUG_PROCESS` 标志，那么新创建的进程会被置于调试模式下。
	// 这意味着父进程（调试器）会接收到新进程产生的所有调试事件，例如进程创建、线程创建、异常等。
	// 调试器可以通过处理这些事件来监控和控制被调试进程的执行。
	DEBUG_PROCESS = 0x00000001
)

// from kernel32.dll
var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procCreateProcess      = modkernel32.NewProc("CreateProcessW")     // 1. CreateProcessW: 用于创建一个新的进程及其主线程。它允许你指定要执行的程序、命令行参数、环境变量等信息，并且可以控制新进程的创建方式，如是否以调试模式启动等。
	procWaitForDebugEvent  = modkernel32.NewProc("WaitForDebugEvent")  // 2. WaitForDebugEvent: 用于等待被调试进程产生的调试事件。当一个进程以调试模式启动（如使用 DEBUG_PROCESS 标志）时，调试器可以使用此函数来接收诸如进程创建、线程创建、异常等调试事件。
	procContinueDebugEvent = modkernel32.NewProc("ContinueDebugEvent") // 3. ContinueDebugEvent: 用于通知操作系统继续执行被调试的进程。当调试器处理完一个调试事件后，需要调用此函数来让被调试进程继续执行。
	procGetThreadContext   = modkernel32.NewProc("GetThreadContext")   // 4. GetThreadContext: 用于获取指定线程的上下文信息，包括寄存器的值、线程状态等。这对于调试器来说非常有用，可以帮助调试器了解线程的执行状态。
	procSetThreadContext   = modkernel32.NewProc("SetThreadContext")   // 5. SetThreadContext: 用于设置指定线程的上下文信息。调试器可以使用此函数来修改线程的寄存器值或其他状态信息，从而控制线程的执行。
	procReadProcessMemory  = modkernel32.NewProc("ReadProcessMemory")  // 6. ReadProcessMemory: 用于从指定进程的内存中读取数据。调试器可以使用此函数来查看被调试进程的内存内容，例如查看变量的值。
	procWriteProcessMemory = modkernel32.NewProc("WriteProcessMemory") // 7. WriteProcessMemory: 用于向指定进程的内存中写入数据。调试器可以使用此函数来修改被调试进程的内存内容，例如修改变量的值。
	procDebugActiveProcess = modkernel32.NewProc("DebugActiveProcess") // 8. DebugActiveProcess: 用于将当前进程附加到一个正在运行的进程上进行调试。通过此函数，调试器可以开始监控和控制一个已经存在的进程的执行。
)

// 创建一个新的子进程并以调试模式启动，这样可以阻塞它直到调试器允许其继续执行
// 返回值：
//   - 进程的句柄
//   - 线程的句柄
//   - 错误信息，如果创建进程失败
func CreateAndBlockProcess(exePath string, cmdLine string) (syscall.Handle, syscall.Handle, error) {
	var (
		si            syscall.StartupInfo
		pi            syscall.ProcessInformation
		exePathPtr, _ = syscall.UTF16PtrFromString(exePath)
		cmdLinePtr, _ = syscall.UTF16PtrFromString(cmdLine)
	)

	// 调用 CreateProcessW 函数创建新进程
	_, _, err := procCreateProcess.Call(
		uintptr(unsafe.Pointer(exePathPtr)), // 指向可执行文件名称的指针，传入可执行文件路径指针
		uintptr(unsafe.Pointer(cmdLinePtr)), // 指向命令行字符串的指针
		uintptr(0),                          // 指向进程安全属性的指针
		uintptr(0),                          // 指向线程安全属性的指针
		uintptr(0),                          // 指示新进程是否继承调用进程的句柄
		uintptr(DEBUG_PROCESS),              // 进程创建标志，使用DEBUG_PROCESS标志以调试模式启动进程
		uintptr(0),                          // 指向环境变量块的指针
		uintptr(0),                          // 指向当前目录名称的指针
		uintptr(unsafe.Pointer(&si)),        // 指向STARTUPINFO结构体的指针，用于指定新进程的窗口和控制台属性
		uintptr(unsafe.Pointer(&pi)),        // 指向PROCESS_INFORMATION结构体的指针，用于接收新进程和主线程的句柄和ID
	)
	if err != nil && err != syscall.Errno(0) {
		return 0, 0, err
	}

	return pi.Process, pi.Thread, nil
}

// ReadProcessMemory 函数用于从指定进程的内存中读取数据。
// 参数:
//   - process: 要读取内存的进程的句柄。
//   - address: 要读取的内存地址。
//   - size: 要读取的字节数。
//
// 返回值:
//   - 读取到的数据字节切片。
//   - 错误信息，如果读取过程中出现错误。
func ReadProcessMemory(process syscall.Handle, address uintptr, size uint) ([]byte, error) {
	// 创建一个指定大小的字节切片，用于存储读取到的数据
	buf := make([]byte, size)
	// 用于存储实际读取的字节数
	var bytesRead uint32

	// 调用Windows API函数 ReadProcessMemory 从指定进程的内存中读取数据
	ret, _, err := procReadProcessMemory.Call(
		uintptr(process),                    // 进程句柄
		address,                             // 要读取的内存地址
		uintptr(unsafe.Pointer(&buf[0])),    // 存储读取数据的缓冲区指针
		uintptr(size),                       // 要读取的字节数
		uintptr(unsafe.Pointer(&bytesRead)), // 存储实际读取字节数的变量指针
	)

	// 如果调用失败（返回值为0），则返回错误
	if ret == 0 {
		return nil, err
	}
	// 返回实际读取的数据
	return buf[:bytesRead], nil

}

// GetThreadContext 函数用于获取指定线程的上下文信息。
// 参数:
//   - thread: 要获取上下文的线程句柄。
//
// 返回值:
//   - CONTEXT 结构体，包含线程的上下文信息。
//   - 错误信息，如果获取上下文失败。
func GetThreadContext(thread syscall.Handle) (*CONTEXT, error) {
	// 创建一个 CONTEXT 结构体实例
	var context CONTEXT
	// 设置 ContextFlags，指定要获取的上下文信息
	context.ContextFlags = CONTEXT_FULL

	// 调用 Windows API 函数 GetThreadContext 获取线程上下文
	ret, _, err := procGetThreadContext.Call(
		uintptr(thread),                   // 线程句柄
		uintptr(unsafe.Pointer(&context)), // 指向 CONTEXT 结构体的指针
	)

	// 如果调用失败（返回值为 0），则返回错误
	if ret == 0 {
		return nil, err
	}

	// 返回获取到的线程上下文
	return &context, nil
}

// WriteProcessMemory 函数用于向指定线程的内存中写入数据。
// 参数:
//   - thread: 要写入内存的线程句柄。
//   - address: 要写入的内存地址。
//   - data: 要写入的数据字节切片。
//
// 返回值:
//   - 写入的字节数。
//   - 错误信息，如果写入过程中出现错误。
func WriteProcessMemory(thread syscall.Handle, address uintptr, data []byte) (uintptr, error) {
	if (data == nil) || (len(data) == 0) {
		return 0, nil
	}
	var bytesWritten uint32
	ret, _, err := procWriteProcessMemory.Call(
		uintptr(thread),
		address,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&bytesWritten)),
	)
	if ret == 0 {
		return 0, err
	}
	return uintptr(bytesWritten), nil
}

// Makebreakpoint 函数用于在指定线程的指定地址设置断点。
// 参数:
//   - thread: 要设置断点的线程句柄。
//   - address: 要设置断点的内存地址。
//
// 返回值:
//   - 错误信息，如果设置断点失败。
func Makebreakpoint(thread syscall.Handle, address uintptr) error {

	breakpointInstruction := byte(0xCC)

	// 写入断点指令到指定地址
	_, err := WriteProcessMemory(thread, address, []byte{breakpointInstruction})
	if err != nil {
		return err
	}

	return nil
}

// setFlag 函数用于设置或清除指定的标志位。
// 参数:
//   - EFlags: 原始的标志位值。
//   - flag: 要设置或清除的标志位。
//   - set: 一个布尔值，如果为 true，则设置标志位；如果为 false，则清除标志位。
//
// 返回值:
//   - 更新后的标志位值。
func setFlag(EFlags uint32, flag uint32, set bool) uint32 {
	if set {
		return EFlags | flag
	}
	return EFlags & ^flag
}

// ReviseThreadContext 函数用于设置指定线程的上下文信息。
// 参数:
//   - thread: 要设置上下文的线程句柄。
//   - context: 要设置的上下文信息。
//   - name: 寄存器名称，目标修改的寄存器名称。
//   - value: 寄存器值，要设置的寄存器值。
//
// 返回值:
//   - 错误信息，如果设置上下文失败。
func ReviseThreadContext(thread syscall.Handle, context *CONTEXT, name string, value uint64) error {
	switch name {
	case "Rip":
		context.Rip = value
	case "Rax":
		context.Rax = value
	case "Rcx":
		context.Rcx = value
	case "Rdx":
		context.Rdx = value
	case "Rbx":
		context.Rbx = value
	case "Rsp":
		context.Rsp = value
	case "Rbp":
		context.Rbp = value
	case "Rsi":
		context.Rsi = value
	case "Rdi":
		context.Rdi = value
	case "R8":
		context.R8 = value
	case "R9":
		context.R9 = value
	case "R10":
		context.R10 = value
	case "R11":
		context.R11 = value
	case "R12":
		context.R12 = value
	case "R13":
		context.R13 = value
	case "R14":
		context.R14 = value
	case "R15":
		context.R15 = value
	case "EFlags":
		context.EFlags = ReviseEFlags(context.EFlags, value, name)
	default:
		return fmt.Errorf("unsupported register name: %s", name)
	}

	// 调用 SetThreadContext 更新线程上下文
	ret, _, err := procSetThreadContext.Call(
		uintptr(thread),
		uintptr(unsafe.Pointer(context)),
	)
	if ret == 0 {
		return err
	}

	return nil
}

func ReviseEFlags(EFlags uint32, value uint64, name string) uint32 {
	set := value != 0
	switch name {
	case "CF":
		return setFlag(EFlags, CF, set)
	case "PF":
		return setFlag(EFlags, PF, set)
	case "AF":
		return setFlag(EFlags, AF, set)
	case "ZF":
		return setFlag(EFlags, ZF, set)
	case "SF":
		return setFlag(EFlags, SF, set)
	case "TF":
		return setFlag(EFlags, TF, set)
	case "IF":
		return setFlag(EFlags, IF, set)
	case "DF":
		return setFlag(EFlags, DF, set)
	case "OF":
		return setFlag(EFlags, OF, set)
	default:
		return EFlags
	}
}

// Disassemble 反汇编单个指令
// 参数:
//   - code: 机器码字节切片
//   - pc: 指令地址（用于计算跳转目标）
//   - bit: 架构位数 (32 或 64)
//
// 返回值:
//   - 汇编指令字符串
//   - 指令长度
//   - 错误信息
func Disassemble(code []byte, pc uint64, bit int) (string, int, error) {
	inst, err := x86asm.Decode(code, bit)
	if err != nil {
		return "", 0, fmt.Errorf("解码失败: %w", err)
	}

	// 获取指令文本
	syntax := x86asm.GNUSyntax(inst, pc, nil)
	return fmt.Sprintf("%-36s // % x", syntax, code[:inst.Len]), inst.Len, nil

}

// DisassembleRange 反汇编指定内存范围
// 参数:
//   - memory: 要反汇编的内存字节切片。
//   - startAddr: 起始地址，用于计算指令的实际地址。
//   - bit: 架构位数 (32 或 64)。
//
// 返回值:
//   - 包含反汇编结果的字符串切片，每个元素对应一条指令。
func DisassembleRange(memory []byte, startAddr uint64, bit int) []string {
	var (
		offset    uint64   // 当前偏移量，用于遍历内存字节切片。
		result    []string // 存储反汇编结果的字符串切片。
		maxLength = 15     // 单条指令最大尝试长度，通常为 8 字节。
	)

	// 遍历内存字节切片，逐条反汇编指令
	for offset < uint64(len(memory)) {
		remaining := len(memory) - int(offset) // 剩余未处理的字节数
		if remaining <= 0 {
			break // 如果没有剩余字节，退出循环
		}

		// 动态调整解码长度，确保不会超出剩余字节范围
		tryLength := remaining
		if tryLength > maxLength {
			tryLength = maxLength
		}

		// 从当前偏移量开始获取待解码的字节切片
		code := memory[offset:]
		inst, err := x86asm.Decode(code, bit) // 尝试解码指令
		if err != nil {
			// 如果解码失败，假设当前字节为无效指令，并向前推进1字节
			result = append(result, fmt.Sprintf("DB 0x%02x      ; 无效指令", code[0]))
			offset++ // 偏移量增加1
			continue
		}

		// 成功解码指令后，调用 Disassemble 获取指令的汇编文本
		syntax, _, err := Disassemble(code, startAddr+offset, bit)
		if err != nil {
			// 如果反汇编过程中出现错误，记录错误信息并跳过该指令
			result = append(result, fmt.Sprintf("Error: %v", err))
			offset++ // 偏移量增加1
			continue
		}

		// 将反汇编结果添加到结果切片中
		result = append(result, syntax)
		offset += uint64(inst.Len) // 偏移量增加当前指令的长度
	}

	// 返回包含所有反汇编结果的字符串切片
	return result
}
