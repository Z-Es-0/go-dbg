/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-10 09:07:43
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-11 14:35:59
 * @FilePath: \ZesOJ\Disassembly\ReadPE\entropy.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package readpe

import (
	"debug/pe"
	"fmt"
	"io"
	"math"
	"os"
)

func op(buffer *[]byte) float64 { // 计算熵
	var ans float64 = 0
	mp := make(map[byte]int, 0)
	var size int = (int)(len(*buffer))
	up := *buffer
	for i := 0; i < size; i++ {
		mp[up[i]]++
	}
	for i := 0; i < 256; i++ {
		ui := float64(mp[byte(i)]) / float64(size)
		if ui > 0 {
			ans += -float64(ui) * math.Log2(float64(ui))
		}

	}
	return ans
}

// 计算文件信息熵
//
//	Accept test
func CalculateEntropy(filePath string) (float64, error) {
	file, err := os.Open(filePath) // 读取文件
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// ------------------------------
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return 0, err
	}

	// test
	// for i := 0; i < len(buffer); i++ {
	// 	fmt.Printf("  %x", buffer[i])
	// }
	result := op(&buffer)
	return result, nil
}

/*
TODO use procexp_judge_pack.cpp
*/
// func GetPackedCheck(entropy float64) int {

// }

// TODO 分别计算data,rdata,text,bss,段的信息熵
func GetsegmentEntropy(filePath string) (map[string]float64, error) {
	mp := map[string]float64{
		".data":  0,
		".rdata": 0,
		".text":  0,
		".bss":   0,
	}
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 解析PE文件头
	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, err
	}
	defer peFile.Close()

	// 查找data段
	var dataSection *pe.Section
	for _, section := range peFile.Sections {
		if section.Name == ".data" {
			dataSection = section
			break
		}
	}

	if dataSection == nil {
		return nil, fmt.Errorf("data section not found")
	}

	// 读取data段内容
	data := make([]byte, dataSection.Size)
	_, err = file.ReadAt(data, int64(dataSection.Offset))
	if err != nil {
		return nil, err
	}

	// 计算data段的信息熵
	mp[".data"] = op(&data)
	return mp, nil
}
