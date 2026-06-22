package cdc

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pglogrepl"
)

func MapTuple[T any](rel *pglogrepl.RelationMessage, tuple *pglogrepl.TupleData) T {
	var result T
	if tuple == nil {
		return result
	}

	rv := reflect.ValueOf(&result).Elem()
	rt := rv.Type()

	colIndex := buildColumnIndex(rel)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := getColumnName(field)

		idx, ok := colIndex[tag]
		if !ok || idx >= len(tuple.Columns) {
			continue
		}

		col := tuple.Columns[idx]
		if col.DataType == 'n' {
			continue
		}

		data := string(col.Data)
		setFieldValue(rv.Field(i), data)
	}

	return result
}

func getColumnName(field reflect.StructField) string {
	if tag := field.Tag.Get("col"); tag != "" {
		return tag
	}
	if tag := field.Tag.Get("json"); tag != "" {
		if idx := strings.Index(tag, ","); idx != -1 {
			tag = tag[:idx]
		}
		if tag != "-" {
			return tag
		}
	}
	if tag := field.Tag.Get("db"); tag != "" {
		return tag
	}
	return field.Name
}

func buildColumnIndex(rel *pglogrepl.RelationMessage) map[string]int {
	idx := make(map[string]int, len(rel.Columns))
	for i, col := range rel.Columns {
		idx[col.Name] = i
	}
	return idx
}

func setFieldValue(field reflect.Value, data string) {
	if !field.CanSet() {
		return
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(data)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := strconv.ParseInt(data, 10, 64); err == nil {
			field.SetInt(v)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := strconv.ParseUint(data, 10, 64); err == nil {
			field.SetUint(v)
		}

	case reflect.Float32, reflect.Float64:
		if v, err := strconv.ParseFloat(data, 64); err == nil {
			field.SetFloat(v)
		}

	case reflect.Bool:
		field.SetBool(data == "t" || data == "true" || data == "1")

	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			for _, layout := range []string{
				time.RFC3339Nano,
				time.RFC3339,
				"2006-01-02 15:04:05.999999",
				"2006-01-02 15:04:05",
				"2006-01-02",
			} {
				if t, err := time.Parse(layout, data); err == nil {
					field.Set(reflect.ValueOf(t))
					break
				}
			}
		}

	case reflect.Array:
		if field.Type() == reflect.TypeOf(uuid.UUID{}) {
			if u, err := uuid.Parse(data); err == nil {
				field.Set(reflect.ValueOf(u))
			}
		}
	}
}
