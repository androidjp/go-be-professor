package cmn

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewMessage(t *testing.T) {
	Convey("should return a new Message", t, func() {
		Convey("given default", func() {
			msg1 := NewMessage("topic-1", []byte{'s'}, WithHeaders(map[string]interface{}{"key1": "val1"}))
			So(msg1.Topic(), ShouldEqual, "topic-1")
		})
	})
}
