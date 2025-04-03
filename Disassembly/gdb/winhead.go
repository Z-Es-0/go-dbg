/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-03 14:47:50
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-04 00:00:12
 * @FilePath: \ZesOJ\Disassembly\gdb\winhead.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

import (
	"fmt"
	"reflect"
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
