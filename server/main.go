package main

import (
	"os"
	"github.com/Sentimentron/repositron/api"
	"net/http"
	"time"
	"log"
	"flag"
	"path/filepath"
	"fmt"
	"github.com/Sentimentron/repositron/utils"
	"github.com/Sentimentron/repositron/content"
	"github.com/Sentimentron/repositron/database"
	"github.com/Sentimentron/repositron/synchronization"
	"github.com/gorilla/mux"
)

func main() {

	// Configure some information about this whole thing
	var dir, store string
	var quota int
	flag.StringVar(&dir, "dir", "static/", "The directory to serve files from. Defaults to static/.")
	flag.StringVar(&store, "store", "const/v1.sqlite", "The Sqlite3 file containing the store.")
	flag.IntVar(&quota, "quota", 1, "Maximum temporary file quota")
	flag.Parse()

	dir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine the absolute path of %s", dir)
		os.Exit(1)
	}

	// Double check that all the flags exist
	if !utils.IsDirectory(dir) {
		fmt.Fprintf(os.Stderr, "Unable to locate dir (tried: %s)", dir)
		os.Exit(1)
	}

	// Create the metadata store
	metadataStore, err := database.CreateStore(store)
	if err != nil {
		log.Fatal(err)
	}

	// Create the on-disk store
	contentStore, err := content.CreateStore(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Create the synchronization store, which stops stuff colliding on append
	syncStore, err := synchronization.CreateMemorySynchronizationStore()
	if err != nil {
		log.Fatal(err)
	}

	// Configure the URLs
	r := mux.NewRouter()

	uiDir, err := filepath.Abs("ui")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine the absolute path of %s", uiDir)
		os.Exit(1)
	}

	// Configure all the URLs on this server
	api.AttachAPIMethods(syncStore, contentStore, metadataStore, uiDir, dir,true, r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  300 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
