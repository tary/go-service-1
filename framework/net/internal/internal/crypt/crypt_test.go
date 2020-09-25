package crypt

import (
	"testing"
)

func Test_crypt(t *testing.T) {
	ed := EncryptData([]byte("aaaaaaa"), false)
	dd := DecryptData(ed, true)

	_ = dd

	//fmt.Debug("ed: ", ed, ", dd: ", dd)
}
