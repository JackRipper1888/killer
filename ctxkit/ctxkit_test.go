package ctxkit

import (
	"fmt"
	"testing"
	"time"
)

func demo() {
	cxt, _ := CtxAdd()
	for {
		select {
		case <-time.After(1 * time.Second):
		case <-cxt.Done():
			fmt.Println(" demo stop ")
			return
		}
	}
}
func TestGetCtx(t *testing.T) {
	go demo()
	go demo()
	go demo()
	time.Sleep(3 * time.Second)
	CancelAll()
	time.Sleep(10 * time.Second)
	fmt.Println("stop")
}
