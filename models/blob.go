package models

import "time"

// BlobType consists of a
type BlobType string

const (
	PermanentBlob BlobType = "permanent"
	TemporaryBlob BlobType = "temp"
)

type Blob struct {
	Id       int64       `db:"id" json:"id"`
	Name     string      `json:"name" validate:"required"`
	Bucket   string      `json:"bucket" validate:"required"`
	Date     time.Time   `json:"uploaded"`
	Class    BlobType    `json:"type" validate:"required,oneof=permanent temp"`
	Checksum string      `db:"sha1" json:"sha1"`
	Uploader string      `json:"owner" validate:"required"`
	Metadata MetadataMap `json:"metadata"`
	Size     int64       `json:"size"`
}
