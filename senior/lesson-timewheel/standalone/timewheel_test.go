package standalone_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"lesson-timewheel/standalone"
	"testing"
	"time"
)

func TestTimeWheel_AddTask(t *testing.T) {
	Convey("单机版时间轮测试用例", t, func() {
		tw := standalone.NewTimeWheel(10, 500*time.Millisecond)
		defer tw.Stop()

		tw.AddTask("test1", func() {
			t.Errorf("test1, %v", time.Now())
		}, time.Now().Add(time.Second))

		tw.AddTask("test2", func() {
			t.Errorf("test2, %v", time.Now())
		}, time.Now().Add(5*time.Second))

		tw.AddTask("test3", func() {
			t.Errorf("test3, %v", time.Now())
		}, time.Now().Add(3*time.Second))

		<-time.After(6 * time.Second)
	})
}
