package ops

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func readString(reader *bufio.Reader) (string, error) {
	tmp, e := reader.ReadString(' ')
	if e != nil {
		return "", e
	}
	key := strings.TrimSpace(tmp)
	return key, nil
}

func readLen(reader *bufio.Reader) (int, error) {
	tmp, e := reader.ReadString(' ')
	if e != nil {
		return 0, e
	}
	keylen, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		return 0, e
	}
	return keylen, nil
}

// op<key.len>空格<value.len>空格<key><value>
func readKeyAndValue(reader *bufio.Reader) (string, []byte, error) {
	keyLen, e := readLen(reader)
	if e != nil {
		return "", nil, e
	}
	valueLen, e := readLen(reader)
	if e != nil {
		return "", nil, e
	}
	keyVal := make([]byte, keyLen)
	_, e = io.ReadFull(reader, keyVal)
	if e != nil {
		return "", nil, e
	}
	key := string(keyVal)

	values := make([]byte, valueLen)
	_, e = io.ReadFull(reader, values)
	if e != nil {
		return "", nil, e
	}
	//log.Println(fmt.Sprintf("key长度:%d,value长度:%d,key值:%s,value值:%s", klen, vlen, key, vaules))
	return key, values, nil
}
