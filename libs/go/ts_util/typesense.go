package ts_util

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/typesense/typesense-go/v4/typesense"
	tapi "github.com/typesense/typesense-go/v4/typesense/api"
)

type Type string

const (
	STRING        Type = "string"
	STRINGARRAY        = "string[]"
	INT32              = "int32"
	INT32ARRAY         = "int32[]"
	INT64              = "int64"
	INT64ARRAY         = "int64[]"
	FLOAT              = "float"
	FLOATARRAY         = "float[]"
	BOOL               = "bool"
	BOOLARRAY          = "bool[]"
	GEOPOINT           = "geopoint"
	GEOPOINTARRAY      = "geopoint[]"
	OBJECT             = "object" // object is comparable to a go struct
	OBJECTARRAY        = "object[]"
	STRINGPTR          = "string*" // special type that can be string or []string
)

func ToField(in any) []tapi.Field {
	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr || val.Kind() != reflect.Struct {
		panic("input should be a struct")
	}
	fields, err := lexField(val.Type())
	if err != nil {
		panic(err)
	}
	return fields
}

func ToFieldWithName(in any) (string, []tapi.Field, error) {
	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr || val.Kind() != reflect.Struct {
		return "", nil, errors.New("input should be a struct")
	}
	name := cutDotToLower(val.Type().String())
	fields, err := lexField(val.Type())
	return name, fields, err
}

func cutDotToLower(s string) string {
	before, after, found := strings.Cut(s, ".")
	if found {
		s = after
	} else {
		s = before
	}
	return strings.ToLower(s)
}

func lexField(typ reflect.Type) ([]tapi.Field, error) {
	var collectionFields []tapi.Field
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue // Skip unexported fields
		}
		fieldType := typeAllowed(field.Type)
		if fieldType == OBJECT {
			// Check if Object is Composition
			if field.Anonymous {
				composited, err := lexField(field.Type)
				if err != nil {
					return nil, err
				}
				collectionFields = append(collectionFields, composited...)
				continue
			}
		}
		tags := field.Tag.Get("json") // tags save the field_name and settings
		apiField, err := parseField(fieldType, tags)
		if err != nil {
			return nil, err
		}
		collectionFields = append(collectionFields, apiField)
	}
	return collectionFields, nil
}

func parseField(T Type, tag string) (tapi.Field, error) {
	params := strings.Split(tag, ",")
	var field tapi.Field
	var True bool = true

	if len(params) == 0 {
		return tapi.Field{}, errors.New("field name has to be provided for matching")
	}

	field.Name = params[0]
	field.Type = string(T)

	for _, key := range params[1:] {
		switch key {
		case "optional": // optional fields, can be null
			field.Optional = &True
		case "facet": // If a field is facet its also automatically indexed
			field.Facet = &True
			field.Index = &True
		case "index":
			field.Index = &True
		case "sort":
			field.Sort = &True
		case "infix":
			field.Infix = &True
		default:
			if ref, ok := strings.CutPrefix(key, "dim:"); ok {
				dim, err := strconv.Atoi(ref)
				if err != nil {
					return tapi.Field{}, err
				}
				field.NumDim = &dim
			}
			if ref, ok := strings.CutPrefix(key, "join:"); ok {
				field.Reference = &ref
				continue
			}

			if ref, ok := strings.CutPrefix(key, "reference:"); ok {
				field.Reference = &ref
				continue
			}

		}
	}

	return field, nil
}

func typeAllowed(t reflect.Type) Type {
	switch t.Kind() {
	case reflect.String:
		return STRING
	case reflect.Int32, reflect.Int:
		return INT32
	case reflect.Int64:
		return INT64
	case reflect.Float32, reflect.Float64:
		return FLOAT
	case reflect.Bool:
		return BOOL
	case reflect.Slice:
		elemType := typeAllowed(t.Elem())
		if elemType != "" {
			return elemType + "[]"
		}
	case reflect.Struct:
		return OBJECT
	case reflect.Pointer:
		return typeAllowed(t.Elem())
	default:
		panic("not handles type")
	}
	fmt.Println(t.Kind())

	return ""
}

// Check if nested field should be supported, this is the case if an Object or ObjectArray has Index = true
func nestedFields(fields []tapi.Field) *bool {
	for _, v := range fields {
		if v.Type == OBJECT || v.Type == OBJECTARRAY {
			return &[]bool{true}[0]
		}
	}
	return nil // equals false
}

// WaitForReady waits up to 10 seconds for Typesense to be ready.
// Panics if Typesense is not ready within the timeout.
func WaitForReady(ctx context.Context, client *typesense.Client) {
	const maxWaitTime = 10 * time.Second
	const retryInterval = 500 * time.Millisecond
	startTime := time.Now()

	for {
		// Try to retrieve collections to check if Typesense is ready
		_, err := client.Collections().Retrieve(ctx, nil)
		if err == nil {
			// Typesense is ready
			return
		}

		// Check if error is 503 (Not Ready or Lagging)
		if strings.Contains(err.Error(), "503") || strings.Contains(err.Error(), "Not Ready") {
			if time.Since(startTime) >= maxWaitTime {
				panic(fmt.Sprintf("typesense not ready after %v: %v", maxWaitTime, err))
			}
			time.Sleep(retryInterval)
			continue
		}

		// Other error, panic immediately
		panic(fmt.Sprintf("typesense health check failed: %v", err))
	}
}

// CreateCollectionIfNotExists checks if a collection exists, and creates it if not.
// Panics if the create operation fails.
func CreateCollectionIfNotExists(ctx context.Context, client *typesense.Client, schema *tapi.CollectionSchema) {
	// Check if collection exists
	_, err := client.Collection(schema.Name).Retrieve(ctx)
	if err == nil {
		// Collection exists, do nothing
		return
	}

	// Collection doesn't exist, create it
	_, err = client.Collections().Create(ctx, schema)
	if err != nil {
		panic(fmt.Sprintf("failed to create collection %s: %v", schema.Name, err))
	}
}
