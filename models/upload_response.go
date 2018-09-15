package models

type BlobUploadResponse struct {
	RedirectURL string `json:"redirectURL"`
	Blob        *Blob
}
