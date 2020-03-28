package batch

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v2"
)

const configYaml = `
- started_at: 00:00
  window: 1m
- started_at: 10:50
  window: 10m
- started_at: 18:30
  window: 30m
- started_at: 22:00
  window: 10s
`

func TestRuleUnmarshal(t *testing.T) {
	Convey("Test load config rules from yaml", t, func() {
		var cfg []Rule
		err := yaml.Unmarshal([]byte(configYaml), &cfg)

		So(err, ShouldBeNil)
		So(cfg, ShouldHaveLength, 4)

		So(cfg[0].StartedAt, ShouldEqual, "00:00")
		So(cfg[0].Window, ShouldEqual, time.Minute)

		So(cfg[1].StartedAt, ShouldEqual, "10:50")
		So(cfg[1].Window, ShouldEqual, 10*time.Minute)

		So(cfg[2].StartedAt, ShouldEqual, "18:30")
		So(cfg[2].Window, ShouldEqual, 30*time.Minute)

		So(cfg[3].StartedAt, ShouldEqual, "22:00")
		So(cfg[3].Window, ShouldEqual, 10*time.Second)
	})
}
