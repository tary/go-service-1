package file

import (
	"os"
)

// PathExists 判断路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
