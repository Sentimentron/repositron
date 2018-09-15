package main

import (
	"flag"
	"fmt"
	"github.com/Sentimentron/repositron/api"
	"github.com/Sentimentron/repositron/content"
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

	// Configure the URLs
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	uiDir, err := filepath.Abs("ui")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine the absolute path of %s", uiDir)
		os.Exit(1)
	}

	r.Handle("/", ui.IndexEndpointFactory(metadataStore, uiDir))
	r.Handle("/upload", ui.UploadEndpointFactory(uiDir))
	r.Handle("/delete/{id:[0-9]+}", ui.DeleteConfirmEndpointFactory(metadataStore, uiDir))
	r.Handle("/del/{id:[0-9]+}", ui.DeleteEndpointFactory(metadataStore, dir))
	r.Handle("/upload/process", ui.ProcessUploadEndpointFactory(metadataStore, dir))

	//s.HandleFunc("/blobs/", BlobsHandler)
	s.Handle("/blobs/byId/{id:[0-9]+}", api.GetBlobDescriptionByIdEndpointFactory(metadataStore)).Methods("GET")
	s.Handle("/blobs/byId/{id:[0-9]+}", api.DeleteBlobByIdEndpointFactory(metadataStore, dir)).Methods("DELETE")
	s.Handle("/blobs/byId/{id:[0-9]+}/content", api.GetBlobContentEndpointFactory(metadataStore)).Methods("GET")
	s.Handle("/blobs/byId/{id:[0-9]+}/content", api.UploadContentEndpointFactory(metadataStore, contentStore)).Methods("PUT").Name("ContentUpload")
	s.Handle("/blobs/search", api.SearchBlobEndpointFactory(metadataStore)).Methods("POST")
	s.Handle("/blobs", api.ListAllBlobsEndpointFactory(metadataStore)).Methods("GET")
	s.Handle("/blobs", api.UploadDescriptionEndpointFactory(metadataStore, s)).Methods("PUT")
	s.Handle("/info", api.DescribeEndpoint()).Methods("GET")

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
