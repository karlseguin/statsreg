package statsreg

import (
	"time"
)

type Configuration struct {
	frequency time.Duration
	output    Output
}

// Create a default configuration object, which can be used as-is
// or configured to suit your needs.
func Configure() *Configuration {
	return &Configuration{
		frequency: time.Minute,
		output:    &File{"stats.json", true},
	}
}

// How frequently the registry will collect stats (default: 1 minute)
func (c *Configuration) Frequency(f time.Duration) *Configuration {
	c.frequency = f
	return c
}

// Where to send the output. Consider using the File configuration
// method: default(File("stats.json", true))
func (c *Configuration) Output(output Output) *Configuration {
	c.output = output
	return c
}

// The file to send the output to and whether the content should be
// truncated first
func (c *Configuration) File(path string, truncate bool) *Configuration {
	c.output = &File{path, truncate}
	return c
}
