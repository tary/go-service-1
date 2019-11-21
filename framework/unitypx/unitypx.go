package unitypx

/*
#cgo LDFLAGS: -L ../../lib/physx -lunitypx
#include "unitypx.h"
*/
import "C"

var sdk C.unitypx_sdk_t

func init() {
	sdk = C.unitypx_sdk_create()
	if sdk == nil {
		panic("初始化unitypx失败!")
	}
}
