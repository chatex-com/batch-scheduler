package batch

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSchedule(t *testing.T) {
	Convey("Create empty scheduler", t, func() {
		schedule, err := NewSchedule(nil)

		So(err, ShouldBeNil)
		So(schedule, ShouldNotBeNil)
		So(schedule, ShouldHaveSameTypeAs, &Schedule{})
	})

	Convey("Create simple config", t, func() {
		cfg := []Rule{
			{StartedAt: "00:00", Window: 10 * time.Minute},
			{StartedAt: "10:20", Window: 30 * time.Minute},
			{StartedAt: "20:33", Window: time.Minute},
		}

		pointer, err := NewSchedule(cfg)
		So(err, ShouldBeNil)

		So(pointer, ShouldNotBeNil)
		So(pointer, ShouldHaveSameTypeAs, &Schedule{})

		schedule := *pointer
		So(schedule, ShouldHaveLength, 3)

		So(schedule[0].StartedAt, ShouldEqual, 0)
		So(schedule[0].Window, ShouldEqual, 10*time.Minute)

		So(schedule[1].StartedAt, ShouldEqual, 620)
		So(schedule[1].Window, ShouldEqual, 30*time.Minute)

		So(schedule[2].StartedAt, ShouldEqual, 1233)
		So(schedule[2].Window, ShouldEqual, time.Minute)
	})

	Convey("Create with wrong config", t, func() {
		schedule, err := NewSchedule([]Rule{
			{StartedAt: "33:88", Window: time.Minute},
		})

		So(err, ShouldNotBeNil)
		So(err, ShouldBeError)
		So(schedule, ShouldBeNil)
	})
}

func TestSchedule_GetWindow(t *testing.T) {
	Convey("Empty schedule", t, func() {
		schedule, _ := NewSchedule(nil)

		So(schedule.GetWindow(time.Now()), ShouldEqual, 0)
	})

	Convey("Schedule with one rule", t, func() {
		schedule, _ := NewSchedule([]Rule{
			{StartedAt: "00:00", Window: time.Minute},
		})

		So(schedule.GetWindow(time.Now()), ShouldEqual, time.Minute)
	})

	Convey("Usual flow", t, func() {
		now := time.Now()

		schedule, _ := NewSchedule([]Rule{
			{StartedAt: "00:00", Window: time.Minute},
			{StartedAt: now.Add(-time.Minute).Format("15:04"), Window: time.Second},
		})

		So(schedule.GetWindow(now), ShouldEqual, time.Second)
	})

	Convey("Usual flow with 00:00 configuration", t, func() {
		now := time.Now()

		schedule, _ := NewSchedule([]Rule{
			{StartedAt: now.Add(-time.Hour).Format("15:04"), Window: time.Minute},
			{StartedAt: now.Add(-time.Minute).Format("15:04"), Window: time.Second},
		})

		So(schedule.GetWindow(now), ShouldEqual, time.Second)
	})

	Convey("Next time frame will be earlier than current window duration", t, func() {
		now := time.Now().Round(time.Minute)

		schedule, _ := NewSchedule([]Rule{
			{StartedAt: now.Format("15:04"), Window: 20 * time.Minute},
			{StartedAt: now.Add(10 * time.Minute).Format("15:04"), Window: time.Minute},
		})

		So(schedule.GetWindow(now), ShouldEqual, 10*time.Minute)
	})

	Convey("Next time frame will be after midnight and earlier than current window duration", t, func() {
		t := time.Now().Round(time.Minute)
		midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

		schedule, _ := NewSchedule([]Rule{
			{StartedAt: "00:00", Window: 5 * time.Minute},
			{StartedAt: midnight.Add(-5 * time.Minute).Format("15:04"), Window: 10 * time.Minute},
		})

		So(schedule.GetWindow(midnight.Add(-time.Minute)), ShouldEqual, time.Minute)
	})
}
