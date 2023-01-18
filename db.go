package main

import (
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func GetDB(driverName, dataSourceName string) *sqlx.DB {
	if db != nil {
		return db
	}

	db, _ = sqlx.Open(driverName, dataSourceName)

	return db
}
