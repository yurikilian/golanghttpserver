package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mitchellh/mapstructure"
	"github.com/yurikilian/bills/internal/logger"
	"reflect"
	"strings"
)

type PsqlStorage[T any] struct {
	db        *sql.DB
	tableName *string
}

func (s *PsqlStorage[T]) Find(id float64) (*T, error) {
	var t *T

	rows, err := s.db.Query(fmt.Sprintf("SELECT * FROM %v WHERE id = $1", *s.tableName), id)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Log.Warn(context.Background(), err.Error())
		}
	}(rows)

	m := FirstRowToMap(rows)

	err = mapstructure.WeakDecode(m, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *PsqlStorage[T]) Create(entity *T) (*T, error) {

	val := reflect.ValueOf(entity).Elem()

	n := val.NumField()

	var columnsSb strings.Builder
	var valuesSb strings.Builder

	values := make([]any, 0, n)

	for i := 0; i < val.NumField(); i++ {
		var columnName string
		if columnTag, ok := val.Type().Field(i).Tag.Lookup("column"); ok {
			columnName = columnTag
		} else {
			columnName = val.Type().Field(i).Name
		}

		value := val.Field(i).Interface()
		values = append(values, value)

		if i > 0 {
			columnsSb.WriteString(", ")
			valuesSb.WriteString(", ")
		}
		columnsSb.WriteString(columnName)
		valuesSb.WriteString(fmt.Sprint("$", i+1))

	}

	query := fmt.Sprintf("INSERT INTO %v(%v) VALUES(%v)", *s.tableName, columnsSb.String(), valuesSb.String())
	_, err := s.db.Exec(query, values...)

	if err != nil {
		return nil, fmt.Errorf("could not insert row on database: %w", err)
	}
	return entity, nil
}

func FirstRowToMap(rows *sql.Rows) *map[string]interface{} {

	cols, _ := rows.Columns()

	data := make(map[string]interface{})

	if rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		for i, colName := range cols {
			data[colName] = columns[i]
		}
	}

	return &data

}

func NewPsqlStorage[T any](db *sql.DB, tableName *string) *PsqlStorage[T] {

	if tableName == nil || len(*tableName) == 0 {
		panic("Invalid table name")
	}

	return &PsqlStorage[T]{
		db:        db,
		tableName: tableName,
	}
}

var _ Storage[interface{}] = (*PsqlStorage[interface{}])(nil)
