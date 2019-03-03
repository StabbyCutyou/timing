package main

import (
	"context"
	"time"

	timing "github.com/StabbyCutyou/timingcontext"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	ctx := timing.NewContext(context.Background())
	FunctionA(ctx)
	spew.Dump(ctx.Timings())

}

// FunctionA is a dummy
func FunctionA(ctx timing.Context) {
	defer ctx.Start().Stop()
	time.Sleep(5 * time.Millisecond)
	FunctionB(ctx)
}

// FunctionB is a dummy
func FunctionB(ctx timing.Context) {
	defer ctx.Start().Stop()
	time.Sleep(500 * time.Millisecond)
	FunctionC(ctx)
}

// FunctionC is a dummy
func FunctionC(ctx timing.Context) {
	defer ctx.Start().Stop()
	time.Sleep(50 * time.Millisecond)
}
