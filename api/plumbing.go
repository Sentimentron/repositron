package api

import (
	"github.com/gorilla/mux"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/server/ui"
	"net/http"
)

// AttachAPIMethods lets you attach Repositron methods to an existing HTTP router.
func AttachAPIMethods(syncStore interfaces.SynchronizationStore,
	contentStore interfaces.ContentStore, metadataStore interfaces.MetadataStore, uiDir string, staticDir string,
	shouldAttachDebugInterface bool, r *mux.Router) {

	// Machine-readable APIs are versioned on a separate prefix
	s := r.PathPrefix("/v1").Subrouter()

	if shouldAttachDebugInterface {
		r.Handle("/", ui.IndexEndpointFactory(metadataStore, uiDir))
		r.Handle("/upload", ui.UploadEndpointFactory(uiDir))
		r.Handle("/delete/{id:[0-9]+}", ui.DeleteConfirmEndpointFactory(metadataStore, uiDir))
		r.Handle("/del/{id:[0-9]+}", ui.DeleteEndpointFactory(metadataStore, contentStore))
		r.Handle("/upload/process", ui.ProcessUploadEndpointFactory(metadataStore, contentStore))
	}

	s.Handle("/blobs/byId/{id:[0-9]+}", GetBlobDescriptionByIdEndpointFactory(metadataStore)).Methods("GET")
	s.Handle("/blobs/byId/{id:[0-9]+}", DeleteBlobByIdEndpointFactory(metadataStore, contentStore)).Methods("DELETE")
	s.Handle("/blobs/byId/{id:[0-9]+}/content", GetBlobContentEndpointFactory(metadataStore)).Methods("GET")
	s.Handle("/blobs/byId/{id:[0-9]+}/content", UploadContentEndpointFactory(metadataStore, contentStore)).Methods("PUT").Name("ContentUpload")
	s.Handle("/blobs/byId/{id:[0-9]+}/content/append", AppendContentEndpointFactory(metadataStore, contentStore, syncStore))
	s.Handle("/blobs/search", SearchBlobEndpointFactory(metadataStore)).Methods("POST")
	s.Handle("/blobs", ListAllBlobsEndpointFactory(metadataStore)).Methods("GET")
	s.Handle("/blobs", UploadDescriptionEndpointFactory(metadataStore, s)).Methods("PUT")
	s.Handle("/info", DescribeEndpoint()).Methods("GET")

	// Set up a URL which will serve static files
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))


}
