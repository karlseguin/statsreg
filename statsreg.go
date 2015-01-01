package statsreg

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Int64Provider func() int64
type StringProvider func() string

type StatsReg struct {
	*Configuration
	sync.Mutex
	stats   map[string]interface{}
	int64s  map[string]Int64Provider
	strings map[string]StringProvider
}

// Create a new registry with the specified configuration
func New(config *Configuration) *StatsReg {
	sr := &StatsReg{
		Configuration: config,
		int64s:        make(map[string]Int64Provider),
		strings:       make(map[string]StringProvider),
		stats:         make(map[string]interface{}),
	}
	go sr.work()
	return sr
}

func (sr *StatsReg) RegisterInt64(name string, provider Int64Provider) {
	sr.Lock()
	defer sr.Unlock()
	sr.int64s[name] = provider
}

func (sr *StatsReg) RegisterString(name string, provider StringProvider) {
	sr.Lock()
	defer sr.Unlock()
	sr.strings[name] = provider
}

func (sr *StatsReg) work() {
	for {
		time.Sleep(sr.frequency)
		sr.Collect()
	}
}

// Force collection to run (this normally runs automatically in it's own
// goroutine at the configured frequency).
func (sr *StatsReg) Collect() {
	sr.Lock()
	defer sr.Unlock()
	for name, provider := range sr.int64s {
		sr.stats[name] = provider()
	}
	for name, provider := range sr.strings {
		sr.stats[name] = provider()
	}
	data, err := json.Marshal(sr.stats)
	if err != nil {
		log.Println("statsreg failed to serialize stats", err)
		return
	}
	if err := sr.output.Write(data); err != nil {
		log.Println("statsreg failed to write data", err)
	}
}
