package Mysqlqb

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// DBConnection holds the database connection.
var DBConnection *sql.DB

// ConnectMySQL initializes the MySQL connection with provided credentials.
func ConnectMySQL(username, password, host, dbname string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbname)

	var err error
	DBConnection, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error creating the database connection: %v", err)
	}

	if err := DBConnection.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
}

// CloseDB closes the database connection.
func CloseDB() {
	if DBConnection != nil {
		err := DBConnection.Close()
		if err != nil {
			log.Fatalf("Error closing the database connection: %v", err)
		}
	}
}

// QueryBuilder represents the query builder.
type QueryBuilder struct {
	table      string
	columns    []string
	joins      []string
	where      []string
	orderBy    string
	groupBy    string
	having     []string
	limit      int
	offset     int
	parameters []interface{}
}

// Table initializes the query builder with a table name.
func Table(table string) *QueryBuilder {
	return &QueryBuilder{
		table:  table,
		limit:  -1,
		offset: -1,
	}
}

// Select adds columns to the SELECT part of the query.
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.columns = append(qb.columns, columns...)

	return qb
}

// Where adds conditions to the WHERE clause.
func (qb *QueryBuilder) Where(condition string, params ...interface{}) *QueryBuilder {
	qb.where = append(qb.where, condition)
	qb.parameters = append(qb.parameters, params...)

	return qb
}

// WhereIn adds an IN condition to the WHERE clause.
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
		qb.parameters = append(qb.parameters, values[i])
	}
	qb.where = append(qb.where, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))

	return qb
}

// WhereNotIn adds a NOT IN condition to the WHERE clause.
func (qb *QueryBuilder) WhereNotIn(column string, values []interface{}) *QueryBuilder {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
		qb.parameters = append(qb.parameters, values[i])
	}
	qb.where = append(qb.where, fmt.Sprintf("%s NOT IN (%s)", column, strings.Join(placeholders, ", ")))

	return qb
}

// WhereNull adds a NULL condition to the WHERE clause.
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	qb.where = append(qb.where, fmt.Sprintf("%s IS NULL", column))

	return qb
}

// OrWhere adds an OR condition to the WHERE clause.
func (qb *QueryBuilder) OrWhere(condition string, params ...interface{}) *QueryBuilder {
	qb.where = append(qb.where, "OR "+condition)
	qb.parameters = append(qb.parameters, params...)

	return qb
}

// WhereLike Like adds a LIKE condition to the WHERE clause.
func (qb *QueryBuilder) WhereLike(column string, value string) *QueryBuilder {
	qb.where = append(qb.where, fmt.Sprintf("%s LIKE ?", column))
	qb.parameters = append(qb.parameters, value)

	return qb
}

// WhereNotLike WhereLikeNotLike NotLike adds a NOT LIKE condition to the WHERE clause.
func (qb *QueryBuilder) WhereNotLike(column string, value string) *QueryBuilder {
	qb.where = append(qb.where, fmt.Sprintf("%s NOT LIKE ?", column))
	qb.parameters = append(qb.parameters, value)

	return qb
}

// WhereBetween adds a BETWEEN condition to the WHERE clause.
func (qb *QueryBuilder) WhereBetween(column string, start, end interface{}) *QueryBuilder {
	qb.where = append(qb.where, fmt.Sprintf("%s BETWEEN ? AND ?", column))
	qb.parameters = append(qb.parameters, start, end)
	return qb
}

// DateBetween adds a BETWEEN condition for date columns in the WHERE clause.
func (qb *QueryBuilder) DateBetween(column string, start string, end string) *QueryBuilder {
	qb.where = append(qb.where, fmt.Sprintf("%s BETWEEN ? AND ?", column))
	qb.parameters = append(qb.parameters, start, end)

	return qb
}

// Join adds a JOIN clause to the query.
func (qb *QueryBuilder) Join(joinType, table, condition string) *QueryBuilder {
	join := fmt.Sprintf("%s JOIN %s ON %s", joinType, table, condition)
	qb.joins = append(qb.joins, join)

	return qb
}

// InnerJoin adds an INNER JOIN clause to the query.
func (qb *QueryBuilder) InnerJoin(table, condition string) *QueryBuilder {
	return qb.Join("INNER", table, condition)
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (qb *QueryBuilder) LeftJoin(table, condition string) *QueryBuilder {
	return qb.Join("LEFT", table, condition)
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (qb *QueryBuilder) RightJoin(table, condition string) *QueryBuilder {
	return qb.Join("RIGHT", table, condition)
}

// GroupBy adds GROUP BY clause to the query.
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupBy = strings.Join(columns, ", ")
	return qb
}

// Having adds HAVING clause to the query.
func (qb *QueryBuilder) Having(condition string, params ...interface{}) *QueryBuilder {
	qb.having = append(qb.having, condition)
	qb.parameters = append(qb.parameters, params...)
	return qb
}

// OrderBy sets the ORDER BY clause.
func (qb *QueryBuilder) OrderBy(order string) *QueryBuilder {
	qb.orderBy = order

	return qb
}

// Limit sets the LIMIT clause.
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit

	return qb
}

// Offset sets the OFFSET clause.
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset

	return qb
}

// Build generates the final SQL query and returns it with its parameters.
func (qb *QueryBuilder) Build() (string, []interface{}) {
	var query strings.Builder

	// SELECT clause
	if len(qb.columns) > 0 {
		query.WriteString("SELECT " + strings.Join(qb.columns, ", "))
	} else {
		query.WriteString("SELECT *")
	}

	// FROM clause
	query.WriteString(" FROM " + qb.table)

	// JOIN clauses
	if len(qb.joins) > 0 {
		query.WriteString(" " + strings.Join(qb.joins, " "))
	}

	// WHERE clause
	if len(qb.where) > 0 {
		query.WriteString(" WHERE " + strings.Join(qb.where, " AND "))
	}

	// ORDER BY clause
	if qb.orderBy != "" {
		query.WriteString(" ORDER BY " + qb.orderBy)
	}

	// LIMIT clause
	if qb.limit >= 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}

	// OFFSET clause
	if qb.offset >= 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", qb.offset))
	}

	return query.String(), qb.parameters
}

// Count returns the count of rows matching the query.
func (qb *QueryBuilder) Count() (int, error) {
	// Modify query to count rows
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", qb.BuildSelectQuery())
	params := qb.parameters

	var count int
	err := DBConnection.QueryRow(countQuery, params...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Error counting rows: %v", err)
	}

	return count, nil
}

func (qb *QueryBuilder) Sum(column string) (float64, error) {
	qb.columns = []string{"SUM(" + column + ")"}
	query, params := qb.Build()
	var sumValue float64
	err := DBConnection.QueryRow(query, params...).Scan(&sumValue)

	return sumValue, err
}

func (qb *QueryBuilder) Max(column string) (float64, error) {
	qb.columns = []string{"MAX(" + column + ")"}
	query, params := qb.Build()
	var maxValue float64
	err := DBConnection.QueryRow(query, params...).Scan(&maxValue)
	return maxValue, err
}

func (qb *QueryBuilder) Min(column string) (float64, error) {
	qb.columns = []string{"MIN(" + column + ")"}
	query, params := qb.Build()
	var minValue float64
	err := DBConnection.QueryRow(query, params...).Scan(&minValue)

	return minValue, err
}

func (qb *QueryBuilder) Avg(column string) (float64, error) {
	qb.columns = []string{"AVG(" + column + ")"}
	query, params := qb.Build()
	var avgValue float64
	err := DBConnection.QueryRow(query, params...).Scan(&avgValue)

	return avgValue, err
}

func (qb *QueryBuilder) Insert(data map[string]interface{}) (sql.Result, error) {
	columns := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	params := make([]interface{}, 0, len(data))

	for column, value := range data {
		columns = append(columns, column)
		placeholders = append(placeholders, "?")
		params = append(params, value)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", qb.table, strings.Join(columns, ","), strings.Join(placeholders, ","))

	return DBConnection.Exec(query, params...)
}

func (qb *QueryBuilder) BulkInsert(data []map[string]interface{}) (sql.Result, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data to insert")
	}

	columns := make([]string, 0)
	for column := range data[0] {
		columns = append(columns, column)
	}

	values := make([]string, 0)
	params := make([]interface{}, 0)

	for _, row := range data {
		placeholders := make([]string, len(row))
		for i, column := range columns {
			placeholders[i] = "?"
			params = append(params, row[column])
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", qb.table, strings.Join(columns, ","), strings.Join(values, ","))

	return DBConnection.Exec(query, params...)
}

func (qb *QueryBuilder) Update(data map[string]interface{}) (sql.Result, error) {
	setClauses := make([]string, 0)
	params := make([]interface{}, 0)

	for column, value := range data {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
		params = append(params, value)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", qb.table, strings.Join(setClauses, ","), strings.Join(qb.where, " AND "))

	return DBConnection.Exec(query, params...)
}

func (qb *QueryBuilder) Delete() (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s", qb.table)

	// Add WHERE clause if exists
	if len(qb.where) > 0 {
		query += " WHERE " + strings.Join(qb.where, " AND ")
	}

	// Print the query for debugging
	qb.PrintQuery()

	// Execute the query with the arguments
	return DBConnection.Exec(query, qb.parameters...)
}

// Get fetches multiple rows and returns them as an array of maps (like Laravel).
func (qb *QueryBuilder) Get() ([]map[string]interface{}, error) {
	query, params := qb.Build()
	rows, err := DBConnection.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Dynamically get column names and values
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		// Prepare a slice for the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create a map for the row
		row := make(map[string]interface{})
		for i, col := range columns {
			if b, ok := values[i].([]byte); ok { // Check if the value is a byte slice
				row[col] = string(b) // Convert byte slice to string
			} else {
				row[col] = values[i] // Otherwise, use the original value
			}
		}

		result = append(result, row)
	}

	return result, nil
}

// First fetches the first row of the result set.
func (qb *QueryBuilder) First() (map[string]interface{}, error) {
	query, params := qb.Build()
	row := DBConnection.QueryRow(query, params...)

	// Dynamically get column names and values
	columns, err := DBConnection.Query(query, params...) // Corrected to handle the error
	if err != nil {
		return nil, err
	}
	defer columns.Close()

	cols, err := columns.Columns() // Get column names
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Scan the row
	if err := row.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	// Create a map for the row
	result := make(map[string]interface{})
	for i, col := range cols {
		if b, ok := values[i].([]byte); ok { // Check if the value is a byte slice
			result[col] = string(b) // Convert byte slice to string
		} else {
			result[col] = values[i] // Otherwise, use the original value
		}
	}

	return result, nil
}

// BuildSelectQuery is a helper for building the core SELECT query.
func (qb *QueryBuilder) BuildSelectQuery() string {
	var query strings.Builder

	// SELECT clause
	if len(qb.columns) > 0 {
		query.WriteString("SELECT " + strings.Join(qb.columns, ", "))
	} else {
		query.WriteString("SELECT *")
	}

	// FROM clause
	query.WriteString(" FROM " + qb.table)

	// JOIN clauses
	if len(qb.joins) > 0 {
		query.WriteString(" " + strings.Join(qb.joins, " "))
	}

	// WHERE clause
	if len(qb.where) > 0 {
		query.WriteString(" WHERE " + strings.Join(qb.where, " AND "))
	}

	return query.String()
}

// PrintQuery prints the built raw SQL query and its parameters.
func (qb *QueryBuilder) PrintQuery() {
	query, params := qb.Build()
	fmt.Println(query, params)
}
