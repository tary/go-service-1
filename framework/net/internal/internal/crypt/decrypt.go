package crypt

// DecryptData 解密算法
func DecryptData(buf []byte, isClient bool) []byte {
	buflen := len(buf)
	var key []byte
	if isClient {
		key = encryptKey
	} else {
		key = decryptKey
	}
	keylen := len(key)

	for i := 0; i < buflen; i++ {
		b := buf[i] ^ key[i%keylen]
		n := byte(i%7 + 1)                 //移位长度(1-7)
		buf[i] = (b >> n) | (b << (8 - n)) // 向右循环移位
	}
	return buf
}
