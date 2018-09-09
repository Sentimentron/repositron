package repositron

import (
	"flag"
	"fmt"
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
	var dir string
	var quota int
	flag.StringVar(&dir, "dir", "static/", "The directory to serve files from. Defaults to static/.")
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

	// Configure the URLs
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	//s.HandleFunc("/blobs/", BlobsHandler)
	//s.HandleFunc("/blobs/{key}", BlobHandler)

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
