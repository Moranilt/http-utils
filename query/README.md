# Query Builder
Simple query builder for SQL.

## Example

### 
```go
  query := New("SELECT * FROM test_table").LeftJoin("test_posts", "test_table.id=test_posts.user_id").Order("name", DESC).Limit("5").Offset("1")
  query.Where().
    OR(EQ("name", "testname"), EQ("age", "12")).
    AND(EQ("id", "123"), EQ("email", "test@mail.com")).
    LIKE("name", "%testname")
  fmt.Println(query.String())
  // Output: SELECT * FROM test_table LEFT JOIN test_posts ON test_table.id=test_posts.user_id WHERE (name='testname' OR age='12') AND (id='123' AND email='test@mail.com') AND name LIKE '%testname' ORDER BY name DESC LIMIT 5 OFFSET 1
```