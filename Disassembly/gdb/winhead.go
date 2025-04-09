/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-03 14:47:50
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-10 02:17:58
 * @FilePath: \ZesOJ\Disassembly\gdb\winhead.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

import (
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

const (
	CONTEXT_CONTROL         = 0x10001                                              // 控制寄存器
	CONTEXT_INTEGER         = 0x10002                                              // 整数寄存器
	CONTEXT_SEGMENTS        = 0x10004                                              // 段寄存器
	CONTEXT_FLOATING_POINT  = 0x10008                                              // 浮点寄存器
	CONTEXT_DEBUG_REGISTERS = 0x10010                                              // 调试寄存器
	CONTEXT_XSTATE          = 0x10020                                              // 扩展状态寄存器
	CONTEXT_FULL            = CONTEXT_CONTROL | CONTEXT_INTEGER | CONTEXT_SEGMENTS //所有用户态寄存器

)

// 标志位掩码
const (
	CF = (uint32)(1 << 0)  // 进位标志
	PF = (uint32)(1 << 2)  // 奇偶校验标志
	AF = (uint32)(1 << 4)  // 辅助进位标志
	ZF = (uint32)(1 << 6)  // 零标志
	SF = (uint32)(1 << 7)  // 符号标志
	TF = (uint32)(1 << 8)  // 陷阱标志
	IF = (uint32)(1 << 9)  // 中断使能标志
	DF = (uint32)(1 << 10) // 方向标志
	OF = (uint32)(1 << 11) // 溢出标志
)

//go:align 16
type CONTEXT struct {
	P1Home uint64
	P2Home uint64
	P3Home uint64
	P4Home uint64
	P5Home uint64
	P6Home uint64

	ContextFlags uint32
	MxCsr        uint32

	SegCs uint16
	SegDs uint16
	SegEs uint16
	SegFs uint16
	SegGs uint16
	SegSs uint16

	EFlags uint32 // 32位标志寄存器

	Dr0 uint64
	Dr1 uint64
	Dr2 uint64
	Dr3 uint64
	Dr6 uint64
	Dr7 uint64

	Rax uint64
	Rcx uint64
	Rdx uint64
	Rbx uint64
	Rsp uint64
	Rbp uint64
	Rsi uint64
	Rdi uint64
	R8  uint64
	R9  uint64
	R10 uint64
	R11 uint64
	R12 uint64
	R13 uint64
	R14 uint64
	R15 uint64
	Rip uint64

	// 联合体处理（512字节）
	unionArea [512]byte // XMM_SAVE_AREA32 和其联合体部分

	VectorRegister [26]M128A

	VectorControl        uint64
	DebugControl         uint64
	LastBranchToRip      uint64
	LastBranchFromRip    uint64
	LastExceptionToRip   uint64
	LastExceptionFromRip uint64
}

// M128A 定义（16字节）
type M128A struct {
	Low  uint64
	High uint64
}

// XMM_SAVE_AREA32 结构体定义
type XMM_SAVE_AREA32 struct {
	ControlWord    uint16
	StatusWord     uint16
	TagWord        byte
	Reserved1      byte
	ErrorOpcode    uint16
	ErrorOffset    uint32
	ErrorSelector  uint16
	Reserved2      uint16
	DataOffset     uint32
	DataSelector   uint16
	Reserved3      uint16
	MxCsr          uint32
	MxCsr_Mask     uint32
	FloatRegisters [8]M128A
	XmmRegisters   [16]M128A
	Reserved4      [96]byte
}

// 联合体访问方法
func (c *CONTEXT) FltSave() *XMM_SAVE_AREA32 {
	return (*XMM_SAVE_AREA32)(unsafe.Pointer(&c.unionArea[0]))
}

func (c *CONTEXT) Q() *[16]M128A {
	return (*[16]M128A)(unsafe.Pointer(&c.unionArea[0]))
}

func (c *CONTEXT) D() *[32]uint64 {
	return (*[32]uint64)(unsafe.Pointer(&c.unionArea[0]))
}

func (c *CONTEXT) S() *[32]uint32 {
	return (*[32]uint32)(unsafe.Pointer(&c.unionArea[0]))
}

// OffsetOfFieldByName 根据字段名计算CONTEXT结构体中对应字段的偏移量
func OffsetOfFieldByName(fieldName string) (uintptr, error) {
	var context CONTEXT
	value := reflect.ValueOf(&context).Elem()
	field := value.FieldByName(fieldName)
	if !field.IsValid() {
		return 0, fmt.Errorf("field %s not found in CONTEXT struct", fieldName)
	}
	return field.UnsafeAddr() - value.UnsafeAddr(), nil
}

// GetCF 获取进位标志 (CF)
func (c *CONTEXT) GetCF() bool {
	return c.EFlags&CF != 0
}

// GetPF 获取奇偶校验标志 (PF)
func (c *CONTEXT) GetPF() bool {
	return c.EFlags&PF != 0
}

// GetAF 获取辅助进位标志 (AF)
func (c *CONTEXT) GetAF() bool {
	return c.EFlags&AF != 0
}

// GetZF 获取零标志 (ZF)
func (c *CONTEXT) GetZF() bool {
	return c.EFlags&ZF != 0
}

// GetSF 获取符号标志 (SF)
func (c *CONTEXT) GetSF() bool {
	return c.EFlags&SF != 0
}

// GetTF 获取陷阱标志 (TF)
func (c *CONTEXT) GetTF() bool {
	return c.EFlags&TF != 0
}

// GetIF 获取中断使能标志 (IF)
func (c *CONTEXT) GetIF() bool {
	return c.EFlags&IF != 0
}

// GetDF 获取方向标志 (DF)
func (c *CONTEXT) GetDF() bool {
	return c.EFlags&DF != 0
}

// GetOF 获取溢出标志 (OF)
func (c *CONTEXT) GetOF() bool {
	return c.EFlags&OF != 0
}

// 调试器相关常量定义
const (
	INFINITE                   = 0xFFFFFFFF // 无限等待
	EXCEPTION_DEBUG_EVENT      = 0x00000001 // 异常调试事件
	CREATE_PROCESS_DEBUG_EVENT = 0x00000003 // 进程创建事件
	EXIT_PROCESS_DEBUG_EVENT   = 0x00000005 // 进程退出事件
	CREATE_THREAD_DEBUG_EVENT  = 0x00000002 // 线程创建事件
	LOAD_DLL_DEBUG_EVENT       = 0x00000006 // DLL加载事件
)

// DEBUG_EVENT 调试事件结构体
type DEBUG_EVENT struct {
	DebugEventCode uint32    // 调试事件类型 (4字节)
	ProcessId      uint32    // 进程ID (4字节)
	ThreadId       uint32    // 线程ID (4字节)
	_              [4]byte   // 对齐填充 (4字节)
	u              [168]byte // 原始字节缓冲区
}

// 定义一个泛型函数，用于获取DEBUG_EVENT结构体中的联合字段
// 该函数接收一个DEBUG_EVENT 指针类型的参数e，并返回一个泛型类型T的值
// 通过unsafe.Pointer将DEBUG_EVENT结构体中的u字段（原始字节缓冲区）转换为泛型类型T的指针
// 然后解引用该指针，返回T类型的值
// 这样可以方便地从DEBUG_EVENT结构体中提取不同类型的联合字段数据
func GetUnion[T any](e *DEBUG_EVENT) T {
	return *(*T)(unsafe.Pointer(&(e.u)))
}

// 包含调试器可以使用的异常信息
type _EXCEPTION_DEBUG_INFO struct {
	ExceptionRecord EXCEPTION_RECORD
	dwFirstChance   uint32
}

// 为结构体 _EXCEPTION_DEBUG_INFO 和它的指针起别名
type EXCEPTION_DEBUG_INFO = _EXCEPTION_DEBUG_INFO
type LPEXCEPTION_DEBUG_INFO = *_EXCEPTION_DEBUG_INFO

// CREATE_THREAD_DEBUG_INFO 线程创建调试信息，包含线程创建时的相关信息
type CREATE_THREAD_DEBUG_INFO struct {
	hThread           syscall.Handle // 新创建线程的句柄
	lpThreadLocalBase uintptr        // 线程局部存储的基地址
	lpStartAddress    uintptr        // 线程的实际入口点地址
}

// EXIT_THREAD_DEBUG_INFO 线程退出调试信息，包含线程退出时的相关信息
type EXIT_THREAD_DEBUG_INFO struct {
	dwExitCode uint32 // 线程退出代码
}

// EXIT_PROCESS_DEBUG_INFO 进程退出调试信息，包含进程退出时的相关信息
type EXIT_PROCESS_DEBUG_INFO struct {
	dwExitCode uint32 // 进程退出代码
}

// LOAD_DLL_DEBUG_INFO DLL加载调试信息，包含DLL加载时的相关信息
type LOAD_DLL_DEBUG_INFO struct {
	hFile                 syscall.Handle // DLL文件句柄
	lpBaseOfDll           uintptr        // DLL映像的基地址
	dwDebugInfoFileOffset uint32         // 调试信息文件的偏移量
	nDebugInfoSize        uint32         // 调试信息的大小
	lpImageName           uintptr        // DLL映像名称的指针
	fUnicode              uint16         // 表示映像名称是否为Unicode编码
}

// UNLOAD_DLL_DEBUG_INFO DLL卸载调试信息，包含DLL卸载时的相关信息
type UNLOAD_DLL_DEBUG_INFO struct {
	lpBaseOfDll uintptr // DLL映像的基地址
}

// OUTPUT_DEBUG_STRING_INFO 调试输出字符串信息，包含调试输出字符串时的相关信息
type OUTPUT_DEBUG_STRING_INFO struct {
	lpDebugStringData  uintptr // 调试输出字符串的指针
	fUnicode           uint16  // 表示字符串是否为Unicode编码
	nDebugStringLength uint16  // 调试输出字符串的长度
}

// RIP_INFO RIP调试信息，包含RIP调试时的相关信息
type RIP_INFO struct {
	dwError uint32 // RIP错误代码
	dwType  uint32 // RIP调试信息类型
}

// EXCEPTION_RECORD 异常记录，用于存储异常的详细信息
type EXCEPTION_RECORD struct {
	// 异常的代码，用于标识不同类型的异常
	ExceptionCode uint32
	// 异常的标志，提供关于异常的额外信息
	ExceptionFlags uint32
	// 指向另一个异常记录的指针，用于嵌套异常
	ExceptionRecord *EXCEPTION_RECORD
	// 异常发生时的处理器状态
	ExceptionAddress uintptr
	// 异常的编号，用于标识异常的类型
	NumberParameters uint32
	// 异常的参数数组，用于存储异常的相关数据
	ParameterArray [15]uintptr
}

// CREATE_PROCESS_DEBUG_INFO 进程创建调试信息，包含进程创建时的相关信息
type CREATE_PROCESS_DEBUG_INFO struct {
	// 进程创建时使用的文件句柄
	hFile syscall.Handle
	// 新创建进程的句柄
	hProcess syscall.Handle
	// 新创建线程的句柄
	hThread syscall.Handle
	// 进程映像的基地址
	lpBaseOfImage uintptr
	// 调试信息文件的偏移量
	dwDebugInfoFileOffset uint32
	// 调试信息的大小
	nDebugInfoSize uint32
	// 线程局部存储的基地址
	lpThreadLocalBase uintptr
	// 进程的实际入口点地址
	lpStartAddress uintptr
	// 进程映像名称的指针
	lpImageName uintptr
	// 表示映像名称是否为Unicode编码
	fUnicode uint16
}

// 新增 OpenThread 封装
func OpenThread(desiredAccess uint32, inheritHandle bool, threadId uint32) (syscall.Handle, error) {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procOpenThread := modkernel32.NewProc("OpenThread")

	var inherit uint32 = 0
	if inheritHandle {
		inherit = 1
	}

	r0, _, e1 := procOpenThread.Call(
		uintptr(desiredAccess),
		uintptr(inherit),
		uintptr(threadId),
	)
	if r0 == 0 {
		return 0, e1
	}
	return syscall.Handle(r0), nil
}
