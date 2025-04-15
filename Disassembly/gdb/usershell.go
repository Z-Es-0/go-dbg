/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-10 00:52:34
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-15 16:49:45
 * @FilePath: \ZesOJ\Disassembly\gdb\usershell.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (d *DbgMachine) ShellMain() {
	for d.Shellgo() {

	}
}

func (d *DbgMachine) Shellgo() bool {
	var cmd string

	fmt.Print("gdb>")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		cmd = scanner.Text()
	} else {
		if err := scanner.Err(); err != nil {
			fmt.Println("读取输入时出错:", err)
		}
	}
	switch cmd {
	case "run", "r":
		return false

	// case "break", "b":
	// 	d.breakgo()
	// 	return true
	// case "memory", "m":
	// 	d.memorygo()
	// 	return true

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
		fmt.Println("break <address>, b <address>: 设置断点")
		fmt.Println("memory  <address> <size>, m  <address> <size>: 查看内存")
		fmt.Println("context, c: 查看线程上下文(当前寄存器状态)")
		fmt.Println("quit, q: 退出程序")
		fmt.Println("list, l: 列出断点")
		fmt.Println("help, h: 查看帮助")

	default:
		op := cmd[0:1]
		switch op {
		case "b":

			if len(cmd) > 5 && (cmd[:4] == "break" || cmd[0] == 'b') {
				if cmd[:4] == "break" {
					cmd = cmd[5:]
				} else {
					cmd = cmd[2:]
				}

				addr := cmd[5:]
				d.breakgo(addr)

			} else {
				fmt.Println("error command the break command is break <address>")
			}

		case "m":
			if len(cmd) > 7 && (cmd[:6] == "memory" || cmd[0] == 'm') {
				if cmd[:6] == "memory" {
					cmd = cmd[7:]
				} else {
					cmd = cmd[2:]
				}

				cmd = strings.TrimSpace(cmd)
				parts := strings.Fields(cmd)

				if len(parts) >= 2 {
					address := parts[0]
					sizeStr := parts[1]

					var size int
					if _, err := fmt.Sscanf(sizeStr, "%d", &size); err != nil {
						fmt.Printf("解析内存大小 %s 为整数时出错: %v\n", sizeStr, err)

					}

					// 尝试将地址字符串转换为 uint64 类型
					var addrUint64 uint64
					if _, err := fmt.Sscanf(address, "%x", &addrUint64); err != nil {
						fmt.Printf("解析地址 %s 为十六进制数字时出错: %v\n", address, err)

					}

					// 调用函数查看指定地址和大小的内存
					d.WatchMenoryandDisassemble(uintptr(addrUint64), size)
					return true
				} else {
					fmt.Println("错误: 内存命令需要地址和大小参数")
					return true
				}

			} else {
				if cmd[0] == 'm' || cmd[:6] == "memory" {
					fmt.Println("this is the memory begin as Rip and the size is 200 bytes")
					d.memorygo()
					return true
				} else {
					fmt.Println("error command the memory command is memory  <address> <size>")
					return true
				}
			}

		}

		return true

	}
	return true
}

func (d *DbgMachine) breakgo(addressstr string) {

	addressstr = strings.TrimSpace(addressstr)
	// if strings.HasPrefix(addressstr, "0x") || strings.HasPrefix(addressstr, "0X") {
	// 	addressstr = addressstr[2:]
	// }

	var addrInt uint64
	if _, err := fmt.Sscanf(addressstr, "%x", &addrInt); err != nil {
		fmt.Printf("解析地址 %s 为十六进制数字时出错: %v\n", addressstr, err)
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

		//d.textdata[uintptr(directive.Address)] = &directive
	}

}

func (d *DbgMachine) WatchMenoryandDisassemble(address uintptr, size int) {
	// 读取内存中的指令
	bytedata, err := ReadProcessMemory(d.process, uintptr(address), uint(size))
	if err != nil {
		fmt.Println("读取内存失败:", err)
		return
	}

	// 反汇编读取的指令
	data := DisassembleRange(bytedata, uint64(address), 64) //  64 位架构

	// 遍历反汇编结果并存储到 textdata 映射中
	for _, directive := range *data {
		// 按16进制地址 - 16进制 - 汇编的格式打印
		hexCodesStr := ""
		for _, code := range directive.HexCodes {
			hexCodesStr += fmt.Sprintf("%02X ", code)
		}
		fmt.Printf("%016X - %-30s - %-25s\n", directive.Address, hexCodesStr, directive.Armcode)

	}
}
