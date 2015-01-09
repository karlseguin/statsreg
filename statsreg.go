package statsreg

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

var LogOverwrite = func(name string) {
	log.Println(fmt.Sprintf("stats %s is already registered, overwriting", name))
}

type Int64Provider func() int64
type StringProvider func() string
type GenericProvider func() interface{}

type StatsReg struct {
	*Configuration
	sync.Mutex
	stats     map[string]interface{}
	providers map[string]GenericProvider
}

// Create a new registry with the specified configuration
func New(config *Configuration) *StatsReg {
	sr := &StatsReg{
		Configuration: config,
		providers:     make(map[string]GenericProvider),
		stats:         make(map[string]interface{}),
	}
	go sr.work()
	return sr
}

// Registers a provider which exports a statistic as an int64
func (sr *StatsReg) RegisterInt64(name string, provider Int64Provider) {
	sr.RegisterGeneric(name, func() interface{} { return provider() })
}

// Registers a provider which exports a statistic as string
func (sr *StatsReg) RegisterString(name string, provider StringProvider) {
	sr.RegisterGeneric(name, func() interface{} { return provider() })
}

// Registers a provider which exports a statistic as anything (must be serializable
// with encoding/json)
func (sr *StatsReg) RegisterGeneric(name string, provider GenericProvider) {
	sr.Lock()
	defer sr.Unlock()
	if _, exists := sr.providers[name]; exists && LogOverwrite != nil {
		LogOverwrite(name)
	}
	sr.providers[name] = provider
}

func (sr *StatsReg) Remove(name string) {
	sr.Lock()
	defer sr.Unlock()
	delete(sr.providers, name)
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
	for name, provider := range sr.providers {
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
