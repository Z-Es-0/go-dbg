/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-06 14:36:16
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-18 03:53:05
 * @FilePath: \ZesOJ\Disassembly\gdb\debug.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gdb

import (
	"syscall"
	"zesdbg/Disassembly/analyse"
)

func Dbginit(exePath string, cmdLine string) bool {
	process, thread, err := CreateAndBlockProcess(exePath, cmdLine)
	if err != nil {
		return false
	}
	strdata, err := analyse.FindHardcodedStrings(exePath)
	if err != nil {
		panic(err)
	}

	dbgMachine := &DbgMachine{
		process: process, // 进程句柄,
		thread:  thread,  // 线程句柄,
		str:     &strdata,

		breakpoints: make(map[uintptr]*Dbgbreak),
		textdata:    make(map[uintptr]*Directive),
	}
	defer syscall.CloseHandle(dbgMachine.process)
	defer syscall.CloseHandle(dbgMachine.thread)

	dbgMachine.Run()
	return true
}
