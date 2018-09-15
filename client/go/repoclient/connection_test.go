package repoclient

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const globalTestURL = "http://localhost:8000/"

func TestConnect(t *testing.T) {
	Convey("Should be able to connect...", t, func() {

		_, err := Connect(globalTestURL)
		So(err, ShouldBeNil)

	})

}
