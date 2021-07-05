package simplejson

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

// Value returns json value, implement driver.Valuer interface.
func (j JSON) Value() (driver.Value, error) {
	if j.data == nil {
		return nil, nil
	}

	return j.MarshalJSON()
}

// Scan scans value into Json, implements sql.Scanner interface.
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return j.UnmarshalJSON(bytes)
}

// GormDataType implements gorm.schema.GormDataTypeInterface.
func (JSON) GormDataType() string {
	return "json"
}

// GormDBDataType implements gorm.migrator.GormDataTypeInterface.
func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	}
	return ""
}

// JSONQueryExpression implements clause.Expression interface to use as querier.
type JSONQueryExpression struct {
	column      string
	keys        []string
	hasKeys     bool
	equals      bool
	equalsValue interface{}
}

// JSONQuery queries column as json.
func JSONQuery(column string) *JSONQueryExpression {
	return &JSONQueryExpression{column: column}
}

// HasKey returns clause.Expression.
func (jsonQuery *JSONQueryExpression) HasKey(keys ...string) *JSONQueryExpression {
	jsonQuery.keys = keys
	jsonQuery.hasKeys = true
	return jsonQuery
}

// Keys returns clause.Expression.
func (jsonQuery *JSONQueryExpression) Equals(value interface{}, keys ...string) *JSONQueryExpression {
	jsonQuery.keys = keys
	jsonQuery.equals = true
	jsonQuery.equalsValue = value
	return jsonQuery
}

// Build implements clause.Expression.
func (jsonQuery *JSONQueryExpression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		switch stmt.Dialector.Name() {
		case "mysql", "sqlite":
			switch {
			case jsonQuery.hasKeys:
				if len(jsonQuery.keys) > 0 {
					builder.WriteString(fmt.Sprintf("JSON_EXTRACT(%s, '$.%s') IS NOT NULL",
						stmt.Quote(jsonQuery.column), strings.Join(jsonQuery.keys, ".")))
				}
			case jsonQuery.equals:
				if len(jsonQuery.keys) > 0 {
					builder.WriteString(fmt.Sprintf("JSON_EXTRACT(%s, '$.%s') = ",
						stmt.Quote(jsonQuery.column), strings.Join(jsonQuery.keys, ".")))
					stmt.AddVar(builder, jsonQuery.equalsValue)
				}
			}
		}
	}
}
