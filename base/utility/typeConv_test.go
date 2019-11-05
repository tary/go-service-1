package utility

import (
	"fmt"
	"testing"

	log "github.com/cihub/seelog"
)

func TestConvert(t *testing.T) {
	var data string

	floatV := float32(3.14)
	data = TypeToString(&floatV)
	log.Debug(data)

	float64V := float64(3.14)
	data = TypeToString(&float64V)
	log.Debug(data)

	intV := int32(0)
	data = TypeToString(intV)
	log.Debug(data)
	data = TypeToString(&intV)
	log.Debug(data)

}

func TestUnquote(t *testing.T) {
	str := "test"
	str = Unquote(str)
	fmt.Println(str)

	str = "\"haha\""
	str = Unquote(str)
	fmt.Println(str)

	str = "\"haha"
	str = Unquote(str)
	fmt.Println(str)
}
