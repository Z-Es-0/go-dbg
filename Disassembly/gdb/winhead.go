/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-03 14:47:50
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-03 22:16:58
 * @FilePath: \ZesOJ\Disassembly\gdb\winhead.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

import (
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

	EFlags uint32

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
