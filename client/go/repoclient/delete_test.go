package repoclient

import (
	"bytes"
	"github.com/Sentimentron/repositron/models"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func TestRepositronConnection_Delete(t *testing.T) {
	Convey("Should be able to delete...", t, func() {

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

			Convey("Should be able to query...", func() {
				newerInfo, err := c.QueryById(newInfo.Id)
				So(err, ShouldBeNil)
				Convey("Should be able to delete noisily...", func() {
					err := c.Delete(newerInfo.Id)
					So(err, ShouldBeNil)
					Convey("Content should not be accessible...", func(){
						var buf bytes.Buffer
						err := c.Download(newerInfo, &buf, true)
						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
