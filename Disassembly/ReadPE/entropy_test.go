/*
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-03-11 00:54:53
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-03-12 22:14:57
 * @FilePath: \ZesOJ\Disassembly\ReadPE\entropy_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package readpe

import (
	"fmt"
	"os"
	"testing"
)

func TestCalculateEntropy(t *testing.T) {
	filePath := "E:/Zesoj/sever/test.exe"
	file, err := os.Open(filePath) // 读取文件
	if err != nil {
		t.Errorf("Open file failed: %v", err)
		return
	}
	cal, err := CalculateEntropy(file)
	if err != nil {
		t.Errorf("CalculateEntropy failed: %v", err)
		return
	} else {
		fmt.Println("answer is ", cal)
	}

}

func TestGetSectionNames(t *testing.T) {
	filePath := "E:/Zesoj/sever/test.exe"
	file, err := os.Open(filePath) // 读取文件
	if err != nil {
		t.Errorf("Open file failed: %v", err)
		return
	}
	defer file.Close()

	sectionNames, err := GetSectionNames(file)
	if err != nil {
		t.Errorf("GetSectionNames failed: %v", err)
		return
	}
	for i := 0; i < len(sectionNames); i++ {
		fmt.Println(sectionNames[i])
	}
}
