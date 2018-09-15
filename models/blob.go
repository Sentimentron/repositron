package models

import (
	"gopkg.in/go-playground/validator.v9"
	"time"
)

// BlobType consists of a
type BlobType string

const (
	PermanentBlob BlobType = "permanent"
	TemporaryBlob BlobType = "temp"
)

type Blob struct {
	Id       int64       `db:"id" json:"id"`
	Name     string      `json:"name" validate:"required" db:"name"`
	Bucket   string      `json:"bucket" validate:"required" db:"bucket"`
	Date     time.Time   `json:"uploaded" db:"date"`
	Class    BlobType    `json:"type" validate:"required,oneof=permanent temp" db:"class"`
	Checksum string      `db:"sha1" json:"sha1"`
	Uploader string      `json:"owner" validate:"required" db:"uploader"`
	Metadata MetadataMap `json:"metadata" db:"metadata"`
	Size     int64       `json:"size" db:"size"`
}

func (b *Blob) Validate() error {

	validate := validator.New()

	return validate.Struct(b)

}
