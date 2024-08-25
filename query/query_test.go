package query

import (
	"fmt"
	"testing"
)

type testItem struct {
	name     string
	callback func(t *testing.T) string
	expected string
}

var tests = []testItem{
	{
		name: "default query",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "limit",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Limit("10")
			return query.String()
		},
		expected: "SELECT * FROM test_table LIMIT 10",
	},
	{
		name: "offset",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Offset("1")
			return query.String()
		},
		expected: "SELECT * FROM test_table OFFSET 1",
	},
	{
		name: "limit offset",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Limit("5").Offset("1")
			return query.String()
		},
		expected: "SELECT * FROM test_table LIMIT 5 OFFSET 1",
	},
	{
		name: "order by",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Order("name", DESC)
			return query.String()
		},
		expected: "SELECT * FROM test_table ORDER BY name DESC",
	},
	{
		name: "limit offset order by",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Order("name", DESC).Limit("5").Offset("1")
			return query.String()
		},
		expected: "SELECT * FROM test_table ORDER BY name DESC LIMIT 5 OFFSET 1",
	},
	{
		name: "where EQ",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().EQ("name", "testname")
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE name = 'testname'",
	},
	{
		name: "where LIKE",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().LIKE("name", "%testname")
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE name LIKE '%testname'",
	},
	{
		name: "where OR",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().OR(EQ("name", "testname"), EQ("age", "12"))
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE (name = 'testname' OR age = '12')",
	},
	{
		name: "where OR AND",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().OR(EQ("name", "testname"), EQ("age", "12")).AND(EQ("id", "123"), EQ("email", "test@mail.com"))
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE (name = 'testname' OR age = '12') AND (id = '123' AND email = 'test@mail.com')",
	},
	{
		name: "not valid ORDER",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Order("test", "random")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid LIMIT",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Limit("s")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid OFFSET",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Offset("s")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid OR",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().OR(EQ("name", "testname"))
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid AND",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().AND(EQ("name", "testname"))
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid DEFAULT AND",
		callback: func(t *testing.T) string {
			query := AND(EQ("name", "testname"))
			return query
		},
		expected: "",
	},
	{
		name: "not valid DEFAULT OR",
		callback: func(t *testing.T) string {
			query := OR(EQ("name", "testname"))
			return query
		},
		expected: "",
	},
	{
		name: "invalid EQ",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().EQ("", "test")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "invalid LIKE",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().LIKE("", "test")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "invalid ORDER",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Order("", "random")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "invalid OR",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().OR()
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "invalid AND",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().AND()
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "nested OR AND",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().OR(EQ("name", "testname"), AND(EQ("age", "10"), EQ("city", "New York")))
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE (name = 'testname' OR (age = '10' AND city = 'New York'))",
	},
	{
		name: "nested AND OR",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().AND(EQ("name", "testname"), OR(EQ("age", "10"), EQ("city", "New York")))
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE (name = 'testname' AND (age = '10' OR city = 'New York'))",
	},
	{
		name: "inner join",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").InnerJoin("posts", "users.id = posts.user_id")
			return query.String()
		},
		expected: "SELECT * FROM users INNER JOIN posts ON users.id = posts.user_id",
	},
	{
		name: "left join",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").LeftJoin("posts", "users.id = posts.user_id")
			return query.String()
		},
		expected: "SELECT * FROM users LEFT JOIN posts ON users.id = posts.user_id",
	},
	{
		name: "right join",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").RightJoin("posts", "users.id = posts.user_id")
			return query.String()
		},
		expected: "SELECT * FROM users RIGHT JOIN posts ON users.id = posts.user_id",
	},
	{
		name: "full join",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").FullJoin("posts", "users.id = posts.user_id")
			return query.String()
		},
		expected: "SELECT * FROM users FULL JOIN posts ON users.id = posts.user_id",
	},
	{
		name: "cross join",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").CrossJoin("posts")
			return query.String()
		},
		expected: "SELECT * FROM users CROSS JOIN posts",
	},
	{
		name: "group by",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").GroupBy("name")
			return query.String()
		},
		expected: "SELECT * FROM users GROUP BY name",
	},
	{
		name: "group by multiple",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").GroupBy("name", "age")
			return query.String()
		},
		expected: "SELECT * FROM users GROUP BY name, age",
	},
	{
		name: "group by without columns",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").GroupBy()
			return query.String()
		},
		expected: "SELECT * FROM users",
	},
	{
		name: "left join empty",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").LeftJoin("", "users.id = posts.user_id")
			return query.String()
		},
		expected: "SELECT * FROM users",
	},
	{
		name: "right join empty",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").RightJoin("", "users.id = posts.user_id")
			return query.String()
		},
		expected: "SELECT * FROM users",
	},
	{
		name: "full join empty",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").FullJoin("", "users.id = posts.user_id")
			return query.String()
		},
		expected: "SELECT * FROM users",
	},
	{
		name: "cross join empty",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").CrossJoin("")
			return query.String()
		},
		expected: "SELECT * FROM users",
	},
	{
		name: "inner join empty",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").InnerJoin("", "users.id = posts.user_id")
			return query.String()
		},
		expected: "SELECT * FROM users",
	},
	{
		name: "having empty",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").Having("")
			return query.String()
		},
		expected: "SELECT * FROM users",
	},
	{
		// HAVING after WHERE
		name: "chained having",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").GroupBy("name").Having("COUNT(id) > 10")
			query.Where().EQ("name", "John")
			return query.String()
		},
		expected: "SELECT * FROM users WHERE name = 'John' GROUP BY name HAVING COUNT(id) > 10",
	},
	{
		// HAVING after GROUP BY
		name: "having after group by",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").GroupBy("name").Having("COUNT(id) > 10")
			return query.String()
		},
		expected: "SELECT * FROM users GROUP BY name HAVING COUNT(id) > 10",
	},
	{
		name: "having without group by",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").Having("COUNT(*) > 10")
			return query.String()
		},
		expected: "SELECT * FROM users",
	},
	{
		name: "using IS NULL",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users")
			query.Where().IS("name", nil)
			return query.String()
		},
		expected: "SELECT * FROM users WHERE name IS NULL",
	},
	{
		name: "using IS NULL where name is pointer string",
		callback: func(t *testing.T) string {
			var name *string
			query := New("SELECT * FROM users")
			query.Where().IS("name", name)
			return query.String()
		},
		expected: "SELECT * FROM users WHERE name IS NULL",
	},
	{
		name: "using IS NULL and OR",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users")
			query.Where().OR(IS("name", nil), EQ("name", "John"))
			return query.String()
		},
		expected: "SELECT * FROM users WHERE (name IS NULL OR name = 'John')",
	},
	{
		name: "using EQ with int type",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users")
			query.Where().EQ("age", 12)
			return query.String()
		},
		expected: "SELECT * FROM users WHERE age = 12",
	},
	{
		name: "using EQ with float64 type",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users")
			query.Where().EQ("age", 12.5)
			return query.String()
		},
		expected: "SELECT * FROM users WHERE age = 12.5",
	},
	{
		name: "using EQ with bool type",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").Where().EQ("blocked", true).Query()
			return query.String()
		},
		expected: "SELECT * FROM users WHERE blocked = true",
	},
	{
		name: "insert columns",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "age")
			return query.String()
		},
		expected: "INSERT INTO table_name (name, age)",
	},
	{
		name: "insert columns with single values with int",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "age").Values("John", 12)
			return query.String()
		},
		expected: "INSERT INTO table_name (name, age) VALUES ('John', 12)",
	},
	{
		name: "insert columns with single values with boolean",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "valid").Values("John", true).Values("Jane", false)
			return query.String()
		},
		expected: "INSERT INTO table_name (name, valid) VALUES ('John', true), ('Jane', false)",
	},
	{
		name: "insert columns with single values with null",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "valid").Values("John", nil)
			return query.String()
		},
		expected: "INSERT INTO table_name (name, valid) VALUES ('John', NULL)",
	},
	{
		name: "insert columns with single values with float64 and float32",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "cost").Values("John", float32(12.5)).Values("Jane", float64(13.55557))
			return query.String()
		},
		expected: "INSERT INTO table_name (name, cost) VALUES ('John', 12.5), ('Jane', 13.55557)",
	},
	{
		name: "insert columns with single values with now() function",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "date").Values("John", "now()")
			return query.String()
		},
		expected: "INSERT INTO table_name (name, date) VALUES ('John', NOW())",
	},
	{
		name: "insert columns with single values with uuid() function",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "id").Values("John", "uuid()")
			return query.String()
		},
		expected: "INSERT INTO table_name (name, id) VALUES ('John', UUID())",
	},
	{
		name: "insert columns with single values with question mark",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "id").Values("?", "?")
			return query.String()
		},
		expected: "INSERT INTO table_name (name, id) VALUES (?, ?)",
	},
	{
		name: "insert columns with multiple values",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "age").Values("John", 12).Values("Doe", 13).Values("Jane", 14)
			return query.String()
		},
		expected: "INSERT INTO table_name (name, age) VALUES ('John', 12), ('Doe', 13), ('Jane', 14)",
	},
	{
		name: "insert columns with multiple values with returning",
		callback: func(t *testing.T) string {
			query := New("INSERT INTO table_name").InsertColumns("name", "age").Values("John", 12).Values("Doe", 13).Values("Jane", 14).Returning("id", "name", "age")
			return query.String()
		},
		expected: "INSERT INTO table_name (name, age) VALUES ('John', 12), ('Doe', 13), ('Jane', 14) RETURNING id, name, age",
	},
	{
		name: "update with set value",
		callback: func(t *testing.T) string {
			query := New("UPDATE table_name").Set("name", "John")
			return query.String()
		},
		expected: "UPDATE table_name SET name = 'John'",
	},
	{
		name: "update with set value with where",
		callback: func(t *testing.T) string {
			query := New("UPDATE table_name").Set("name", "John").Where().EQ("age", 12).Query()
			return query.String()
		},
		expected: "UPDATE table_name SET name = 'John' WHERE age = 12",
	},
	{
		name: "update with set multiple value with where",
		callback: func(t *testing.T) string {
			query := New("UPDATE table_name").Set("name", "John").Set("age", 20).Where().EQ("age", 12).Query()
			return query.String()
		},
		expected: "UPDATE table_name SET name = 'John', age = 20 WHERE age = 12",
	},
	{
		name: "update with set multiple value with where and returning",
		callback: func(t *testing.T) string {
			query := New("UPDATE table_name").Set("name", "John").Set("age", 20).Where().EQ("age", 12).Query().Returning("id", "name", "age")
			return query.String()
		},
		expected: "UPDATE table_name SET name = 'John', age = 20 WHERE age = 12 RETURNING id, name, age",
	},
	{
		name: "where in clause with single value",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM table_name").Where().IN("name", "John").Query()
			return query.String()
		},
		expected: "SELECT * FROM table_name WHERE name IN ('John')",
	},
	{
		name: "where in clause with array value",
		callback: func(t *testing.T) string {
			names := []string{"John", "Jane", "Kevin"}
			query := New("SELECT * FROM table_name").Where().IN("name", names).Query()
			return query.String()
		},
		expected: "SELECT * FROM table_name WHERE name IN ('John','Jane','Kevin')",
	},
	{
		name: "where in clause with two array values",
		callback: func(t *testing.T) string {
			names := []string{"John", "Jane", "Kevin"}
			query := New("SELECT * FROM table_name").Where().IN("name", names, names).Query()
			return query.String()
		},
		expected: "SELECT * FROM table_name WHERE name IN ('John','Jane','Kevin','John','Jane','Kevin')",
	},
	{
		name: "where in clause with multiple values",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM table_name").Where().IN("name", "John", "Jane", "Kevin").Query()
			return query.String()
		},
		expected: "SELECT * FROM table_name WHERE name IN ('John','Jane','Kevin')",
	},
	{
		name: "where in clause with zero values",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM table_name").Where().IN("name").Query()
			return query.String()
		},
		expected: "SELECT * FROM table_name",
	},
	{
		name: "IN combined with other conditions",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM users").Where().EQ("active", true).IN("role", "admin", "moderator").Query()
			return query.String()
		},
		expected: "SELECT * FROM users WHERE active = true AND role IN ('admin','moderator')",
	},
}

func TestQuery(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.callback(t)
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

func ExampleQuery_String() {
	query := New("SELECT * FROM test_table").LeftJoin("test_posts", "test_table.id=test_posts.user_id").Order("name", DESC).Limit("5").Offset("1").Where().
		OR(EQ("name", "testname"), EQ("age", "12")).
		AND(EQ("id", "123"), EQ("email", "test@mail.com")).
		IN("name", "John", "Doe", "Jane").
		LIKE("name", "%testname").Query()
	fmt.Println(query.String())
	// Output: SELECT * FROM test_table LEFT JOIN test_posts ON test_table.id=test_posts.user_id WHERE (name = 'testname' OR age = '12') AND (id = '123' AND email = 'test@mail.com') AND name IN ('John','Doe','Jane') AND name LIKE '%testname' ORDER BY name DESC LIMIT 5 OFFSET 1
}

func ExampleQuery_InsertColumns() {
	query := New("INSERT INTO table_name").InsertColumns("name", "age").
		Values("John", 12).Values("Doe", 13).Values("Jane", 14).
		Returning("id", "name", "age")
	fmt.Println(query.String())
	// Output: INSERT INTO table_name (name, age) VALUES ('John', 12), ('Doe', 13), ('Jane', 14) RETURNING id, name, age
}

func ExampleQuery_Where() {
	query := New("SELECT * FROM test_table").Where().OR(EQ("name", "testname"), EQ("age", "12")).Query()
	fmt.Println(query.String())

	query.Where().AND(EQ("id", "123"), EQ("email", "test@mail.com"))
	fmt.Println(query.String())

	query.Where().LIKE("name", "%testname").Query()
	fmt.Println(query.String())
	// Output:
	// SELECT * FROM test_table WHERE (name = 'testname' OR age = '12')
	// SELECT * FROM test_table WHERE (name = 'testname' OR age = '12') AND (id = '123' AND email = 'test@mail.com')
	// SELECT * FROM test_table WHERE (name = 'testname' OR age = '12') AND (id = '123' AND email = 'test@mail.com') AND name LIKE '%testname'
}

func ExampleQuery_Set() {
	query := New("UPDATE table_name").Set("name", "John").Where().EQ("age", 12).Query()
	fmt.Println(query.String())

	query.Set("age", 20)
	fmt.Println(query.String())

	query.Set("country", "USA").Set("city", "New York").Set("accepted", true).Set("updated_at", "now()")
	fmt.Println(query.String())
	// Output:
	// UPDATE table_name SET name = 'John' WHERE age = 12
	// UPDATE table_name SET name = 'John', age = 20 WHERE age = 12
	// UPDATE table_name SET name = 'John', age = 20, country = 'USA', city = 'New York', accepted = true, updated_at = NOW() WHERE age = 12
}
