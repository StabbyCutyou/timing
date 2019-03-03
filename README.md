Timing Context
---------------

This is an attempt to make a lightweight, simple to use drop-in context replacement that you can use to record callstack timings.

You should be able to drop this into any place you already use contexts, and for code you control the signature of, a change from context.Context to timing.Context should be all you need.

All you have to do is add `defer ctx.Start().Stop()` to the beginning of any method call you want to trace the timing of, and that's it!

You can pull all the timing calculations from the context once you're done with it by calling the `Timings()` method, which will return to you a list of fully qualified function names

For an example, checkout `cmd/sample/main.go`

# Note

API potentially unstable until v1.0.0