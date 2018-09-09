package main

import (
	"flag"
	"fmt"
	"github.com/Sentimentron/repositron/api"
	"github.com/Sentimentron/repositron/database"
	"github.com/Sentimentron/repositron/ui"
	"github.com/Sentimentron/repositron/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

	// Create the store
	dataStore, err := database.CreateStore(store)
	if err != nil {
		log.Fatal(err)
	}

	// Configure the URLs
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	uiDir, err := filepath.Abs("ui")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine the absolute path of %s", uiDir)
		os.Exit(1)
	}

	r.Handle("/", ui.IndexEndpointFactory(dataStore, uiDir))
	r.Handle("/upload", ui.UploadEndpointFactory(uiDir))
	r.Handle("/delete/{id}", ui.DeleteConfirmEndpointFactory(dataStore, uiDir))
	r.Handle("/del/{id}", ui.DeleteEndpointFactory(dataStore))
	r.Handle("/upload/process", ui.ProcessUploadEndpointFactory(dataStore, dir))

	//s.HandleFunc("/blobs/", BlobsHandler)
	s.Handle("/blobs/{id}/content", api.BlobContentEndpointFactory(dataStore))

	// Set up a handler which will serve permanent files
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  300 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
