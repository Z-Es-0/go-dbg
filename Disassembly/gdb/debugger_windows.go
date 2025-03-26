/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-21 21:49:22
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-26 17:15:49
 * @FilePath: \ZesOJ\Disassembly\gdb\debugger_windows.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

// 引入windows包以使用Windows API
import (
	"syscall"
	"unsafe"
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
