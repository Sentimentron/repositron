package api

import (
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/gorilla/mux"
	"net/http"
)

func BlobContentEndpointFactory(store interfaces.MetadataStore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]

		redirectString := fmt.Sprintf("/static/%s", id)
		http.Redirect(w, r, redirectString, 307)

	})

}
