package ctxs_test

import (
	"context"
	"fmt"
	"time"

	"github.com/wal1251/pkg/core/ctxs"
)

func ExampleStartMeasureContext() {
	// В реальной жизни оцениваемую функцию можно обернуть в middleware, которая выполнит замер длительности вызова.
	measureSeconds := func(ctx context.Context, target func()) time.Duration {
		ctx = ctxs.StartMeasureContext(ctx)
		target()
		return ctxs.ElapsedFromContext(ctx).Round(time.Second)
	}

	// Замерим
	fmt.Println(measureSeconds(context.TODO(), func() {
		time.Sleep(time.Second)
	}))

	// Output:
	// 1s
}
