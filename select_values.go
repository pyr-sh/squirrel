package squirrel

import (
	"fmt"
	"strconv"
	"strings"
)

// e.g. SELECT FROM (VALUES ('a', 'b'), ('c', 'd')) AS (col1, col2)
type selectValues struct {
	values [][]interface{}
}

func (sv selectValues) ToSql() (string, []interface{}, error) {
	var res strings.Builder
	var args []interface{}
	res.WriteString("VALUES ")
	for i, row := range sv.values {
		if i != 0 {
			res.WriteString(", ")
		}

		res.WriteRune('(')
		for j, col := range row {
			if j != 0 {
				res.WriteString(", ")
			}
			switch v := col.(type) {
			case Sqlizer:
				q, a, err := v.ToSql()
				if err != nil {
					return "", nil, err
				}
				res.WriteString(q)
				args = append(args, a...)
			case string:
				res.WriteRune('\'')
				res.WriteString(v)
				res.WriteRune('\'')
			case int:
				res.WriteString(strconv.Itoa(v))
			case int8:
				res.WriteString(strconv.FormatInt(int64(v), 10))
			case int16:
				res.WriteString(strconv.FormatInt(int64(v), 10))
			case int32:
				res.WriteString(strconv.FormatInt(int64(v), 10))
			case int64:
				res.WriteString(strconv.FormatInt(int64(v), 10))
			case uint:
				res.WriteString(strconv.FormatUint(uint64(v), 10))
			case uint8:
				res.WriteString(strconv.FormatUint(uint64(v), 10))
			case uint16:
				res.WriteString(strconv.FormatUint(uint64(v), 10))
			case uint32:
				res.WriteString(strconv.FormatUint(uint64(v), 10))
			case uint64:
				res.WriteString(strconv.FormatUint(uint64(v), 10))
			default:
				return "", nil,
					fmt.Errorf("unhandled type %T, consider using a string and casting it to the desired type", v)
			}
		}
		res.WriteRune(')')
	}
	return res.String(), args, nil
}

type fromValuesSelectAlias struct {
	aliasExpr
	definition []string
}

func (a fromValuesSelectAlias) ToSql() (string, []interface{}, error) {
	sql, args, err := a.aliasExpr.ToSql()
	if err != nil {
		return "", nil, err
	}
	if len(a.definition) > 0 {
		var b strings.Builder
		b.WriteRune('(')
		for i, col := range a.definition {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(col)
		}
		b.WriteRune(')')
		sql += b.String()
	}
	return sql, args, err
}
