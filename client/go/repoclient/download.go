package repoclient

import (
	"crypto/sha256"
	"fmt"
	"github.com/Sentimentron/repositron/models"
	"io"
	"log"
	"net/http"
)

type DownloadWriter struct {
	underlyingWriter       io.Writer
	lastReportedPercentage string
	downloadedSoFar        int64
	totalDownloadSize      int64
}

func NewDownloadWriter(w io.Writer, size int64) *DownloadWriter {
	return &DownloadWriter{w, "", 0, size}
}

func (d *DownloadWriter) Write(p []byte) (int, error) {
	// Write to the underlying writer
	n, err := d.underlyingWriter.Write(p)
	if err != nil {
		log.Printf("Error writing: %v", err)
		return n, err
	}

	d.downloadedSoFar += int64(n)
	defer d.reportProgress()

	return n, err
}

func (d *DownloadWriter) reportProgress() {
	percentage := (100.0 * float32(d.downloadedSoFar) / float32(d.totalDownloadSize))
	pc := fmt.Sprintf("%.2f%% downloaded...\r", percentage)
	if pc == d.lastReportedPercentage {
		return
	}
	d.lastReportedPercentage = pc
	fmt.Print(pc)
}



func (c *RepositronConnection) Download(b *models.Blob, w io.Writer, verbose bool) error {

	if verbose {
		w = NewDownloadWriter(w, b.Size)
	}

	// Request the object
	contentUrl := c.GetURL(fmt.Sprintf("/v1/blobs/byId/%d/content", b.Id))
	if verbose {
		log.Printf("Downloading from: %s", contentUrl)
	}
	resp, err := http.Get(contentUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a teewriter so we can verify the checksum
	h := sha256.New()
	tee := io.TeeReader(resp.Body, h)
	n, err := io.Copy(w, tee)
	checksum := fmt.Sprintf("%x", h.Sum(nil))
	if n != b.Size {
		return fmt.Errorf("did not read enough (expected %d byte(s), got %d)", b.Size, n)
	} else if checksum != b.Checksum {
		return fmt.Errorf("checksum mismatch (expected %s, got %s)", b.Checksum, checksum)
	}

	return nil
}
