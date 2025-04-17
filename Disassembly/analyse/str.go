package analyse

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"unicode/utf16"
)

type StrDate struct {
	Addr    uint64
	Str     string
	Segment string
}

// FindHardcodedStrings 分析PE文件中的硬编码字符串
func FindHardcodedStrings(filePath string) ([]StrDate, error) {
	// 打开目标文件
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 解析PE文件
	peFile, err := pe.NewFile(f)
	if err != nil {
		return nil, err
	}
	defer peFile.Close()

	// var stringsFound []string

	var res []StrDate

	// 遍历所有节区
	for _, section := range peFile.Sections {

		name := section.Name
		// var stringsFound []string

		switch name {
		case ".rdata":
			data, err := section.Data()
			if err != nil {
				continue
			}

			// 扫描ASCII字符串并记录信息
			asciiStrings := scanForAsciiStrings(data)
			for _, s := range asciiStrings {
				startIndex := bytes.Index(data, []byte(s))
				if startIndex != -1 {
					addr := uint64(section.VirtualAddress) + uint64(startIndex)
					res = append(res, StrDate{
						Addr:    addr,
						Str:     s,
						Segment: section.Name,
					})
				}
			}

			// 扫描Unicode字符串并记录信息
			unicodeStrings := scanForUnicodeStrings(data)
			for _, s := range unicodeStrings {
				utf16Bytes := make([]byte, len(s)*2)
				for i, r := range s {
					binary.LittleEndian.PutUint16(utf16Bytes[i*2:], uint16(r))
				}
				startIndex := bytes.Index(data, utf16Bytes)
				if startIndex != -1 {
					addr := uint64(section.VirtualAddress) + uint64(startIndex)
					res = append(res, StrDate{
						Addr:    addr,
						Str:     s,
						Segment: section.Name,
					})
				}
			}

		}
	}
	return res, nil
}

// 扫描ASCII字符串（可打印字符连续序列）
func scanForAsciiStrings(data []byte) []string {
	var result []string
	var currentStr []byte
	const minLength = 4 // 最小有效字符串长度

	for _, b := range data {
		if b >= 0x20 && b <= 0x7E { // 可打印ASCII范围
			currentStr = append(currentStr, b)
		} else {
			if len(currentStr) >= minLength {
				result = append(result, string(currentStr))
			}
			currentStr = nil
		}
	}
	// 处理最后可能存在的字符串
	if len(currentStr) >= minLength {
		result = append(result, string(currentStr))
	}
	return result
}

// 扫描Unicode字符串（小端序）
func scanForUnicodeStrings(data []byte) []string {
	var result []string
	var currentStr []uint16
	const minLength = 4

	for i := 0; i < len(data)-1; i += 2 {
		// 读取小端序的UTF-16字符
		u := binary.LittleEndian.Uint16(data[i:])

		// 基本多文种平面字符且可打印
		if u >= 0x20 && u <= 0x7E && u != 0x7F {
			currentStr = append(currentStr, u)
		} else {
			if len(currentStr) >= minLength {
				result = append(result, string(utf16.Decode(currentStr)))
			}
			currentStr = nil
		}
	}
	return result
}

func PrintStr(str string) {
	strings, err := FindHardcodedStrings(str)
	if err != nil {
		panic(err)
	}

	for _, s := range strings {
		fmt.Printf("%-20x %-20s %s\n", s.Addr, s.Segment, s.Str)
	}
}

func SelectStr(data *[]StrDate, str string) ([]StrDate, error) {
	// 编译正则表达式
	re, err := regexp.Compile(str)
	if err != nil {
		return nil, err
	}

	var result []StrDate
	// 遍历数据，查找符合正则表达式的字符串
	for _, item := range *data {
		if re.MatchString(item.Str) {
			result = append(result, item)
		}
	}

	return result, nil
}
