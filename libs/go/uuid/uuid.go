package uuid

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

type UUID = uuid.UUID

func New() UUID {
	return uuid.New()
}

func Parse(s string) (UUID, error) {
	return uuid.Parse(s)
}

func ToNull(u uuid.UUID) NullableUUID {
	return NullableUUID{
		UUID:  u,
		Valid: u != uuid.Nil,
	}
}

// NullableUUID represents a UUID that can be null
type NullableUUID struct {
	UUID  UUID
	Valid bool // Valid is true if UUID is not NULL
}

// Value implements the driver.Valuer interface
func (n NullableUUID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.UUID[:], nil
}

// Scan implements the sql.Scanner interface
func (n *NullableUUID) Scan(value interface{}) error {
	if value == nil {
		n.UUID = UUID{}
		n.Valid = false
		return nil
	}

	if bytes, ok := value.([]byte); ok {
		n.UUID = UUID(bytes)
		n.Valid = true
		return nil
	}

	if str, ok := value.(string); ok {
		var err error
		n.UUID, err = uuid.Parse(str)
		if err != nil {
			return err
		}
		n.Valid = true
		return nil
	}

	return nil
}

// MarshalJSON implements json.Marshaler interface
func (n NullableUUID) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.UUID.String())
}

// UnmarshalJSON implements json.Unmarshaler interface
func (n *NullableUUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.UUID = UUID{}
		n.Valid = false
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, err := uuid.Parse(str)
	if err != nil {
		return err
	}

	n.UUID = parsed
	n.Valid = true
	return nil
}

// NewNullableUUID creates a new NullableUUID from a UUID
func NewNullableUUID(u UUID) NullableUUID {
	return NullableUUID{
		UUID:  u,
		Valid: true,
	}
}

// NullNullableUUID creates a null NullableUUID
func NullNullableUUID() NullableUUID {
	return NullableUUID{
		UUID:  UUID{},
		Valid: false,
	}
}
