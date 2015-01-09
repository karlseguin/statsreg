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
  sr.RegisterGeneric("goroutines", func() interface{} {
    return runtime.NumGoroutine() // import "runtime"
  })
}
```

## Functions
- `RegisterInt64(name string, provider func() int64)` - registers a provider that returns an int64
- `RegisterString(name string, provider func() string)` - registers a provider that returns an string
- `RegisterGeneric(name string, provider func() interface{})` - registers a provider that returns an interface{}. The value must be serializable using encoding/json
- `Remove(name string)` - removes the provider. Safe to call if it doesn't exist

Methods are thread-safe.

### Overwriting
By default, statsreg will log a message when registering an alread-registered name. This makes sense for static statistics (those registered at init). However, you may want to dynamically register statistics, sometimes with the same name. You can disable the logging by setting the global var `LogOvewrite` to `nil` (this is actually a function, so you could implement your own custom logging logic, but why?).
