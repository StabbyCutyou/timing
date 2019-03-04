package timing

import (
	"context"
	"sync"
	"time"
)

type key string

const keyTiming key = "$timing.tc"
const keyLock key = "$timing.mu"

func New(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, keyTiming, &tc{stack: make([]frame, 0), timings: make(map[string]Record)})
	return context.WithValue(ctx, keyLock, &sync.Mutex{})
}

func Timings(ctx context.Context) map[string]Record {
	t := ctx.Value(keyTiming)
	tctx, ok := t.(*tc)
	if !ok {
		return nil
	}
	return tctx.timings
}

func Start(ctx context.Context) context.Context {
	m := ctx.Value(keyLock)
	mu, ok := m.(*sync.Mutex)
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

func Stop(ctx context.Context) {
	m := ctx.Value(keyLock)
	mu, ok := m.(*sync.Mutex)
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
