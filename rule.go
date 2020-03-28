package batch

import (
	"time"
)

type Rule struct {
	StartedAt string        `yaml:"started_at"`
	Window    time.Duration `yaml:"window"`
}
