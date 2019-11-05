package dailytimer

import (
	"fmt"
	"testing"
	//	"time"
)

type T struct {
}

func (t *T) test1() {
	fmt.Println("xxx")
}

func (t *T) test2(a int) {
	fmt.Println("xxx ", a)
}

func TestAddTimer(t *testing.T) {
	/*t1 := &T{}
	tw := New("1s", 60)
	tw.Start()
	//tw.AddTimerTask("2018-12-28 17:45:00", time.Duration(2)*time.Second, "1", false, t1.test1)
	tw.AddTimerTask("2018-12-28 17:46:00", time.Duration(1)*time.Second, "2", true, t1.test2, 1)
	for {

	}*/
}
