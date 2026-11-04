package base64

import (
	"fmt"
	"testing"
)

func TestSha512(t *testing.T) {
	tmp, _ := Decode([]byte("1w=="))
	fmt.Println(tmp)
}
