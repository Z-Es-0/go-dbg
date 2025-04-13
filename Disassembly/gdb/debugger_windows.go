/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-21 21:49:22
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-13 20:56:41
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
	DEBUG_PROCESS                           = 0x00000001
	EXCEPTION_GUARD_PAGE                    = 0x80000001 //保护页异常
	EXCEPTION_IN_PAGE_ERROR                 = 0x80000006 //页错误异常
	EXCEPTION_ILLEGAL_INSTRUCTION           = 0xC000008C //非法指令异常
	EXCEPTION_PRIV_INSTRUCTION              = 0xC0000096 //特权指令异常
	EXCEPTION_BREAKPOINT_WITH_CODE_CHANGE   = 0x80010004 //代码改变的断点异常
	EXCEPTION_DATATYPE_MISALIGNMENT         = 0x80000002 //数据类型对齐异常
	EXCEPTION_NONCONTINUABLE_EXCEPTION      = 0xC000000D //不可继续的异常
	EXCEPTION_STACK_OVERFLOW                = 0xC00000FD //堆栈溢出异常
	EXCEPTION_INVALID_DISPOSITION           = 0xC00000FE //无效处理异常
	DBG_CONTINUE                            = 0x00010002 //正常继续执行
	DBG_EXCEPTION_NOT_HANDLED               = 0x80010001 //异常未处理
	EXCEPTION_BREAKPOINT                    = 0x80000003 //断点异常
	EXCEPTION_SINGLE_STEP                   = 0x80000004 //单步异常
	EXCEPTION_ACCESS_VIOLATION              = 0xC0000005 //访问冲突异常
	EXCEPTION_INT_DIVIDE_BY_ZERO            = 0xC0000094 //整数除以零异常
	EXCEPTION_INT_OVERFLOW                  = 0xC0000095 //整数溢出异常
	EXCEPTION_INT_INVALID_OPERATION         = 0xC0000096 //无效操作异常
	EXCEPTION_INT_NO_MAPPING                = 0xC0000097 //无映射异常
	EXCEPTION_INT_PRIV_INSTRUCTION          = 0xC0000098 //特权指令异常
	EXCEPTION_INT_INVALID_STATE             = 0xC0000099 //无效状态异常
	EXCEPTION_INT_NONCONTINUABLE_EXCEPTION  = 0xC000009A //不可继续的异常
	EXCEPTION_INT_STACK_OVERFLOW            = 0xC00000FD //堆栈溢出异常
	EXCEPTION_INT_INVALID_DISPOSITION       = 0xC00000FE //无效处理异常
	EXCEPTION_INT_ARRAY_BOUNDS_EXCEEDED     = 0xC00000FF //数组越界异常
	EXCEPTION_INT_FLT_DENORMAL_OPERAND      = 0xC0000100 //浮点数异常
	EXCEPTION_INT_FLT_DIVIDE_BY_ZERO        = 0xC0000101 //浮点数除以零异常
	EXCEPTION_INT_FLT_INEXACT_RESULT        = 0xC0000102 //浮点数结果不精确异常
	EXCEPTION_INT_FLT_INVALID_OPERATION     = 0xC0000103 //浮点数无效操作异常
	EXCEPTION_INT_FLT_OVERFLOW              = 0xC0000104 //浮点数溢出异常
	EXCEPTION_INT_FLT_STACK_CHECK           = 0xC0000105 //浮点数堆栈检查异常
	EXCEPTION_INT_FLT_UNDERFLOW             = 0xC0000106 //浮点数下溢异常
	EXCEPTION_INT_FLT_PRECISION             = 0xC0000107 //浮点数精度异常
	EXCEPTION_INT_ILLEGAL_INSTRUCTION       = 0xC0000108 //非法指令异常
	EXCEPTION_INT_DATATYPE_MISALIGN         = 0xC000010A //数据类型对齐异常
	EXCEPTION_INT_MACHINE_CHECK             = 0xC000010B //机器检查异常
	EXCEPTION_INT_SIMD_FP_EXCEPTION         = 0xC000010C //SIMD浮点异常
	EXCEPTION_INT_PRIV_INSTRUCTION2         = 0xC000010D //特权指令异常2
	EXCEPTION_INT_INVALID_HANDLE            = 0xC000010E //无效句柄异常
	EXCEPTION_INT_INVALID_PARAMETER         = 0xC000010F //无效参数异常
	EXCEPTION_INT_ATTEMPTED_SEGMENT_OVERRUN = 0xC0000110 //段越界异常
	EXCEPTION_INT_ATTEMPTED_STACK_OVERRUN   = 0xC0000111 //堆栈越界异常
	EXCEPTION_INT_INVALID_HANDLE2           = 0xC0000112 //无效句柄异常2
	EXCEPTION_INT_INVALID_HANDLE3           = 0xC0000113 //无效句柄异常3
	EXCEPTION_INT_INVALID_HANDLE4           = 0xC0000114 //无效句柄异常4
	EXCEPTION_INT_INVALID_HANDLE5           = 0xC0000115 //无效句柄异常5
	EXCEPTION_INT_INVALID_HANDLE6           = 0xC0000116 //无效句柄异常6
	EXCEPTION_INT_INVALID_HANDLE7           = 0xC0000117 //无效句柄异常7
	EXCEPTION_INT_INVALID_HANDLE8           = 0xC0000118 //无效句柄异常8
	EXCEPTION_INT_INVALID_HANDLE9           = 0xC0000119 //无效句柄异常9
	EXCEPTION_INT_INVALID_HANDLE10          = 0xC000011A //无效句柄异常10
	EXCEPTION_INT_INVALID_HANDLE11          = 0xC000011B //无效句柄异常11

	EXCEPTION_INT3         = 0xCC       //INT3指令
	EXCEPTION_BREAKPOINT2  = 0x80000002 //断点异常2
	EXCEPTION_SINGLE_STEP2 = 0x80000006 //单步异常2
	EXCEPTION_BREAKPOINT3  = 0x80000007 //断点异常3
	EXCEPTION_SINGLE_STEP3 = 0x80000008 //单步异常3
	EXCEPTION_BREAKPOINT4  = 0x80000009 //断点异常4
	EXCEPTION_SINGLE_STEP4 = 0x8000000A //单步异常4
	EXCEPTION_BREAKPOINT5  = 0x8000000B //断点异常5
	EXCEPTION_SINGLE_STEP5 = 0x8000000C //单步异常5
)

// from kernel32.dll
var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetProcessId       = modkernel32.NewProc("GetProcessId")       // GetProcessId: 用于获取指定进程的 PID。
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

// Disassemble 函数用于将单个机器码指令反汇编为汇编指令
// 参数:
//   - code: 包含机器码的字节切片，代表要反汇编的指令
//   - pc: 指令的地址，用于计算跳转目标等相对地址
//   - bit: 架构的位数，可取值为 32 或 64，用于指定指令的解码模式
//
// 返回值:
//   - Directive 结构体，包含反汇编后的指令信息，如汇编指令字符串、指令长度、指令地址等
//   - 错误信息，如果在解码或反汇编过程中出现错误，则返回相应的错误对象
func Disassemble(code []byte, pc uint64, bit int) (Directive, error) {
	// 尝试对输入的机器码进行解码
	inst, err := x86asm.Decode(code, bit)
	// 初始化一个 Directive 结构体，用于存储反汇编结果
	res := Directive{}
	// 若解码过程中出现错误，则返回带有错误信息的空结果
	if err != nil {
		return res, fmt.Errorf("解码失败: %w", err)
	}

	// 使用 GNU 语法将解码后的指令转换为汇编指令字符串
	syntax := x86asm.GNUSyntax(inst, pc, nil)
	// 设置指令长度
	res.Length = uint32(inst.Len)
	// 设置指令地址
	res.Address = pc
	// 设置汇编指令字符串
	res.Armcode = syntax
	// 设置机器码的十六进制表示
	res.HexCodes = code[:inst.Len]

	// 返回包含反汇编结果的结构体和 nil 错误
	return res, nil
}

// DisassembleRange 反汇编指定内存范围
// 参数:
//   - memory: 要反汇编的内存字节切片。
//   - startAddr: 起始地址，用于计算指令的实际地址。
//   - bit: 架构位数 (32 或 64)。
//
// 返回值:
//   - 包含反汇编结果的 Directive 结构体切片的地址，每个元素对应一条指令。
func DisassembleRange(memory []byte, startAddr uint64, bit int) *[]Directive {
	var (
		offset    uint64      // 当前偏移量，用于遍历内存字节切片。
		result    []Directive // 存储反汇编结果的 Directive 结构体切片。
		maxLength = 15        // 单条指令最大尝试长度
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
			// 如果解码失败，假设当前字节为无效指令
			directive := Directive{
				Length:   1,
				Address:  startAddr + offset,
				Armcode:  fmt.Sprintf("DB 0x%02x      ; 无效指令", code[0]),
				HexCodes: []byte{code[0]},
			}
			result = append(result, directive)
			offset++ // 偏移量增加1
			continue
		}

		// 成功解码指令后，调用 Disassemble 获取指令的汇编文本
		directive, err := Disassemble(code, startAddr+offset, bit)
		if err != nil {
			// 如果反汇编过程中出现错误，记录错误信息
			directive = Directive{
				Length:   1,
				Address:  startAddr + offset,
				Armcode:  fmt.Sprintf("Error: %v", err),
				HexCodes: []byte{code[0]},
			}
			result = append(result, directive)
			offset++ // 偏移量增加1
			continue
		}

		// 将反汇编结果添加到结果切片中
		result = append(result, directive)
		offset += uint64(inst.Len) // 偏移量增加当前指令的长度
	}

	// 返回包含所有反汇编结果的 Directive 结构体切片
	return &result

}

// GetProcessID 函数用于获取子进程的 PID。
// 参数:
//   - process: 子进程的句柄。
//
// 返回值:
//   - 子进程的 PID。
//   - 错误信息，如果获取 PID 失败。
func GetProcessID(process syscall.Handle) (uint32, error) {
	var pid uint32

	// 调用 GetProcessId 函数获取子进程的 PID
	ret, _, err := syscall.Syscall(procGetProcessId.Addr(), 1, uintptr(process), 0, 0)
	if ret == 0 {
		return 0, err
	}

	pid = uint32(ret)
	return pid, nil
}

// 修复后的 GetThreadID 函数，用于获取指定线程的线程ID。
// 参数:
//   - thread: 要获取线程ID的线程句柄。
//
// 返回值:
//   - 线程的ID。
//   - 错误信息，如果获取线程ID失败。
func GetThreadID(thread syscall.Handle) (uint32, error) {
	// 定义一个变量用于存储线程ID
	var tid uint32

	// 引入 kernel32.dll 中的 GetThreadId 函数
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetThreadId := modkernel32.NewProc("GetThreadId")

	// 调用 Syscall 函数来执行 GetThreadId 系统调用，获取线程ID
	ret, _, err := syscall.Syscall(procGetThreadId.Addr(), 1, uintptr(thread), 0, 0)

	// 如果返回值为 0，表示调用失败，返回错误信息
	if ret == 0 {
		return 0, err
	}

	// 将返回值转换为 uint32 类型并赋值给 tid
	tid = uint32(ret)
	// 返回线程ID和 nil 错误信息
	return tid, nil
}

// ContinueDebugEvent 函数用于通知操作系统继续执行被调试的进程。
// 参数:
//   - processId: 被调试进程的 PID。
//   - threadId: 被调试线程的 ID。
//   - continueStatus: 继续执行的状态，通常为 DBG_CONTINUE。
//
// 返回值:
//   - 错误信息，如果继续执行过程中出现错误。
func ContinueDebugEvent(processId, threadId uint32, continueStatus uint32) error {
	ret, _, err := procContinueDebugEvent.Call(
		uintptr(processId),      // 被调试进程的 PID
		uintptr(threadId),       // 被调试线程的 ID
		uintptr(continueStatus), // 继续执行的状态
	)
	if ret == 0 {
		return err
	}
	return nil
}

// GetRip 函数用于获取当前线程的指令指针寄存器（RIP）。
// 返回值:
//   - 当前线程的 RIP 寄存器值。
func GetRip() uint64 {
	var ctx CONTEXT
	ctx.ContextFlags = CONTEXT_FULL
	ret, _, err := procGetThreadContext.Call(
		uintptr(0),                    // 线程句柄
		uintptr(unsafe.Pointer(&ctx)), // 指向 CONTEXT 结构体的指针
	)
	if ret == 0 {
		panic(err)
	}
	return ctx.Rip
}

// WaitForDebug 函数用于等待指定进程的调试事件。
// 参数:
//   - debugEvent: 指向 DEBUG_EVENT 结构体的指针，用于接收调试事件信息。
//
// 返回值:
//   - 指向 DEBUG_EVENT 结构体的指针，包含接收到的调试事件信息。
//   - 错误信息，如果等待过程中出现错误。
func WaitForDebug(debugEvent *DEBUG_EVENT) (*DEBUG_EVENT, error) {
	// 等待调试事件（Windows API调用）
	// 第一个参数是指向 DEBUG_EVENT 结构体的指针，用于接收调试事件信息
	// 第二个参数是等待时间，使用 INFINITE 表示无限期等待
	ret, _, err := procWaitForDebugEvent.Call(
		uintptr(unsafe.Pointer(debugEvent)),
		uintptr(INFINITE),
	)
	// 如果调用失败（返回值为0），则返回错误
	if ret == 0 {
		return nil, err
	}
	return debugEvent, nil
}

func PrintContext(ctx *CONTEXT) {
	fmt.Printf("-------- 寄存器: ------ \n")
	fmt.Printf("Rip: 0x%X\n", ctx.Rip)
	fmt.Printf("Rax: 0x%X\n", ctx.Rax)
	fmt.Printf("Rcx: 0x%X\n", ctx.Rcx)
	fmt.Printf("Rdx: 0x%X\n", ctx.Rdx)
	fmt.Printf("Rbx: 0x%X\n", ctx.Rbx)
	fmt.Printf("Rsp: 0x%X\n", ctx.Rsp)
	fmt.Printf("Rbp: 0x%X\n", ctx.Rbp)
	fmt.Printf("Rsi: 0x%X\n", ctx.Rsi)
	fmt.Printf("Rdi: 0x%X\n", ctx.Rdi)
	fmt.Printf("R8: 0x%X\n", ctx.R8)
	fmt.Printf("R9: 0x%X\n", ctx.R9)
	fmt.Printf("R10: 0x%X\n", ctx.R10)
	fmt.Printf("R11: 0x%X\n", ctx.R11)
	fmt.Printf("R12: 0x%X\n", ctx.R12)
	fmt.Printf("R13: 0x%X\n", ctx.R13)
	fmt.Printf("R14: 0x%X\n", ctx.R14)
	fmt.Printf("R15: 0x%X\n\n", ctx.R15)
	fmt.Printf("-------- 段寄存器: ------ \n")
	fmt.Printf("SegCs: 0x%X\n", ctx.SegCs)
	fmt.Printf("SegDs: 0x%X\n", ctx.SegDs)
	fmt.Printf("SegEs: 0x%X\n", ctx.SegEs)
	fmt.Printf("SegCs: 0x%X\n\n", ctx.SegCs)
	//fmt.Printf("EFlags: 0x%X\n", ctx.EFlags)
	// fmt.Printf("FltSave: 0x%X\n", ctx.FltSave())
	fmt.Printf("-------- 标志寄存器: ------ \n")
	fmt.Printf("CF: %v\n", ctx.GetCF())
	fmt.Printf("ZF: %v\n", ctx.GetZF())
	fmt.Printf("SF: %v\n", ctx.GetSF())
	fmt.Printf("OF: %v\n", ctx.GetOF())
	fmt.Printf("PF: %v\n", ctx.GetPF())
	fmt.Printf("AF: %v\n", ctx.GetAF())
	fmt.Printf("CF: %v\n", ctx.GetCF())
	fmt.Printf("DF: %v\n", ctx.GetDF())
	fmt.Printf("IF: %v\n", ctx.GetIF())
	fmt.Printf("TF: %v\n", ctx.GetTF())
	fmt.Printf("IF: %v\n\n", ctx.GetIF())
	// 可以继续添加其他标志寄存器的打印
}
