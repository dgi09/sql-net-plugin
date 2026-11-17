package internal

import (
	"strings"
	"unicode"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

type Field struct {
	Idx           int
	NetType       string
	NetRecordName string
	NetArgName    string
	SqlName       string
	IsArray       bool
	IsNullable    bool
}

func MapSqlType(t string) string {
	t = strings.ReplaceAll(t, "pg_catalog.", "")

	switch t {
	case "boolean":
		return "bool"

	case "integer":
		return "int"

	case "int":
		return "int"

	case "int4":
		return "int"

	case "int8":
		return "long"

	case "smallint":
		return "short"

	case "bigint":
		return "long"

	case "bigserial":
		return "long"

	case "real":
		return "float"

	case "double precision":
		return "double"

	case "numeric":
		return "decimal"

	case "money":
		return "decimal"

	case "uuid":
		return "Guid"

	case "timestamp":
		return "DateTime"

	case "timestamptz":
		return "DateTime"

	case "date":
		return "DateTime"

	case "text":
		return "string"

	case "character varying":
		return "string"

	case "varchar":
		return "string"

	case "character":
		return "string"

	case "citext":
		return "string"

	case "json":
		return "string"

	case "jsonb":
		return "string"

	case "bytea":
		return "byte"

	case "interval":
		return "TimeSpan"
	}

	return t
}

func ToNetRecordName(cname string) string {
	subWords := strings.Split(cname, "_")

	result := ""

	for _, sw := range subWords {
		word := []rune(sw)
		word[0] = unicode.ToUpper(word[0])

		result += string(word)
	}

	return result
}

func ToNetArgName(cname string) string {
	subWords := strings.Split(cname, "_")

	result := ""

	for i, sw := range subWords {
		word := []rune(sw)
		if i == 0 {
			word[0] = unicode.ToLower(word[0])
		} else {
			word[0] = unicode.ToUpper(word[0])
		}

		result += string(word)
	}

	return result
}

func MapColumn(idx int, col *plugin.Column) Field {
	columnName := col.Name
	sqlType := col.Type.GetName()

	return Field{
		Idx:           idx,
		NetType:       MapSqlType(sqlType),
		NetRecordName: ToNetRecordName(columnName),
		NetArgName:    ToNetArgName(columnName),
		SqlName:       columnName,
		IsArray:       col.IsArray,
		IsNullable:    !col.NotNull,
	}
}

func MapColumns(cols []*plugin.Column) []Field {
	res := make([]Field, 0)

	for i, c := range cols {
		res = append(res, MapColumn(i, c))
	}

	return res
}

func MapParameters(params []*plugin.Parameter) []Field {
	res := make([]Field, 0)

	for i, p := range params {
		res = append(res, MapColumn(i, p.Column))
	}

	return res
}
