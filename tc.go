package timing

import (
	"context"
	"runtime"
	"sync"
	"time"
)

type key string

const keyTiming key = "$timing.tc"
const keyLock key = "$timing.mu"

// Record will track metadata about the timings and count of a function
type Record struct {
	Duration  time.Duration `json:"duration"`
	CallCount int           `json:"count"`
}

// tc stands for timing context
type tc struct {
	timings map[string]Record
	stack   []frame
}

// frame as a name is like, whatever, but i keep it in a stack soooooooo
// anyways, you never touch them so who cares.
type frame struct {
	f string        // funcname
	t time.Time     // current time to start counting duration from on next calc
	d time.Duration // current amount of calculated duration
}

// adapted this approach from the following SO link, the answer from user svenwltr, Jul 24 '16 at 11:06
// https://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang
func getCallerFuncPC(stack int) string {
	pc, _, _, ok := runtime.Caller(stack)
	if d := runtime.FuncForPC(pc); ok && d != nil {
		return d.Name()
	}
	return "" // TODO better default value?
}

// WithTiming initializes a context ready to track callstack timings within a single goroutine
func WithTiming(ctx context.Context) context.Context {
	return context.WithValue(
		context.WithValue(
			ctx,
			keyTiming,
			&tc{stack: make([]frame, 0), timings: make(map[string]Record)},
		),
		keyLock,
		&sync.Mutex{},
	)
}

// Timings returns the timing collections up to this point. Note that until Stop() is
// evaluated, the data is incomplete and unsafe to use.
func Timings(ctx context.Context) map[string]Record {
	tctx, ok := ctx.Value(keyTiming).(*tc)
	if !ok {
		return nil
	}
	return tctx.timings
}

// Start will begin tracking time for a callstack frame. It is meant to be called
// by passing it into Stop()
func Start(ctx context.Context) context.Context {
	mu, ok := ctx.Value(keyLock).(*sync.Mutex)
	if !ok {
		return ctx
	}
	mu.Lock()
	defer mu.Unlock()
	t := ctx.Value(keyTiming)
	tctx, ok := t.(*tc)
	if !ok {
		return ctx
	}
	if len(tctx.stack) > 0 {
		tctx.stack[len(tctx.stack)-1].d += time.Since(tctx.stack[len(tctx.stack)-1].t)
	}
	tctx.stack = append(tctx.stack, frame{f: getCallerFuncPC(2), t: time.Now()})
	return context.WithValue(ctx, keyTiming, tctx)
}

// Stop finishes calculating the timing for a callstack frame. It is meant to be called
// by passing in Start() as the context, and defered as early as possible in a function
func Stop(ctx context.Context) {
	mu, ok := ctx.Value(keyLock).(*sync.Mutex)
	if !ok {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	t := ctx.Value(keyTiming)
	tctx, ok := t.(*tc)
	if !ok {
		return
	}
	var f frame
	f, tctx.stack = tctx.stack[len(tctx.stack)-1], tctx.stack[:len(tctx.stack)-1]
	if len(tctx.stack) > 0 {
		// signal to the prior frame that the rest of it's calculations begin now
		tctx.stack[len(tctx.stack)-1].t = time.Now()
	}
	r := tctx.timings[f.f]
	r.CallCount++
	r.Duration += time.Since(f.t) + f.d
	tctx.timings[f.f] = r
}
