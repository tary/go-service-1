package crypt

// EncryptData 加密算法
func EncryptData(buf []byte, isClient bool) []byte {
	buflen := len(buf)

	var key []byte
	if isClient {
		key = decryptKey
	} else {
		key = encryptKey
	}
	keylen := len(key)

	for i := 0; i < buflen; i++ {
		n := byte(i%7 + 1)                       //移位长度(1-7)
		b := (buf[i] << n) | (buf[i] >> (8 - n)) // 向左循环移位
		buf[i] = b ^ key[i%keylen]
	}

	return buf
}
