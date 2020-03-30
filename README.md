# Batch Scheduler

<a href="https://opensource.org/licenses/Apache-2.0" rel="nofollow"><img src="https://img.shields.io/badge/license-Apache%202-blue" alt="License" style="max-width:100%;"></a>
![unit-tests](https://github.com/chatex-com/batch-scheduler/workflows/unit-tests/badge.svg)
![linter](https://github.com/chatex-com/batch-scheduler/workflows/linter/badge.svg)

## Example

```go
package main

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/chatex-com/batch-scheduler"
)

// trigger events:
// - from 00:00 every minute
// - from 10:00 every 10 minutes
// - from 18:00 every 20 seconds
// - from 22:00 every hour
// NOTE: First event will be triggered immediately
const configYaml = `
- started_at: 00:00
  window: 1m
- started_at: 10:00
  window: 10m
- started_at: 18:00
  window: 20s
- started_at: 22:00
  window: 60m
`

func main() {
	var config []batch.Rule
	_ = yaml.Unmarshal([]byte(configYaml), &config)

	scheduler, err := batch.NewScheduler(config)
	if err != nil {
		panic(err)
	}

	ch := scheduler.Run()
	// scheduler.Stop()
	for range ch {
		fmt.Println(time.Now())
	}
}
```
