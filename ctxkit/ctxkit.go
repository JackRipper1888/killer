package ctxkit

import (
	"context"
	"fmt"
	"runtime"
)

var (
	contextList = []*contextData{}
	parentCtx   = context.Background()
)

type contextData struct {
	ctx    context.Context
	cancel func()
}

func GetCtx(num int) context.Context {
	fmt.Println(num)
	if num < 0 || num > len(contextList) {
		panic("num out of range")
	}
	return contextList[num].ctx
}

func InitCtxPool(txNum int) func() {
	i := 0
	for i <= txNum {
		i++
		Ctx, Cancel := context.WithCancel(parentCtx)
		contextList = append(contextList, &contextData{
			ctx:    Ctx,
			cancel: Cancel,
		})
	}
	return CancelAll
}

func CtxAdd() (context.Context, func()) {
	Ctx, Cancel := context.WithCancel(parentCtx)
	contextList = append(contextList, &contextData{
		ctx:    Ctx,
		cancel: Cancel,
	})
	return Ctx, Cancel
}

func CancelAll() {
	defer runtime.GC()
	for _, k := range contextList {
		k.cancel()
	}
}
