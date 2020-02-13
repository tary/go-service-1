package safecontainer

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

type msgFireInfo struct {
	name    string
	content interface{}
}

type T struct {
	index int
}

func Test_list(t *testing.T) {

	sl := NewSafeList()

	sl.Put(&msgFireInfo{
		name:    "test",
		content: &T{index: 1},
	})

	fmt.Println(sl.IsEmpty())

	sl.Pop()

	fmt.Println(sl.IsEmpty())

}


func Test_mulist(t *testing.T) {

	sl := NewSafeList_M()
	assert.Equal(t,sl.IsEmpty(),true,"must empty")
	sl.Put(1)
	assert.Equal(t,sl.IsEmpty(),false,"must not empty")
	sl.Put(2)
	sl.Put(3)

	v1,err1:= sl.Pop()
	assert.Equal(t,v1,1,"pop1 error")
	assert.Equal(t,err1,nil,"pop1 error")

	v2,err2:= sl.Pop()
	assert.Equal(t,v2,2,"pop2 error")
	assert.Equal(t,err2,nil,"pop2 error")

	v3,err3:= sl.Pop()
	assert.Equal(t,v3,3,"pop3 error")
	assert.Equal(t,err3,nil,"pop3 error")

}

//atomic list
func BenchmarkList(b *testing.B) {
	sl := NewSafeList()

	b.ResetTimer()

	//for i := 0; i < b.N; i++ {
	//	sl.Put(1)
	//	sl.Pop()
	//}

	b.RunParallel(func(bp *testing.PB) {
		for bp.Next() {
			sl.Put(1)
			sl.Pop()
		}
	})

}


//mutex list
func BenchmarkMuList(b *testing.B) {
	sl := NewSafeList_M()

	b.ResetTimer()

	//for i := 0; i < b.N; i++ {
	//	sl.Put(1)
	//	sl.Pop()
	//}
	b.RunParallel(func(bp *testing.PB) {
		for bp.Next() {
			sl.Put(1)
			sl.Pop()
		}
	})

}


//chan
func BenchmarkChan(b *testing.B) {
	c := make(chan int, 1000)

	b.ResetTimer()

	b.RunParallel(func(bp *testing.PB) {
		for bp.Next() {
			c <- 1
			<-c
		}
	})
}
