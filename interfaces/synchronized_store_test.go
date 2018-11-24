package interfaces

import ("testing"
. "github.com/smartystreets/goconvey/convey"
	"../content"
	"../database"
	"os"
	"io/ioutil"
	"log"
)

func createStoresForTest() (ContentStore, MetadataStore) {

	// Create a content store for testing in a temporary directory
	tmpDirPrefix := os.TempDir()
	tmpDir, err := ioutil.TempDir(tmpDirPrefix, "repoTest-")
	if err != nil {
		panic(err)
	}
	contentStore, err := content.CreateStore(tmpDir)
	if err != nil {
		panic(err)
	}

	// Create a metadata store, also in a temporary directory, for testing
	tmpFile, err := ioutil.TempFile("", "repo")
	log.Printf("Creating temporary file at: %s", tmpFile.Name())
	os.Remove(tmpFile.Name())
	handle, err := database.CreateStore(tmpFile.Name())
	if err != nil {
		panic(err)
	}
	return contentStore, handle
}

func TestCreateCombinedStore(t *testing.T) {

	Convey("Should be able to create the combined store...", t, func(){
		contentStore, metadataStore := createStoresForTest()
		var combinedStore BlobStore

	})


}
