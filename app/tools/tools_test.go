package tools

import (
	"fmt"
	"testing"
)

func TestGetUUID(t *testing.T) {
	GetUUID()
}

func TestGetUid(t *testing.T) {
	id := GetUid()
	fmt.Printf("id:%d", id)
}
