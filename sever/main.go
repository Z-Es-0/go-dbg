/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-12 13:35:04
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-12 16:11:24
 * @FilePath: \ZesOJ\Disassembly\sever\main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"fmt"
	"os"

	readpe "zesdbg/Disassembly/ReadPE"
)

func main() {
	filePath := "test.exe" // 添加文件路径定义
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	pe, err := readpe.ReadPE(file)
	fmt.Println(pe.DOSHeader.AddressOfNewExeHeader)

	// 初始化PE解析
	// peFile, _ := pe.NewFile(file)

}
