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

type tc struct {
	context.Context
	enabled bool
	timings map[string]Record
	stack   []frame
	m       sync.Mutex
}

type frame struct {
	f string
	t time.Time
	d time.Duration
}

// NewContext makes a new timing context
func NewContext(ctx context.Context) Context {
	return &tc{Context: ctx, timings: make(map[string]Record), enabled: true}
}

func getCallerFuncPC(stack int) string {
	pc, _, _, ok := runtime.Caller(stack)
	if d := runtime.FuncForPC(pc); ok && d != nil {
		return d.Name()
	}
	return "" // TODO better default value?
}

// Start will start recording a new call and cap off the time calculated for the prior call
func (c *tc) Start() Context {
	// moving one deeper - cap off timing on prior stack if present
	c.m.Lock()
	defer c.m.Unlock()
	// Add one to the stack
	if len(c.stack) > 0 {
		// record how long the prior frame was running
		c.stack[len(c.stack)-1].d = time.Since(c.stack[len(c.stack)-1].t)
	}
	c.stack = append(c.stack, frame{f: getCallerFuncPC(2), t: time.Now()})
	return c
}

// Stop will stop recording the current call and resume calculating time for the prior call
func (c *tc) Stop() {
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

// Timings returns the calculated function timings
func (c *tc) Timings() map[string]Record {
	return c.timings
}
