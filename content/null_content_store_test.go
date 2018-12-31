package content

import (
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

// This test file checks that the NullContentStore conforms to the interfaces.ContentStore
// interface.

func TestNullContentStore_ConformsToInterface(t *testing.T) {

	Convey("Null content store should conform to the interface...", t, func() {

		var contentStore interfaces.ContentStore

		Convey("Should be able to create it...", func() {
			contentStore = new(NullContentStore)
		})

	})

}

func TestNullContentStore_AppendBlobContent(t *testing.T) {
	Convey("NullContentStore should do nothing when appending something...", t, func() {
		var contentStore interfaces.ContentStore = new(NullContentStore)
		var record models.Blob

		written, err := contentStore.WriteBlobContent(&record, strings.NewReader("Hello world"))
		So(written, ShouldBeNil)
		So(err, ShouldBeNil)
	})
}
