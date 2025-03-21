/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-21 18:33:43
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-21 19:49:10
 * @FilePath: \ZesOJ\Disassembly\ReadPE\makePEdate.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package readpe

import (
	"debug/pe"
	"os"
)

type PEInfo struct {
	header         *PEHeader
	sectionname    *[]string
	sectionentropy *map[string]float64
	sectionsize    *map[string]uint32
	sectiondata    *map[string][]byte
	sectionispack  *map[string]bool
}

func GetPEInfo(file *os.File) (*PEInfo, error) { // 获取PE文件信息
	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, err
	}
	defer peFile.Close()

	peinfo := &PEInfo{}
	peinfo.header, _ = ReadPE(file)
	peinfo.sectionname, _ = GetSectionNames(file)
	peinfo.sectionentropy, _ = GetsegmentEntropy(file)
	peinfo.sectionispack = IsPacked(peinfo.sectionentropy)
	peinfo.sectionsize, _ = GetSectionSize(file)
	peinfo.sectiondata, _ = GetSectionData(file)

	return peinfo, nil
}
