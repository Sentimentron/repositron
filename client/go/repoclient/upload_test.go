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

func TestRepositronConnection_AppendDidNotExistBefore(t *testing.T) {
	Convey("Should be able to append to a file which we have a description of, but no data", t, func(){

		c, err := Connect(globalTestURL)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		metadata := models.MetadataMap{}
		metadata["key"] = "value"

		originalContent := "APPEND TO ME\n"
		fixedContent := "THESE THINGS YOU SEE\n"
		content := strings.NewReader(originalContent)
		info := models.Blob{
			Id:       0,
			Bucket:   "__testing",
			Date:     time.Now(),
			Class:    "temp",
			Checksum: "",
			Uploader: "__tester",
			Metadata: metadata,
			Size:     int64(len(originalContent)),
			Name:     "__test_upload_file",
		}

		Convey("Should be able to append silently....", func(){
			newInfo, err := c.Upload(&info, content, false)
			So(err, ShouldBeNil)
			So(newInfo, ShouldNotBeNil)

			Convey("Should be able to append to this document...", func(){
				newerInfo, err := c.Append(newInfo, int64(len(fixedContent)), strings.NewReader(fixedContent), false)
				So(err, ShouldBeNil)
				So(newerInfo.Size, ShouldEqual, len(fixedContent)+len(originalContent))
				So(newerInfo.Checksum, ShouldEqual, "da711a5fd60b670b8afbdc39b1eec97a9755addfbe6dc736804ab3d8c552c001")
			})

		})

	})
}