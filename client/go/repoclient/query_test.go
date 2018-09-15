package repoclient

import (
	"testing"
	"github.com/Sentimentron/repositron/models"
	"strings"
	. "github.com/smartystreets/goconvey/convey"

	"time"
)

func TestRepositronConnection_Query(t *testing.T) {
	Convey("Should be able to query stuff...", t, func() {

		c, err := Connect(globalTestURL)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		metadata := models.MetadataMap{}
		metadata["key"] = "value"

		fixedContent1 := "<html><body>hi</body></html>"
		content1 := strings.NewReader(fixedContent1)

		fixedContent2 := "<html><body>world</body></html>"
		content2 := strings.NewReader(fixedContent2)

		info1 := models.Blob{
			Id:       0,
			Bucket:   "__testing",
			Date:     time.Now(),
			Class:    "temp",
			Checksum: "",
			Uploader: "__tester",
			Metadata: metadata,
			Size:     int64(len(fixedContent1)),
			Name:     "__test_upload_file",
		}
		info2 := models.Blob{
			Id:       0,
			Bucket:   "__testing",
			Date:     time.Now(),
			Class:    "temp",
			Checksum: "",
			Uploader: "__tester",
			Metadata: metadata,
			Size:     int64(len(fixedContent2)),
			Name:     "__test_upload_file_2",
		}

		Convey("Should be able to upload silently...", func() {
			newInfo1, err := c.Upload(&info1, content1, false)
			So(newInfo1, ShouldNotBeNil)
			So(err, ShouldBeNil)

			newInfo2, err := c.Upload(&info2, content2, false)
			So(err, ShouldBeNil)
			So(newInfo2, ShouldNotBeNil)


			Convey("Should be able to query...", func() {
				newerInfo1, err := c.QueryById(newInfo1.Id)
				So(err, ShouldBeNil)
				newerInfo2, err := c.QueryById(newInfo2.Id)
				So(err, ShouldBeNil)

				Convey("Should be able to search by bucket...", func(){
					ids, err := c.QueryByBucket( "__testing")
					So(err, ShouldBeNil)
					So(newerInfo1.Id, ShouldBeIn, ids)
				})
				Convey("Should be able to search by name...", func(){
					ids, err := c.QueryByName( "__test_upload_file")
					So(err, ShouldBeNil)
					So(newerInfo1.Id, ShouldBeIn, ids)
					So(newerInfo2.Id, ShouldNotBeIn, ids)
				})
				Convey("Should be able to search by checksum...", func(){
					ids, err := c.QueryByChecksum( "95d70659530e385bfae5d6eefe689d95ac463cb0c58235f19eef71bdaa725126")
					So(err, ShouldBeNil)
					So(newerInfo1.Id, ShouldBeIn, ids)
					So(newerInfo2.Id, ShouldNotBeIn, ids)
				})
			})
		})
	})
}

