package main

import (
	"context"
	"time"

	"github.com/StabbyCutyou/timing"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	ctx := timing.New(context.Background())
	FunctionA(ctx)
	spew.Dump(timing.Timings(ctx))
	ctx = timing.New(context.Background())
	FunctionSlowChain(ctx)
	spew.Dump(timing.Timings(ctx))
}

// FunctionA is a dummy
func FunctionA(ctx context.Context) {
	defer timing.Stop(timing.Start(ctx))
	time.Sleep(5 * time.Millisecond)
	FunctionB(ctx)
}

// FunctionB is a dummy
func FunctionB(ctx context.Context) {
	defer timing.Stop(timing.Start(ctx))
	time.Sleep(500 * time.Millisecond)
	FunctionC(ctx)
}

// FunctionC is a dummy
func FunctionC(ctx context.Context) {
	defer timing.Stop(timing.Start(ctx))
	time.Sleep(50 * time.Millisecond)
}

func FunctionSlowChain(ctx context.Context) {
	defer timing.Stop(timing.Start(ctx))
	time.Sleep(1 * time.Millisecond)
	link1(ctx)
}

func link1(ctx context.Context) {
	defer timing.Stop(timing.Start(ctx))
	time.Sleep(10 * time.Millisecond)
	link2(ctx)
}

func link2(ctx context.Context) {
	defer timing.Stop(timing.Start(ctx))
	time.Sleep(100 * time.Millisecond)
	link3(ctx)
	link3(ctx)
	link3(ctx)
}

func link3(ctx context.Context) {
	defer timing.Stop(timing.Start(ctx))
	time.Sleep(500 * time.Millisecond)
	link4(ctx)
	link4(ctx)
}

func link4(ctx context.Context) {
	defer timing.Stop(timing.Start(ctx))
	time.Sleep(1000 * time.Millisecond)
}