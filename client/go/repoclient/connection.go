package repoclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sentimentron/repositron/models"
	"net/http"
)

const SupportedAPIVersion = "1.0"

var UnsupportedAPIVersionError = errors.New("unsupported api verson")

type RepositronConnection struct {
	BaseURL string
}

func (c *RepositronConnection) GetURL(sub string) string {
	return fmt.Sprintf("%s/%s", c.BaseURL, sub)
}

func Connect(baseURL string) (*RepositronConnection, error) {

	var desc models.APIDescription

	ret := &RepositronConnection{baseURL}

	// Request the information
	resp, err := http.Get(ret.GetURL("v1/info"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the description
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&desc)
	if err != nil {
		return nil, err
	}

	// Validate that we can connect
	if desc.APIVersion != SupportedAPIVersion {
		return nil, UnsupportedAPIVersionError
	}

	return ret, nil
}
