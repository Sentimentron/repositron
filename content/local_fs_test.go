package content

import (
	"bytes"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreateStore2(t *testing.T) {
	Convey("Should be able to Create a local_fs store...", t, func() {

		var store interfaces.ContentStore

		// Create a temporary directory containing the store
		tmpDirPrefix := os.TempDir()
		tmpDir, err := ioutil.TempDir(tmpDirPrefix, "repoTest-")
		So(err, ShouldBeNil)

		// Create the store and check for issues
		store, err = CreateStore(tmpDir)
		So(err, ShouldBeNil)
		So(store, ShouldNotBeNil)

		Convey("Should not be able to open a store in a directory that does not exist...", func() {
			_, err := CreateStore("/does/not/exist")
			So(err, ShouldNotBeNil)
		})

	})
}

func getStoreForTesting() *FileSystemContentStore {
	// Create a temporary directory containing the store
	tmpDirPrefix := os.TempDir()
	tmpDir, err := ioutil.TempDir(tmpDirPrefix, "repoTest-")
	if err != nil {
		panic(err)
	}
	So(err, ShouldBeNil)

	// Create the store and check for issues
	store, err := CreateStore(tmpDir)
	if err != nil {
		panic(err)
	}

	return store
}

func TestFileSystemContentStore_WriteBlobContent(t *testing.T) {
	Convey("Should be able to write content into the store...", t, func() {
		store := getStoreForTesting()
		blob := &models.Blob{
			Id:       1,
			Name:     "test_file",
			Bucket:   "test_bucket",
			Date:     time.Now(),
			Class:    models.BlobType("temp"),
			Checksum: "",
			Uploader: "",
			Metadata: nil,
			Size:     0,
		}

		content := "some content"
		reader := strings.NewReader(content)

		read, err := store.WriteBlobContent(blob, reader)
		So(err, ShouldBeNil)
		So(read.Size, ShouldEqual, len(content))

		Convey("Should be able to read what was written...", func() {
			// Should be able to retrieve that content later
			writer := &bytes.Buffer{}
			written, err := store.RetrieveBlobContent(blob, writer)

			So(err, ShouldBeNil)
			So(written, ShouldEqual, len(content))
			So(string(writer.Bytes()), ShouldEqual, "some content")
		})
	})
}

func TestFileSystemContentStore_InsertBlobContent(t *testing.T) {
	Convey("Should be able to insert content into blobs...", t, func() {
		store := getStoreForTesting()
		blob := &models.Blob{
			Id:       1,
			Name:     "test_file",
			Bucket:   "test_bucket",
			Date:     time.Now(),
			Class:    models.BlobType("temp"),
			Checksum: "",
			Uploader: "",
			Metadata: nil,
			Size:     0,
		}

		content := "some content"
		reader := strings.NewReader(content)

		written, err := store.InsertBlobContent(blob, 64, reader)
		So(err, ShouldBeNil)
		So(written.Size, ShouldEqual, len(content)+64)

		Convey("Should be able to retrieve that content...", func() {
			buf := new(bytes.Buffer)
			read, err := store.RetrieveBlobContent(blob, buf)
			So(read, ShouldEqual, 64+len(content))
			So(err, ShouldBeNil)

			emptyBuf := make([]byte, 64)

			So(buf.Bytes()[:64], ShouldResemble, emptyBuf)
			So(buf.Bytes()[64:], ShouldResemble, []byte(content))

		})

	})
}

func TestFileSystemContentStore_DeleteBlobContent(t *testing.T) {
	Convey("Should be able to delete blobs when they're no longer used...", t, func() {
		store := getStoreForTesting()
		blob := &models.Blob{
			Id:       1,
			Name:     "test_file",
			Bucket:   "test_bucket",
			Date:     time.Now(),
			Class:    models.BlobType("temp"),
			Checksum: "",
			Uploader: "",
			Metadata: nil,
			Size:     0,
		}

		content := "some content"
		reader := strings.NewReader(content)

		_, err := store.WriteBlobContent(blob, reader)
		So(err, ShouldBeNil)

		Convey("Should be able to delete this blob...", func() {
			err := store.DeleteBlobContent(blob)
			So(err, ShouldBeNil)

			Convey("Blob content should no longer be accessible...", func() {
				writer := &bytes.Buffer{}
				_, err := store.RetrieveBlobContent(blob, writer)
				So(err, ShouldNotBeNil)
			})
		})

	})
}

func TestFileSystemContentStore_AppendBlobContent(t *testing.T) {
	Convey("Should be able to append content to blobs...", t, func() {
		store := getStoreForTesting()
		blob := &models.Blob{
			Id:       1,
			Name:     "test_file",
			Bucket:   "test_bucket",
			Date:     time.Now(),
			Class:    models.BlobType("temp"),
			Checksum: "",
			Uploader: "",
			Metadata: nil,
			Size:     0,
		}

		content := "some content"
		reader := strings.NewReader(content)

		read, err := store.WriteBlobContent(blob, reader)
		So(err, ShouldBeNil)
		So(read.Size, ShouldEqual, len(content))

		Convey("Should be able to append to that blob...", func() {
			additionalContent := " shall be appended"
			reader := strings.NewReader(additionalContent)

			written, err := store.AppendBlobContent(read, reader)
			So(written.Size, ShouldEqual, len(additionalContent)+len(content))
			So(err, ShouldBeNil)

			Convey("Should be able to read that content back...", func() {
				writer := &bytes.Buffer{}
				written, err := store.RetrieveBlobContent(blob, writer)

				So(err, ShouldBeNil)
				So(string(writer.Bytes()), ShouldEqual, "some content shall be appended")
				So(written, ShouldEqual, len(content)+len(additionalContent))
			})
		})
	})
}
