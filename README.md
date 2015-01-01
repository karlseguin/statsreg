# statsreg
A collector and writer of simple statistics.

We use it, for example, to generate a JSON file that tells us things such as as how often we miss a fetch from a pool or have to grow a supposedly fixed-length buffer.

```go
// Run every minute.
// Save output to stats.json, truncate the file each time
func init() {
  config := statsreg.Configure().
    Frequency(time.Minute).
    File("stats.json", true)
  sr := statsreg.New(config)

  // bp.Misses is a func() int64 that records how often a get
  // hit a drained pool
  sr.RegisterInt64("main_pool_misses", bp.Misses)

  // note that both debug.ReadGCStats() and runtime.NumGoroutine()
  // require fairly exclusive locks and, under load, you might
  // want to consider whether or not they should be on at all times
  // (for the record runtime.ReadMemStats() is much more intrusive)
  sr.RegisterInt64("gcs", func() int64 {
    stats := new(debug.GCStats) // import "runtime/debug"
    debug.ReadGCStats(stats)
    return stats.NumGC
  })
  sr.RegisterInt64("goroutines", func() int64 {
    return int64(runtime.NumGoroutine()) // import "runtime"
  })
}
```

Statsreg is thread-safe. Providers can be added at any time
