package readpe

import (
	"os"
	"testing"
)

func TestGetPEInfo(t *testing.T) {
	filePath := "E:/Zesoj/sever/test.exe"
	file, err := os.Open(filePath) // 读取文件
	if err != nil {
		t.Errorf("Open file failed: %v", err)
		return
	}
	peinfo, err := GetPEInfo(file)
	if err != nil {
		t.Errorf("GetPEInfo failed: %v", err)
		return
	} else {
		t.Log(peinfo.header.DOSHeader.AddressOfNewExeHeader)
		t.Log(peinfo.header.DOSHeader.AddressOfRelocationTable)
		t.Log(peinfo.header.DOSHeader.AddressOfRelocationTable)
		t.Log(*peinfo.sectionname)
	}
}
