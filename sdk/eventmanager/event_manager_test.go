package eventmanager

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldGetLastEvents(t *testing.T) {
	Convey("should get last events from event manager", t, func() {
		evts, err := LastEvents("", "", "1d")
		So(err, ShouldBeNil)
		So(len(evts), ShouldBeGreaterThan, 0)
	})

}
