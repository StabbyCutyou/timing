package timing

import (
	"context"
	"testing"
)

func getCallerPC() string {
	return getCallerFuncPC(2)
}

func BenchmarkPC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCallerPC()
	}
}

func BenchmarkWithTiming(b *testing.B) {
	ctx := WithTiming(context.Background())
	for i := 0; i < b.N; i++ {
		traceMyStack(ctx)
	}
}

func BenchmarkWithoutTiming(b *testing.B) {
	ctx := WithoutTiming(context.Background())
	for i := 0; i < b.N; i++ {
		traceMyStack(ctx)
	}
}

func traceMyStack(ctx Context) {
	defer ctx.Start().Stop()
}

func TestWithoutTiming(t *testing.T) {
	ctx := WithoutTiming(context.Background())
	ctx.Start()
	ctx.Stop()
	if ctx.Timings() != nil {
		t.Fatal("timing dict should have been nil, is not")
	}
}
