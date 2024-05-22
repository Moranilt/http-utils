package query

import (
	"fmt"
	"strings"

	"github.com/Moranilt/http-utils/validators"
)

const (
	// DESC is descending order
	DESC = "DESC"
	// ASC is ascending order
	ASC = "ASC"
)

// ValidOrderType checks if the given string is a valid ordering type
func ValidOrderType(val string) bool {
	val = strings.ToUpper(val)
	return val == DESC || val == ASC
}

type Query struct {
	// main is the base SQL query string
	main []string

	// where is the WHERE clause
	where *Where

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
}

// New creates a new Query with the given base SQL query string
func New(q string) *Query {
	return &Query{
		main:  []string{q},
		where: &Where{},
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

// Where returns the WHERE clause for adding conditions
func (q *Query) Where() *Where {
	return q.where
}

// String returns the full SQL query string
func (q *Query) String() string {
	var where string
	var groupBy string
	var having string

	if len(q.where.chunks) > 0 {
		where = "WHERE " + strings.Join(q.where.chunks, " AND ")
	}

	if q.groupBy != "" {
		groupBy = q.groupBy
	}

	if q.having != "" && q.groupBy != "" {
		having = q.having
	}

	var (
		items     = []string{where, groupBy, q.order, having, q.limit, q.offset}
		mainQuery = strings.Join(q.main, " ")

		result strings.Builder
	)

	result.WriteString(mainQuery)

	for _, item := range items {
		if len(item) != 0 {
			result.WriteString(" " + item)
		}
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

// The Where struct represents the WHERE clause in a SQL query.
// It contains a chunks slice to hold the individual WHERE conditions.
type Where struct {
	chunks []string
}

// EQ adds an equality condition to the WHERE clause.
func (w *Where) EQ(fieldName string, value any) *Where {
	if isEmpty(fieldName) {
		return w
	}

	w.chunks = append(w.chunks, EQ(fieldName, value))
	return w
}

// LIKE adds a LIKE condition to the WHERE clause.
func (w *Where) LIKE(fieldName string, value any) *Where {
	if isEmpty(fieldName) {
		return w
	}
	w.chunks = append(w.chunks, LIKE(fieldName, value))
	return w
}

// OR adds an OR condition to the WHERE clause.
func (w *Where) OR(args ...string) *Where {
	if len(args) < 2 {
		return w
	}
	w.chunks = append(w.chunks, OR(args...))
	return w
}

// AND adds an AND condition to the WHERE clause.
func (w *Where) AND(args ...string) *Where {
	if len(args) < 2 {
		return w
	}
	w.chunks = append(w.chunks, AND(args...))
	return w
}

func (w *Where) IS(fieldName string, value any) *Where {
	if isEmpty(fieldName) {
		return w
	}

	w.chunks = append(w.chunks, IS(fieldName, value))
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
		panic("not supported type. Supports only numbers, string, bool")
	}
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
