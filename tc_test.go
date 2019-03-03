package timing

import (
	"context"
	"testing"
)

func BenchmarkPC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCallerPC()
	}
}

func BenchmarkTrace(b *testing.B) {
	ctx := NewContext(context.Background())
	for i := 0; i < b.N; i++ {
		traceMyStack(ctx)
	}
}

func getCallerPC() string {
	return getCallerFuncPC(2)
}

func traceMyStack(ctx Context) {
	defer ctx.Start().Stop()
}
