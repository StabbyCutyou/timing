package timing

import (
	"context"
	"time"
)

type key string

const keyTiming key = "$timing.tc"

func New(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyTiming, &tc{stack: make([]frame, 0), timings: make(map[string]Record)})
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
