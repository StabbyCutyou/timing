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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		traceMyStack(ctx)
	}
}

func BenchmarkWithTimingStart(b *testing.B) {
	ctx := WithTiming(context.Background())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Stop(Start(ctx))
	}
}

func traceMyStack(ctx context.Context) {
	defer Stop(Start(ctx))
}

func TestContextCancel(t *testing.T) {
	ctx, cncl := context.WithCancel(context.Background())
	tctx := WithTiming(ctx)
	cncl()
	<-tctx.Done() // blocks forever if it's not done lol
}

func TestContextValue(t *testing.T) {
	ctx := context.Background()
	vctx := context.WithValue(ctx, struct{}{}, "shanksy")
	tctx := WithTiming(vctx)
	v := tctx.Value(struct{}{})
	if v != "shanksy" {
		t.Fatal("value was not shanksy")
	}
}

func BenchmarkNonTimingContext(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		traceMyStackPure(ctx)
	}
}

func traceMyStackPure(ctx context.Context) {
	defer Stop(Start(ctx))
}
