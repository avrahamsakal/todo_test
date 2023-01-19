package models

import (
	"github.com/jmoiron/sqlx"
)

func (meta Metadata) GetTableName() string {
	return "metadata"
}

type Metadata struct {
	Model
	Key   string `db:"text"`
	Value string `db:"text"`
}

func (meta Metadata) CanMetadataRead(MetadataId int) bool  { return false }
func (meta Metadata) CanMetadataWrite(MetadataId int) bool { return false }

// GetMetadataByKey retrieves the value associated with the given key
func GetMetadataByKey(db *sqlx.DB, key string) (string, error) {
	var value string
	if err := db.Get(&value, `
		SELECT value FROM metadata WHERE key = "?"
	`, key); err != nil {
		return "", err
	}

	return value, nil
}
