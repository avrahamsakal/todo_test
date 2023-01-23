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

// Restrict to metadata R/W access to user 0 (admin user);
// essentially this gives us a secure admin panel
// (CanUserReads employ kal vaChomers with CanUserWrite)
func (meta Metadata) CanUserRead(userId int) bool  { return meta.CanUserWrite(userId) }
func (meta Metadata) CanUserWrite(userId int) bool { return userId == 0 }

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
