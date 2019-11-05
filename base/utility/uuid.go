package utility

import (
	"fmt"
	"math/rand"
	"strconv"
)

// NewUUID newuuid
func NewUUID() uint64 {
	num := rand.Uint64()%(0xefffffffffffffff) + 0x1000000000000000
	return num
}

// UUIDToString  UUIDToString
func UUIDToString(num uint64) string {
	result := fmt.Sprintf("%x", num)
	return result
}

// UUIDFromString UUIDFromString
func UUIDFromString(tokenStr string) (uint64, error) {

	num, err := strconv.ParseUint(tokenStr, 16, 64)
	return num, err
}
