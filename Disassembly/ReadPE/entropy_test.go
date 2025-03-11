package readpe

import (
	"fmt"
	"testing"
)

func TestCalculateEntropy(t *testing.T) {
	filePath := "./test.exe"
	cal, err := CalculateEntropy(filePath)
	if err != nil {
		t.Errorf("CalculateEntropy failed: %v", err)
		return
	} else {
		fmt.Println("answer is ", cal)
	}

}
