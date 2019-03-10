Timing Context
---------------

This is an attempt to make a lightweight, simple to use drop-in context replacement that you can use to record callstack timings.

You should be able to drop this into any place you already use contexts, and for code you control the signature of, a change from context.Context to timing.Context should be all you need.

All you have to do is add `defer timing.Stop(timing.Start(ctx))` to the beginning of any method call you want to trace the timing of, and initialize any context once with timing.WithTiming(context.Context), and that's it!

You can pull all the timing calculations from the context once you're done with it by calling the `Timings()` method, which will return to you a list of fully qualified function names

For an example, checkout `cmd/sample/main.go`

# Note

API potentially unstable until v1.0.0

Not sure if this is even a good/correct approach? (needs to be good enough, anyways, but not perfect)

# Roadmap
A way to merge a sub-timing-context so you can safely track sub go routines you kick off, and collect the info underneath the original timing table to help paint full pictures. 