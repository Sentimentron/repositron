package repoclient

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"github.com/Sentimentron/repositron/models"
	"time"
)

func TestRepositronConnection_UploadVerbose(t *testing.T) {
	Convey("Should be able to upload...", t, func() {

		c, err := Connect(globalTestURL)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		metadata := models.MetadataMap{}
		metadata["key"] = "value"

		fixedContent := "<html><body>hi</body></html>"
		content := strings.NewReader(fixedContent)
		info := models.Blob{
			Id: 0,
			Bucket: "__testing",
			Date: time.Now(),
			Class: "temp",
			Checksum: "",
			Uploader: "__tester",
			Metadata: metadata,
			Size: int64(len(fixedContent)),
			Name: "__test_upload_file",
		}

		Convey("Should be able to upload silently...", func(){
			err := c.Upload(&info, content, false)
			So(err, ShouldBeNil)
		})
		Convey("Should be able to upload loudly...", func(){
			err := c.Upload(&info, content, true)
			So(err, ShouldBeNil)
		})

	})
}