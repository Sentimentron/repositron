package models

type Query struct {
	Name     string      `json:"name"`
	Bucket   string      `json:"bucket"`
	Checksum string      `json:"checksum"`
	Metadata MetadataMap `json:"metadata"`
}
