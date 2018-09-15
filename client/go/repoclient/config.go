package repoclient

import (
	"encoding/json"
	"github.com/Sentimentron/repositron/utils"
	"log"
	"os"
	"os/user"
	"path"
)

type ClientConfiguration struct {
	BaseURL       string `json:"baseURL"`
	ConfigVersion string `json:"configVersion"`
}

func BuildDefaultClientConfigurationPath() string {

	// Retrieve the user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Print("client error: Could not determine the current user")
		log.Fatal(err)
	}

	// Create the .config directory, if it doesn't exist
	configDir := path.Join(usr.HomeDir, ".config")
	if !utils.IsDirectory(configDir) {
		err = os.Mkdir(configDir, 0700)
		if err != nil {
			log.Print("client error: Could not create .config for the current user")
			log.Fatal(err)
		}
	}

	// Return the final path
	return path.Join(configDir, "repositron.conf")
}

func ReadClientConfiguration(path string) *ClientConfiguration {

	// Open the configuration file
	f, err := os.Open(path)
	if err != nil {
		log.Printf("client error: Could not open configuration at: %s", path)
		log.Fatal(err)
	}

	// Decode the configuration file
	dec := json.NewDecoder(f)
	var c ClientConfiguration
	err = dec.Decode(&c)
	if err != nil {
		log.Printf("client error: Could not decode configuration at: %s", path)
		log.Fatal(err)
	}

	if c.ConfigVersion != "1" {
		log.Fatalf("client error: Could not use configuration at: %s (unsupported version)", path)
	}

	return &c
}

func WriteClientConfiguration(c *ClientConfiguration, path string) {

	f, err := os.Create(path)
	if err != nil {
		log.Printf("client error: could not create configuration at %s", path)
		log.Fatal(err)
	}

	enc := json.NewEncoder(f)
	err = enc.Encode(c)
	if err != nil {
		log.Printf("client error: could not create configuration at %s", path)
		log.Fatal(err)
	}

	f.Close()
}
