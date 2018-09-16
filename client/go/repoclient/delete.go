package repoclient

import (
	"fmt"
	"net/http"
)

func (c *RepositronConnection) Delete(blobId int64) error {

	client := &http.Client{}

	contentUrl := c.GetURL(fmt.Sprintf("v1/blobs/byId/%d", blobId))

	req, err := http.NewRequest("DELETE", contentUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Bad status code: expected %d, got %d", 202, resp.StatusCode)
	}

	return nil

}
