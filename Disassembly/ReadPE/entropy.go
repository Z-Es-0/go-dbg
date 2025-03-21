/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-10 09:07:43
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-21 19:44:51
 * @FilePath: \ZesOJ\Disassembly\ReadPE\entropy.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package readpe

import (
	"debug/pe"
	"io"
	"math"
	"os"
)

// GetSectionNames 函数用于返回可执行文件的所有段名称
func GetSectionNames(file *os.File) (*[]string, error) {

	// 解析PE文件头
	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, err
	}
	defer peFile.Close()

	// 用于存储段名称的切片
	sectionNames := make([]string, 0, len(peFile.Sections))
	for _, section := range peFile.Sections {
		sectionNames = append(sectionNames, section.Name)
	}

	return &sectionNames, nil
}

// 计算熵
func op(buffer *[]byte) float64 {
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
func CalculateEntropy(file *os.File) (float64, error) {

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

// 计算各段熵
func GetsegmentEntropy(file *os.File) (*map[string]float64, error) {
	sptr, err := GetSectionNames(file)
	s := *sptr
	if err != nil {
		return nil, err
	}

	mp := map[string]float64{}
	for i := 0; i < len(s); i++ {
		mp[s[i]] = 0
	}

	// 解析PE文件头
	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, err
	}
	defer peFile.Close()

	headmp := make(map[string][]byte, 0)
	for i := 0; i < len(s); i++ {
		headmp[s[i]] = make([]byte, 0)
	}
	for _, section := range peFile.Sections {
		data, err := section.Data()
		if err != nil {
			return nil, err
		}
		headmp[section.Name] = data
	}

	// 计算各段的信息熵
	for i := 0; i < len(s); i++ {

		// 将 slice 传递给 op 函数，因为 op 函数接收的是一个指向 slice 的指针。
		buffer := headmp[s[i]]
		mp[s[i]] = op(&buffer)
	}
	return &mp, nil
}

func IsPacked(mpptr *map[string]float64) *map[string]bool {
	mp := *mpptr
	ans := make(map[string]bool, len(mp))

	for i, v := range mp {
		if v > 7.0 {
			ans[i] = true
		} else {
			ans[i] = false
		}

	}

	return &ans

}

// GetSectionSize 函数用于返回可执行文件中各段的大小
func GetSectionSize(file *os.File) (*map[string]uint32, error) {

	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, err
	}
	defer peFile.Close()

	sectionSizes := make(map[string]uint32, len(peFile.Sections))
	for _, section := range peFile.Sections {
		sectionSizes[section.Name] = section.Size
	}

	return &sectionSizes, nil
}

// GetSectionData 函数用于返回可执行文件中各段的数据
func GetSectionData(file *os.File) (*map[string][]byte, error) {

	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, err
	}
	defer peFile.Close()

	sectionData := make(map[string][]byte, len(peFile.Sections))
	for _, section := range peFile.Sections {
		data, err := section.Data()
		if err != nil {
			return nil, err
		}
		sectionData[section.Name] = data
	}

	return &sectionData, nil
}
