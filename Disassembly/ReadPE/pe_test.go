/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-09 23:58:33
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-12 21:45:30
 * @FilePath: \ZesOJ\Disassembly\ReadPE\pe_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package readpe

import (
	"os"
	"testing"
	"time"
)

func TestReadPE(t *testing.T) {
	filePath := "E:/Zesoj/sever/test.exe"
	file, _ := os.Open(filePath)
	defer file.Close()
	pe, err := ReadPE(file)
	// ntHeaderOffset := int64(pe.DOSHeader.Lfanew)
	// t.Logf("NT头文件偏移: 0x%X", ntHeaderOffset)
	if err != nil {
		// peSignatureOffset := int64(pe.DOSHeader.Lfanew) + 4 // PE签名位于NT头起始位置后4字节
		// t.Logf("PE签名文件偏移: 0x%X", peSignatureOffset)
		t.Errorf("ReadPE failed: %v", err)
		return
	}
	if pe == nil {
		t.Errorf("ReadPE returned nil PEHeader")
	}
	printPEInfo(t, pe)
}

func printPEInfo(t *testing.T, pe *PEHeader) {
	// DOS头信息
	t.Log("\n=== DOS头 ===")
	t.Logf("魔术字: 0x%X (MZ签名)", pe.DOSHeader.MZSignature)
	t.Logf("NT头偏移: 0x%X", pe.DOSHeader.AddressOfNewExeHeader)

	// 文件头信息
	t.Log("\n=== NT头 ===")
	t.Logf("机器类型: 0x%X (%s)", pe.NTHeader.Machine, machineTypeToString(pe.NTHeader.Machine))
	t.Logf("节区数量: %d", pe.NTHeader.NumberOfSections)
	t.Logf("时间戳: 0x%X (%s)", pe.NTHeader.TimeDateStamp,
		time.Unix(int64(pe.NTHeader.TimeDateStamp), 0).UTC())

	// 可选头信息
	t.Log("\n=== 可选头 ===")
	// 打印OptionalHeader32或OptionalHeader64头的信息
	if pe.OptionalHeader32 != nil {
		oh32 := pe.OptionalHeader32
		t.Logf("32bit可选头")
		// t.Logf("主操作系统版本: %d", oh32.MajorOperatingSystemVersion)
		// t.Logf("次操作系统版本: %d", oh32.MinorOperatingSystemVersion)
		t.Logf("主链接器版本: %d", oh32.MajorLinkerVersion)
		t.Logf("次链接器版本: %d", oh32.MinorLinkerVersion)
		t.Logf("代码大小: 0x%X", oh32.SizeOfCode)
		t.Logf("初始化数据大小: 0x%X", oh32.SizeOfInitializedData)
		t.Logf("未初始化数据大小: 0x%X", oh32.SizeOfUninitializedData)
		t.Logf("入口点地址: 0x%X", oh32.AddressOfEntryPoint)
		t.Logf("基地址: 0x%X", oh32.ImageBase)
		t.Logf("子系统: 0x%X (%s)", oh32.Subsystem, subsystemToString(oh32.Subsystem))
	} else if pe.OptionalHeader64 != nil {
		oh64 := pe.OptionalHeader64
		t.Logf("64bit可选头")
		// t.Logf("主操作系统版本: %d", oh64.MajorOperatingSystemVersion)
		// t.Logf("次操作系统版本: %d", oh64.MinorOperatingSystemVersion)
		t.Logf("主链接器版本: %d", oh64.MajorLinkerVersion)
		t.Logf("次链接器版本: %d", oh64.MinorLinkerVersion)
		t.Logf("代码大小: 0x%X", oh64.SizeOfCode)
		t.Logf("初始化数据大小: 0x%X", oh64.SizeOfInitializedData)
		t.Logf("未初始化数据大小: 0x%X", oh64.SizeOfUninitializedData)
		t.Logf("入口点地址: 0x%X", oh64.AddressOfEntryPoint)
		t.Logf("基地址: 0x%X", oh64.ImageBase)
		t.Logf("子系统: 0x%X (%s)", oh64.Subsystem, subsystemToString(oh64.Subsystem))
		t.Logf("堆保留大小: 0x%X", oh64.SizeOfHeapReserve)
		t.Logf("堆提交大小: 0x%X", oh64.SizeOfHeapCommit)
		t.Logf("栈保留大小: 0x%X", oh64.SizeOfStackReserve)
		t.Logf("栈提交大小: 0x%X", oh64.SizeOfStackCommit)
	}
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
