package models

import "time"

// BlobType consists of a
type BlobType string

const (
	PermanentBlob BlobType = "permanent"
	TemporaryBlob BlobType = "temp"
)

type Blob struct {
	Id       int64       `db:"id",json:"id"`
	Name     string      `json:"name"`
	Bucket   string      `json:"bucket"`
	Date     time.Time   `json:"uploaded"`
	Class    BlobType    `json:"type"`
	Checksum string      `db:"sha1",json:"sha1"`
	Uploader string      `json:"owner"`
	Metadata MetadataMap `json:"metadata"`
	Size     int64       `json:"size"`
}
