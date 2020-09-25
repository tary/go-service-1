package sflist

import (
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

	if sl.IsEmpty() {
		t.Error("sl is empty")
	}

	sl.Pop()

	if !sl.IsEmpty() {
		t.Error("sl is not empty")
	}
}

func BenchmarkListParallel(b *testing.B) {
	sl := NewSafeList()
	b.ResetTimer()

	b.RunParallel(func(bp *testing.PB) {
		for bp.Next() {
			sl.Put(12345)
			sl.Pop()
		}
	})

}

func BenchmarkChanParallel(b *testing.B) {
	c := make(chan int, 1000)
	b.ResetTimer()

	b.RunParallel(func(bp *testing.PB) {
		for bp.Next() {
			c <- 12345
			<-c
		}
	})
}
