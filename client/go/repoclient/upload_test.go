package repoclient

import (
	"github.com/Sentimentron/repositron/models"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func TestRepositronConnection_Upload(t *testing.T) {
	Convey("Should be able to upload...", t, func() {

		c, err := Connect(globalTestURL)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		metadata := models.MetadataMap{}
		metadata["key"] = "value"

		fixedContent := "<html><body>hi</body></html>"
		content := strings.NewReader(fixedContent)
		info := models.Blob{
			Id:       0,
			Bucket:   "__testing",
			Date:     time.Now(),
			Class:    "temp",
			Checksum: "",
			Uploader: "__tester",
			Metadata: metadata,
			Size:     int64(len(fixedContent)),
			Name:     "__test_upload_file",
		}

		Convey("Should be able to upload silently...", func() {
			newInfo, err := c.Upload(&info, content, false)
			So(newInfo, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(newInfo.Id, ShouldBeGreaterThan, 0)
		})
		Convey("Should be able to upload loudly...", func() {
			newInfo, err := c.Upload(&info, content, true)
			So(newInfo, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(newInfo.Id, ShouldBeGreaterThan, 0)
		})

	})
}
