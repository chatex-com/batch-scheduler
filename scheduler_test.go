package batch

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewScheduler(t *testing.T) {
	Convey("Create empty scheduler", t, func() {
		scheduler, err := NewScheduler(nil)

		So(err, ShouldBeNil)
		So(scheduler, ShouldNotBeNil)
		So(scheduler, ShouldHaveSameTypeAs, &Scheduler{})
	})

	Convey("Create simple config", t, func() {
		scheduler, err := NewScheduler([]Rule{
			{StartedAt: "00:00", Window: 10 * time.Minute},
			{StartedAt: "10:20", Window: 30 * time.Minute},
			{StartedAt: "20:33", Window: time.Minute},
		})

		So(err, ShouldBeNil)

		So(scheduler, ShouldNotBeNil)
		So(scheduler, ShouldHaveSameTypeAs, &Scheduler{})
		So(scheduler.schedule, ShouldHaveLength, 3)
	})

	Convey("Create with wrong config", t, func() {
		scheduler, err := NewScheduler([]Rule{
			{StartedAt: "33:77", Window: time.Minute},
		})

		So(err, ShouldNotBeNil)
		So(err, ShouldBeError)
		So(scheduler, ShouldBeNil)
	})
}

func TestScheduler_Run(t *testing.T) {
	Convey("Test infinity run without config", t, func() {
		scheduler, _ := NewScheduler(nil)

		ch := scheduler.Run()

		t1 := time.Now()
		<-ch
		t2 := time.Now()
		<-ch
		t3 := time.Now()

		So(t2.Sub(t1), ShouldAlmostEqual, 0, time.Millisecond)
		So(t3.Sub(t2), ShouldAlmostEqual, 0, time.Millisecond)
	})

	Convey("Test run with config", t, func() {
		scheduler, _ := NewScheduler([]Rule{
			{StartedAt: "00:00", Window: 100 * time.Millisecond},
		})

		ch := scheduler.Run()

		t1 := time.Now()
		<-ch
		t2 := time.Now()
		<-ch
		t3 := time.Now()

		So(t2.Sub(t1), ShouldAlmostEqual, 0, time.Millisecond)
		So(t3.Sub(t2), ShouldAlmostEqual, 100 * time.Millisecond, 5 * time.Millisecond)
	})
}

func TestScheduler_Stop(t *testing.T) {
	Convey("Test stop with config", t, func() {
		scheduler, _ := NewScheduler([]Rule{
			{StartedAt: "00:00", Window: 20 * time.Millisecond},
		})

		ch := scheduler.Run()

		go func() {
			<-time.After(100 * time.Millisecond)
			scheduler.Stop()
		}()

		var count int
		t1 := time.Now()
		for range ch {
			count++
		}
		t2 := time.Now()

		So(t2.Sub(t1), ShouldAlmostEqual, 100 * time.Millisecond, 5 * time.Millisecond)
		So(count, ShouldAlmostEqual, 5, 1)
	})
}
