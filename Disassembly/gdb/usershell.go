/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-10 00:52:34
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-11 14:58:49
 * @FilePath: \ZesOJ\Disassembly\gdb\usershell.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

import (
	"fmt"
	"os"
)

func (d *DbgMachine) ShellMain() {
	for d.Shellgo() {

	}
}

func (d *DbgMachine) Shellgo() bool {
	var cmd string
	fmt.Print("gdb>")
	fmt.Scanln(&cmd)
	switch cmd {
	case "run", "r":
		return false

	case "break", "b":
		d.breakgo()
		return true
	case "memory", "m":
		d.memorygo()
		return true

	case "context", "c":
		context, err := GetThreadContext(d.thread)
		if err != nil {
			fmt.Println("获取线程上下文失败:", err)
			return true
		}
		PrintContext(context)
		return true

	case "quit", "q":
		os.Exit(0)
		return false

	case "list", "l":
		for _, directive := range d.breakpoints {
			fmt.Printf("断点地址: 0x%X\n", directive.address)
			fmt.Println("指令: ", directive.rawcode.Armcode)
		}

		// 反汇编读取的指令

	// case "stepping","s":

	case "help", "h":
		fmt.Println("run, r: 运行程序")
		fmt.Println("break, b: 设置断点")
		fmt.Println("memory, m: 查看内存")
		fmt.Println("context, c: 查看线程上下文")
		fmt.Println("quit, q: 退出程序")
		fmt.Println("list, l: 列出断点")
		fmt.Println("help, h: 查看帮助")
	}
	return true
}

func (d *DbgMachine) breakgo() {
	var addr string
	fmt.Print("输入断点地址: ")
	fmt.Scanln(&addr)
	// 将输入的字符串作为16进制数字解析
	var addrInt uint64
	if _, err := fmt.Sscanf(addr, "%x", &addrInt); err != nil {
		fmt.Printf("解析地址 %s 为十六进制数字时出错: %v\n", addr, err)
		return
	}

	// 转换为 uintptr 类型
	addrPtr := uintptr(addrInt)

	err := d.SetBreakpoint(addrPtr)
	if err != nil {
		fmt.Printf("设置断点失败: %v\n", err)
		return
	}

}

func (d *DbgMachine) memorygo() {
	if d.textdata == nil {
		d.textdata = make(map[uintptr]*Directive)
	}

	// 获取线程上下文
	context, err := GetThreadContext(d.thread)
	if err != nil {
		fmt.Println("获取线程上下文失败:", err)
		return
	}

	// 获取指令指针 (RIP)
	rip := context.Rip

	// 读取内存中的指令
	bytedata, err := ReadProcessMemory(d.process, uintptr(rip), 200)
	if err != nil {
		fmt.Println("读取内存失败:", err)
		return
	}

	// 反汇编读取的指令
	data := DisassembleRange(bytedata, rip, 64) //  64 位架构

	// 遍历反汇编结果并存储到 textdata 映射中
	for _, directive := range *data {
		// 按16进制地址 - 16进制 - 汇编的格式打印
		hexCodesStr := ""
		for _, code := range directive.HexCodes {
			hexCodesStr += fmt.Sprintf("%02X ", code)
		}
		fmt.Printf("%016X - %-30s - %-25s\n", directive.Address, hexCodesStr, directive.Armcode)

		// d.textdata[uintptr(directive.Address)] = &directive
	}

}
