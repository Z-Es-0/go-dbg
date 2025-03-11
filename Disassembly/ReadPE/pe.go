/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-09 22:00:41
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-11 14:02:43
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
func ReadPE(filePath string) (*PEHeader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, fmt.Errorf("解析 PE 文件失败: %v", err)
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
