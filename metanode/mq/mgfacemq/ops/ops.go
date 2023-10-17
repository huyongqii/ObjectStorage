package ops

import (
	"bufio"

	"com.mgface.disobj/metanode/store"
)

// Get abnf规则]G<key.len> <key>
func Get(store store.Store, reader *bufio.Reader) (interface{}, error) {
	key, value, _ := readKeyAndValue(reader)
	return store.Get(key, value)
}

func Set(store store.Store, reader *bufio.Reader) error {
	key, value, e := readKeyAndValue(reader)
	store.Set(key, value)
	return e
}

func Del(store store.Store, reader *bufio.Reader) error {
	key, value, e := readKeyAndValue(reader)
	store.Del(key, value)
	return e
}

func Put(store store.Store, reader *bufio.Reader) (interface{}, error) {
	fname, _ := readString(reader)

	fsize, _ := readLen(reader)
	return store.Put(fname, fsize, reader)
}

func Syn(store store.Store, reader *bufio.Reader) (interface{}, error) {
	synLen, _ := readLen(reader)
	return store.Sync(synLen, reader)
}
