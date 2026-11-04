package BootConfig

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBootConfig(t *testing.T) {
	var obj Root
	js, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(js))
	var obj_out Root
	err = json.Unmarshal(js, &obj_out)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkBootConfig(b *testing.B) {
	for i := 0; i < 100; i++ {
		var obj Root
		js, err := json.MarshalIndent(obj, "", "\t")
		if err != nil {
			b.Error(err)
		}
		var obj_out Root
		err = json.Unmarshal(js, &obj_out)
		if err != nil {
			b.Error(err)
		}
	}
}
