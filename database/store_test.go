package database

import (
	"github.com/Sentimentron/repositron/models"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestCreateStore(t *testing.T) {
	Convey("Given an arbitrary file...", t, func() {
		tmpFile, err := ioutil.TempFile("", "repo")
		So(err, ShouldBeNil)
		log.Printf("Creating temporary file at: %s", tmpFile.Name())
		os.Remove(tmpFile.Name())

		Convey("Should be able to create a database there...", func() {

			handle, err := CreateStore(tmpFile.Name())
			So(err, ShouldBeNil)
			So(handle, ShouldNotBeNil)

			Convey("Should be able to close the store...", func() {
				err = handle.Close()
				So(err, ShouldBeNil)

				Convey("Should be able to re-open the database there too (though this is not allowed)", func() {
					handle, err := CreateStore(tmpFile.Name())
					So(err, ShouldBeNil)
					So(handle, ShouldNotBeNil)
				})
			})

		})
	})
}

func TestStore_RetrieveBlobById(t *testing.T) {
	Convey("Given a blank store...", t, func() {
		tmpFile, err := ioutil.TempFile("", "repo")
		So(err, ShouldBeNil)
		log.Printf("Creating temporary file at: %s", tmpFile.Name())
		os.Remove(tmpFile.Name())

		handle, err := CreateStore(tmpFile.Name())
		So(err, ShouldBeNil)

		Convey("Should return the right error if no matching records...", func() {
			b, err := handle.RetrieveBlobById(1231)
			So(err, ShouldEqual, NoMatchingBlobsError)
			So(b, ShouldBeNil)
		})
	})
}

func TestStore_StoreBlobRecord(t *testing.T) {
	Convey("Given a blank store...", t, func() {

		tmpFile, err := ioutil.TempFile("", "repo")
		So(err, ShouldBeNil)
		log.Printf("Creating temporary file at: %s", tmpFile.Name())
		os.Remove(tmpFile.Name())

		handle, err := CreateStore(tmpFile.Name())
		So(err, ShouldBeNil)
		So(handle, ShouldNotBeNil)

		Convey("Should be able to insert a WIP-blob record", func() {

			metadata := make(map[string]interface{})
			metadata["some"] = "val"

			b := &models.Blob{
				0,
				"my_test_file",
				"test_bucket",
				time.Now(),
				models.TemporaryBlob,
				"",
				"default",
				metadata,
				-1,
			}

			inserted, err := handle.StoreBlobRecord(b)
			So(err, ShouldBeNil)
			So(inserted, ShouldNotBeNil)
			So(inserted.Id, ShouldBeGreaterThan, 0)

			Convey("Should then be able to finalize it:", func() {
				c := *inserted
				c.Checksum = "asdfasdfasdfasdf"
				c.Size = 40

				updated, err := handle.FinalizeBlobRecord(&c)
				So(err, ShouldBeNil)
				So(updated, ShouldNotBeNil)

				So(updated.Checksum, ShouldEqual, "asdfasdfasdfasdf")
				So(updated.Size, ShouldEqual, 40)

				Convey("Should then be able to get it via checksum...", func() {
					cur, err := handle.GetBlobIdsMatchingChecksum("asdfasdfasdfasdf")
					So(err, ShouldBeNil)
					So(cur, ShouldNotBeNil)
					So(len(cur), ShouldEqual, 1)
					So(cur, ShouldResemble, []int64{c.Id})
				})

				Convey("Should then be able to get it via name...", func() {
					cur, err := handle.GetBlobIdsMatchingName("my_test_file")
					So(err, ShouldBeNil)
					So(cur, ShouldNotBeNil)
					So(len(cur), ShouldEqual, 1)
					So(cur, ShouldResemble, []int64{c.Id})
				})

				Convey("Should then be able to get it via bucket...", func() {
					cur, err := handle.GetBlobIdsMatchingBucket("test_bucket")
					So(err, ShouldBeNil)
					So(cur, ShouldNotBeNil)
					So(len(cur), ShouldEqual, 1)
					So(cur, ShouldResemble, []int64{c.Id})
				})

				Convey("Should be able to compare them...", func() {
					cur, err := handle.RetrieveBlobById(c.Id)
					So(err, ShouldBeNil)
					So(cur, ShouldResemble, c)
				})

			})

		})
	})
}

func TestStore_DeleteById(t *testing.T) {
	Convey("Given a blank store...", t, func() {

		tmpFile, err := ioutil.TempFile("", "repo")
		So(err, ShouldBeNil)
		log.Printf("Creating temporary file at: %s", tmpFile.Name())
		os.Remove(tmpFile.Name())

		handle, err := CreateStore(tmpFile.Name())
		So(err, ShouldBeNil)

		Convey("And an inserted item...", func() {
			metadata := make(map[string]interface{})
			metadata["some"] = "val"

			b := &models.Blob{
				0,
				"my_test_file",
				"test_bucket",
				time.Now(),
				models.TemporaryBlob,
				"",
				"default",
				metadata,
				-1,
			}

			inserted, err := handle.StoreBlobRecord(b)
			So(err, ShouldBeNil)
			So(inserted, ShouldNotBeNil)

			c := *inserted
			c.Checksum = "asdfasdfasdfasdf"
			c.Size = 40

			updated, err := handle.FinalizeBlobRecord(&c)
			So(err, ShouldBeNil)

			Convey("Should be able to delete that item...", func() {
				err := handle.DeleteById(updated.Id)
				So(err, ShouldBeNil)

				_, err = handle.RetrieveBlobById(updated.Id)
				So(err, ShouldEqual, NoMatchingBlobsError)
			})
		})
	})
}
