package models

import (
	"database/sql"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// Not an entity
type IModel[M any] interface {
	Get() IModel[M]
	GetId() int
	SetId(int) IModel[M]
	GetTableName() string
	CanUserRead(int) bool
	CanUserWrite(int) bool
	Load(db *sqlx.DB, flags ...bool) (IModel[M], error) // Propagating (i.e. non-lazy) .Load()s has a problem of circular dependencies and loading (potentially massive) collections vs. objects, but we could get away with it in this application if we wanted to
	SetDeletedAt(*time.Time) IModel[M]
}
type Model struct {
	Id        int        `db:"id"` // @TODO: Make this *int
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (m Model) SetDeletedAt(t *time.Time) IModel[any] {
	m.DeletedAt = t
	return m
}
func (m Model) GetTableName() string         { return "" }
func (m Model) Get() IModel[any]             { return m }
func (m Model) GetId() int                   { return m.Id }
func (m Model) SetId(id int) IModel[any]     { m.Id = id; return m }
func (m Model) CanUserRead(userId int) bool  { return false }
func (m Model) CanUserWrite(userId int) bool { return false }
func (m Model) Load(db *sqlx.DB, flags ...bool) (IModel[any], error) {
	return m, nil
}

func Read[M IModel[any]](db *sqlx.DB, m M) (M, error) {
	tableName := m.GetTableName()
	id := m.GetId()

	if err := db.Get(&m, `
		SELECT * FROM `+tableName+`
		WHERE id = ?
		LIMIT 1
	`, id); err != nil {
		var null M
		return null, err
	}

	return m, nil
}

func Count[M IModel[any]](db *sqlx.DB, m M, id int) (int64, error) {
	var count int64
	if err := db.Get(&count, "SELECT COUNT(*) FROM "+m.GetTableName()); err != nil {
		return count, err
	}
	return count, nil
}

// Update is really Upsert, due to time constraints and project req's
// @TODO add parameter "upsert bool" default to false
func Update[M IModel[any]](db *sqlx.DB, m M) (sql.Result, error) {
	fields := getDBFields(m.Get())      // e.g. []string{"id", "snake_case"}
	csv := strings.Join(fields, ", ")   // e.g. "id, name, description"
	csvc := strings.Join(fields, ", :") // e.g. ":id, :name, :description"
	sql := "INSERT IGNORE INTO " + m.GetTableName() +
		" (" + csv + ") VALUES (:" + csvc + ")"
	return sqlx.NamedExec(db, sql, m)
}

// getDBFields reflects on a struct and returns the values of fields with `db` tags,
// or a map[string]interface{} and returns the keys.
func getDBFields(value interface{}) []string {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	fields := []string{}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i).Tag.Get("db")
			if field != "" {
				fields = append(fields, field)
			}
		}
		return fields
	}
	if v.Kind() == reflect.Map {
		for _, keyv := range v.MapKeys() {
			fields = append(fields, keyv.String())
		}
		return fields
	}
	return []string{}
}
