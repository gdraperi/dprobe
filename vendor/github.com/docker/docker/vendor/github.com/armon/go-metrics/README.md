go-metrics
==========

This library provides a `metrics` package which can be used to instrument code,
expose application metrics, and profile runtime performance in a flexible manner.

Sinks
=====

The `metrics` package makes use of a `MetricSink` interface to support delivery
to any type of backend. Currently the following sinks are provided:

* StatsiteSink : Sinks to a statsite instance (TCP)
* StatsdSink: Sinks to a statsd / statsite instance (UDP)
* InmemSink : Provides in-memory aggregation, can be used to export stats
* FanoutSink : Sinks to multiple sinks. Enables writing to multiple statsite instances for example.
* BlackholeSink : Sinks to nowhere

In addition to the sinks, the `InmemSignal` can be used to catch a signal,
and dump a formatted output of recent metrics. For example, when a process gets
a SIGUSR1, it can dump to stderr recent performance metrics for debugging.

Examples
========

Here is an example of using the package:

    func SlowMethod() ***REMOVED***
        // Profiling the runtime of a method
        defer metrics.MeasureSince([]string***REMOVED***"SlowMethod"***REMOVED***, time.Now())
***REMOVED***

    // Configure a statsite sink as the global metrics sink
    sink, _ := metrics.NewStatsiteSink("statsite:8125")
    metrics.NewGlobal(metrics.DefaultConfig("service-name"), sink)

    // Emit a Key/Value pair
    metrics.EmitKey([]string***REMOVED***"questions", "meaning of life"***REMOVED***, 42)


Here is an example of setting up an signal handler:

    // Setup the inmem sink and signal handler
    inm := NewInmemSink(10*time.Second, time.Minute)
    sig := DefaultInmemSignal(inm)
    metrics.NewGlobal(metrics.DefaultConfig("service-name"), inm)

    // Run some code
    inm.SetGauge([]string***REMOVED***"foo"***REMOVED***, 42)
	inm.EmitKey([]string***REMOVED***"bar"***REMOVED***, 30)

	inm.IncrCounter([]string***REMOVED***"baz"***REMOVED***, 42)
	inm.IncrCounter([]string***REMOVED***"baz"***REMOVED***, 1)
	inm.IncrCounter([]string***REMOVED***"baz"***REMOVED***, 80)

	inm.AddSample([]string***REMOVED***"method", "wow"***REMOVED***, 42)
	inm.AddSample([]string***REMOVED***"method", "wow"***REMOVED***, 100)
	inm.AddSample([]string***REMOVED***"method", "wow"***REMOVED***, 22)

    ....

When a signal comes in, output like the following will be dumped to stderr:

    [2014-01-28 14:57:33.04 -0800 PST][G] 'foo': 42.000
    [2014-01-28 14:57:33.04 -0800 PST][P] 'bar': 30.000
	[2014-01-28 14:57:33.04 -0800 PST][C] 'baz': Count: 3 Min: 1.000 Mean: 41.000 Max: 80.000 Stddev: 39.509
    [2014-01-28 14:57:33.04 -0800 PST][S] 'method.wow': Count: 3 Min: 22.000 Mean: 54.667 Max: 100.000 Stddev: 40.513

