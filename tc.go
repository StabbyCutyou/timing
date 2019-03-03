package timing

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// Context represents a context capable of capturing timing infos
type Context interface {
	context.Context
	Start() Context
	Stop()
	Timings() map[string]Record
}

// Record will track metadata about the timings and count of a function
type Record struct {
	Duration  time.Duration `json:"duration"`
	CallCount int           `json:"count"`
}

// tc stands for timing context
type tc struct {
	context.Context
	m       sync.Mutex
	timings map[string]Record
	stack   []frame
	enabled bool //readonly
}

// frame as a name is like, whatever, but i keep it in a stack soooooooo
// anyways, you never touch them so who cares.
type frame struct {
	f string        // funcname
	t time.Time     // current time to start counting duration from on next calc
	d time.Duration // current amount of calculated duration
}

// WithTiming makes a new timing context with timing enabled
func WithTiming(ctx context.Context) Context {
	return &tc{Context: ctx, timings: make(map[string]Record), enabled: true}
}

// WithoutTiming makes a new timing context without timing enabled
// You would use this to create timing contexts that are disabled, and thus
// all Timing related calls are no-ops. Timings() will return nil. Everything else
// behaves off of the underlying context.
func WithoutTiming(ctx context.Context) Context {
	return &tc{Context: ctx}
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

// Start will start recording a new call and cap off the time calculated for the prior call
func (c *tc) Start() Context {
	if !c.enabled {
		return c
	}
	// moving one deeper - cap off timing on prior stack if present
	c.m.Lock()
	defer c.m.Unlock()
	if len(c.stack) > 0 {
		// record how long the prior frame was running
		c.stack[len(c.stack)-1].d += time.Since(c.stack[len(c.stack)-1].t)
	}
	// Add one to the stack
	c.stack = append(c.stack, frame{f: getCallerFuncPC(2), t: time.Now()})
	return c
}

// Stop will stop recording the current call and resume calculating time for the prior call
func (c *tc) Stop() {
	if !c.enabled {
		return
	}
	c.m.Lock()
	defer c.m.Unlock()
	// pop one from the stack
	var f frame
	f, c.stack = c.stack[len(c.stack)-1], c.stack[:len(c.stack)-1]
	if len(c.stack) > 0 {
		// signal to the prior frame that the rest of it's calculations begin now
		c.stack[len(c.stack)-1].t = time.Now()
	}
	r := c.timings[f.f]
	r.CallCount++
	r.Duration += time.Since(f.t) + f.d
	c.timings[f.f] = r
}

// Timings returns the calculated function timings. Ideally, you only call this once all
// timings have finished, as all records won't be 100% complete until the final Stop is called.
func (c *tc) Timings() map[string]Record {
	return c.timings
}
