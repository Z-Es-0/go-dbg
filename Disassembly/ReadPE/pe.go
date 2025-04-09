/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-09 22:00:41
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-06 19:37:33
 * @FilePath: \ZesOJ\Disassembly\ReadPE\pe.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package readpe

import (
	"fmt"
	"os"

	"github.com/Binject/debug/pe"
)

// PEHeader 包含 PE 文件的头信息
type PEHeader struct {
	DOSHeader        pe.DosHeader
	NTHeader         pe.FileHeader
	OptionalHeader32 *pe.OptionalHeader32
	OptionalHeader64 *pe.OptionalHeader64
}

// ReadPE 读取 PE 文件头
func ReadPE(file *os.File) (*PEHeader, error) {
	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, fmt.Errorf("解析 PE 文件失败: %v", err) // 更新错误描述
	}
	defer peFile.Close()

	peHeader := &PEHeader{
		DOSHeader: peFile.DosHeader,
		NTHeader:  peFile.FileHeader,
	}

	switch optHeader := peFile.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		peHeader.OptionalHeader32 = optHeader
	case *pe.OptionalHeader64:
		peHeader.OptionalHeader64 = optHeader
	default:
		return nil, fmt.Errorf("未知的可选头格式")
	}

	return peHeader, nil
}

// 机器类型转换
func machineTypeToString(machine uint16) string {
	switch machine {
	case 0x8664:
		return "x86-64"
	case 0x14C:
		return "x86"
	case 0xAA64:
		return "ARM64"
	default:
		return "未知架构"
	}
}

// 子系统类型转换
func subsystemToString(subsystem uint16) string {
	switch subsystem {
	case 1:
		return "设备驱动"
	case 2:
		return "Windows GUI"
	case 3:
		return "Windows CUI"
	case 5:
		return "OS/2 CUI"
	case 7:
		return "POSIX CUI"
	case 9:
		return "Windows CE GUI"
	case 10:
		return "EFI应用"
	default:
		return "未知子系统"
	}
}

func PrintPEInfo(pe *PEHeader) {
	// DOS头信息
	fmt.Println("\n=== DOS头 ===")
	fmt.Printf("魔术字: 0x%X (MZ签名)\n", pe.DOSHeader.MZSignature)
	fmt.Printf("NT头偏移: 0x%X\n", pe.DOSHeader.AddressOfNewExeHeader)

	// 文件头信息
	fmt.Println("\n=== NT头 ===")
	fmt.Printf("机器类型: 0x%X (%s)\n", pe.NTHeader.Machine, machineTypeToString(pe.NTHeader.Machine))
	fmt.Printf("节区数量: %d\n", pe.NTHeader.NumberOfSections)
	fmt.Printf("时间日期戳: 0x%X\n", pe.NTHeader.TimeDateStamp)
	fmt.Printf("字符表大小: 0x%X\n", pe.NTHeader.SizeOfOptionalHeader)
	fmt.Printf("特征值: 0x%X\n", pe.NTHeader.Characteristics)

	// 可选头信息
	fmt.Println("\n=== 可选头 ===")
	if pe.OptionalHeader32 != nil {
		oh32 := pe.OptionalHeader32
		fmt.Printf("32bit可选头\n")
		fmt.Printf("主链接器版本: %d\n", oh32.MajorLinkerVersion)
		fmt.Printf("次链接器版本: %d\n", oh32.MinorLinkerVersion)
		fmt.Printf("代码大小: 0x%X\n", oh32.SizeOfCode)
		fmt.Printf("初始化数据大小: 0x%X\n", oh32.SizeOfInitializedData)
		fmt.Printf("未初始化数据大小: 0x%X\n", oh32.SizeOfUninitializedData)
		fmt.Printf("入口点地址: 0x%X\n", oh32.AddressOfEntryPoint)
		fmt.Printf("基地址: 0x%X\n", oh32.ImageBase)
		fmt.Printf("子系统: 0x%X (%s)\n", oh32.Subsystem, subsystemToString(oh32.Subsystem))
	} else if pe.OptionalHeader64 != nil {
		oh64 := pe.OptionalHeader64
		fmt.Printf("64bit可选头\n")
		fmt.Printf("主链接器版本: %d\n", oh64.MajorLinkerVersion)
		fmt.Printf("次链接器版本: %d\n", oh64.MinorLinkerVersion)
		fmt.Printf("代码大小: 0x%X\n", oh64.SizeOfCode)
		fmt.Printf("初始化数据大小: 0x%X\n", oh64.SizeOfInitializedData)
		fmt.Printf("未初始化数据大小: 0x%X\n", oh64.SizeOfUninitializedData)
		fmt.Printf("入口点地址: 0x%X\n", oh64.AddressOfEntryPoint)
		fmt.Printf("基地址: 0x%X\n", oh64.ImageBase)
		fmt.Printf("子系统: 0x%X (%s)\n", oh64.Subsystem, subsystemToString(oh64.Subsystem))
		fmt.Printf("堆保留大小: 0x%X\n", oh64.SizeOfHeapReserve)
		fmt.Printf("堆提交大小: 0x%X\n", oh64.SizeOfHeapCommit)
		fmt.Printf("栈保留大小: 0x%X\n", oh64.SizeOfStackReserve)
		fmt.Printf("栈提交大小: 0x%X\n", oh64.SizeOfStackCommit)
	}
}
