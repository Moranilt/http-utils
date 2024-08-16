package query

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Moranilt/http-utils/validators"
)

const (
	// DESC is descending order
	DESC = "DESC"
	// ASC is ascending order
	ASC = "ASC"
)

type WhereClause interface {
	// EQ creates an equality condition for the specified field and value
	EQ(fieldName string, value any) WhereClause

	// LIKE creates a LIKE condition for the specified field and value
	LIKE(fieldName string, value any) WhereClause

	// OR creates an OR condition with the given arguments
	OR(args ...string) WhereClause

	// AND creates an AND condition with the given arguments
	AND(args ...string) WhereClause

	// IS creates an IS condition for the specified field and value
	IS(fieldName string, value any) WhereClause

	// Query returns the Query object associated with this WhereClause
	Query() *Query
}

// ValidOrderType checks if the given string is a valid ordering type
func ValidOrderType(val string) bool {
	val = strings.ToUpper(val)
	return val == DESC || val == ASC
}

type Query struct {
	// main is the base SQL query string
	main []string

	// where is the WHERE clause
	where []string

	// groupBy is the GROUP BY clause
	groupBy string

	// other fields
	having string

	// order is the ORDER BY clause
	order string

	// limit is the LIMIT clause
	limit string

	// offset is the OFFSET clause
	offset string

	// insertColumns is the columns to be inserted into the table
	insertColumns []string

	// values is the values to be inserted into the table
	values [][]any

	// returning is the RETURNING clause
	returning []string
}

// New creates a new Query with the given base SQL query string
func New(q string) *Query {
	return &Query{
		main: []string{q},
	}
}

func (q *Query) InnerJoin(table string, on string) *Query {
	if isEmpty(table) || isEmpty(on) {
		return q
	}

	q.main = append(q.main, fmt.Sprintf("INNER JOIN %s ON %s", table, on))
	return q
}

func (q *Query) LeftJoin(table string, on string) *Query {
	if isEmpty(table) || isEmpty(on) {
		return q
	}

	q.main = append(q.main, fmt.Sprintf("LEFT JOIN %s ON %s", table, on))
	return q
}

func (q *Query) RightJoin(table string, on string) *Query {
	if isEmpty(table) || isEmpty(on) {
		return q
	}

	q.main = append(q.main, fmt.Sprintf("RIGHT JOIN %s ON %s", table, on))
	return q
}

func (q *Query) FullJoin(table string, on string) *Query {
	if isEmpty(table) || isEmpty(on) {
		return q
	}

	q.main = append(q.main, fmt.Sprintf("FULL JOIN %s ON %s", table, on))
	return q
}

func (q *Query) CrossJoin(table string) *Query {
	if isEmpty(table) {
		return q
	}

	q.main = append(q.main, fmt.Sprintf("CROSS JOIN %s", table))
	return q
}

func (q *Query) Returning(fields ...string) *Query {
	if len(fields) == 0 {
		return q
	}

	q.returning = fields
	return q
}

func (q *Query) Query() *Query {
	return q
}

// Sets the columns to be inserted in an INSERT statement.
func (q *Query) InsertColumns(columns ...string) *Query {
	q.insertColumns = columns
	return q
}

// Appends the provided values to the list of values to be inserted in an
// INSERT statement.
func (q *Query) Values(values ...any) *Query {
	q.values = append(q.values, values)
	return q
}

// Where returns the WHERE clause for adding conditions
func (q *Query) Where() WhereClause {
	return q
}

// Returns the full SQL query string
func (q *Query) String() string {
	var result strings.Builder

	result.WriteString(strings.Join(q.main, " "))

	appendIfNotEmpty := func(s string) {
		if s != "" {
			result.WriteByte(' ')
			result.WriteString(s)
		}
	}

	if len(q.where) > 0 {
		result.WriteString(" WHERE ")
		result.WriteString(strings.Join(q.where, " AND "))
	}

	appendIfNotEmpty(q.groupBy)

	if q.having != "" && q.groupBy != "" {
		appendIfNotEmpty(q.having)
	}

	if len(q.insertColumns) > 0 {
		result.WriteString(" (")
		result.WriteString(strings.Join(q.insertColumns, ", "))
		result.WriteByte(')')
	}

	if len(q.values) > 0 {
		result.WriteByte(' ')
		result.WriteString(buildValues(q.values))
	}

	appendIfNotEmpty(q.order)
	appendIfNotEmpty(q.limit)
	appendIfNotEmpty(q.offset)

	if len(q.returning) > 0 {
		result.WriteString(" RETURNING ")
		result.WriteString(strings.Join(q.returning, ", "))
	}

	return result.String()
}

// Order adds an ORDER BY clause
func (q *Query) Order(by string, orderType string) *Query {
	if !ValidOrderType(orderType) {
		return q
	}
	q.order = fmt.Sprintf("ORDER BY %s %s", by, orderType)
	return q
}

// Limit adds a LIMIT clause
func (q *Query) Limit(val string) *Query {
	if !validators.ValidInt(val) {
		return q
	}
	q.limit = fmt.Sprintf("LIMIT %s", val)
	return q
}

// Offset adds an OFFSET clause
func (q *Query) Offset(val string) *Query {
	if !validators.ValidInt(val) {
		return q
	}
	q.offset = fmt.Sprintf("OFFSET %s", val)
	return q
}

// EQ adds an equality condition to the WHERE clause.
func (w *Query) EQ(fieldName string, value any) WhereClause {
	if isEmpty(fieldName) {
		return w
	}

	w.where = append(w.where, EQ(fieldName, value))
	return w
}

// LIKE adds a LIKE condition to the WHERE clause.
func (w *Query) LIKE(fieldName string, value any) WhereClause {
	if isEmpty(fieldName) {
		return w
	}
	w.where = append(w.where, LIKE(fieldName, value))
	return w
}

// OR adds an OR condition to the WHERE clause.
func (w *Query) OR(args ...string) WhereClause {
	if len(args) < 2 {
		return w
	}
	w.where = append(w.where, OR(args...))
	return w
}

// AND adds an AND condition to the WHERE clause.
func (w *Query) AND(args ...string) WhereClause {
	if len(args) < 2 {
		return w
	}
	w.where = append(w.where, AND(args...))
	return w
}

func (w *Query) IS(fieldName string, value any) WhereClause {
	if isEmpty(fieldName) {
		return w
	}

	w.where = append(w.where, IS(fieldName, value))
	return w
}

// AND creates an AND condition string from the provided arguments.
func AND(args ...string) string {
	if len(args) < 2 {
		return ""
	}
	return fmt.Sprintf("(%s)", strings.Join(args, " AND "))
}

// EQ creates an equality condition string.
func EQ(fieldName string, value any) string {
	return buildComparison(fieldName, value, "=")
}

// LIKE creates a LIKE condition string.
func LIKE(fieldName string, value any) string {
	return buildComparison(fieldName, value, "LIKE")
}

func buildComparison(fieldName string, value any, comparison string) string {
	if isEmpty(fieldName) {
		panic("fieldName cannot be empty")
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			value = nil
		} else if v.IsValid() {
			value = v.Elem().Interface()
		} else {
			panic("not supported type. Supports only numbers, string, bool, nil")
		}
	}
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%s %s %d", fieldName, comparison, value)
	case float32, float64:
		return fmt.Sprintf("%s %s %v", fieldName, comparison, value)
	case string:
		if value == "NULL" || value == "null" {
			return fmt.Sprintf("%s %s NULL", fieldName, comparison)
		}
		return fmt.Sprintf("%s %s '%s'", fieldName, comparison, value)
	case bool:
		return fmt.Sprintf("%s %s %t", fieldName, comparison, value)
	case nil:
		return fmt.Sprintf("%s %s NULL", fieldName, comparison)
	default:
		panic("not supported type. Supports only numbers, string, bool, nil")
	}
}

func wrapValue(value any) string {
	if value == nil {
		return "NULL"
	}

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "NULL"
		}
		value = v.Elem().Interface()
	}

	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)
	case string:
		switch strings.ToUpper(v) {
		case "NULL", "NOW()", "UUID()":
			return strings.ToUpper(v)
		case "?":
			return "?"
		default:
			return "'" + v + "'"
		}
	case bool:
		return strconv.FormatBool(v)
	default:
		panic("unsupported type: only numbers, string, bool, nil are supported")
	}
}

func buildValues(values [][]any) string {
	var builder strings.Builder
	builder.WriteString("VALUES ")

	for i, row := range values {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteByte('(')
		for j, value := range row {
			if j > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(wrapValue(value))
		}
		builder.WriteByte(')')
	}

	return builder.String()
}

// OR creates an OR condition string from the provided arguments.
func OR(args ...string) string {
	if len(args) < 2 {
		return ""
	}
	return fmt.Sprintf("(%s)", strings.Join(args, " OR "))
}

// IS creates an IS condition string from the provided arguments.
func IS(fieldName string, value any) string {
	return buildComparison(fieldName, value, "IS")
}

// GroupBy adds a GROUP BY clause to the query
func (q *Query) GroupBy(columns ...string) *Query {
	if len(columns) == 0 {
		return q
	}

	q.groupBy = fmt.Sprintf("GROUP BY %s", strings.Join(columns, ", "))
	return q
}

// Having adds a HAVING clause to the query
func (q *Query) Having(condition string) *Query {
	if condition == "" {
		return q
	}

	q.having = fmt.Sprintf("HAVING %s", condition)
	return q
}

func isEmpty(s string) bool {
	return len(s) == 0
}
