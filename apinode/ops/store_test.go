package ops

import (
	"github.com/klauspost/reedsolomon"
	"runtime"
	"testing"
)

func TestReedsolomon(t *testing.T) {
	enc, _ := reedsolomon.New(4, 2, reedsolomon.WithMaxGoroutines(runtime.NumCPU()))
	bigfile := "123456789ABCDEF"
	bigfile1 := []byte(bigfile)
	// Split the file
	split, _ := enc.Split(bigfile1)
	t.Log("split->", len(split))
	for i, v := range split {
		t.Log(i, string(v))
	}
	err := enc.Encode(split)
	t.Log(err == nil)
	ok, err1 := enc.Verify(split)
	t.Log(err1, "->", ok)
	t.Log("破坏数据：")
	split[0] = nil
	split[1] = nil
	for i, v := range split[:4] {
		t.Log("数据分片", i, string(v))
	}
	t.Log("还原数据：")
	enc.Reconstruct(split)
	data := make([]byte, 0)
	for i, v := range split[:4] {
		t.Log(i, string(v))
		data = append(data, v...)
	}
	t.Log("结果:", string(data))
}
