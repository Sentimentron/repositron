package models

type BlobSearch struct {
	Name     *string `json:"name"`
	Checksum *string `json:"checksum"`
	Bucket   *string `json:"bucket"`
}
